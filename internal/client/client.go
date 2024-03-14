package client

import (
	"context"
	"fmt"
	"io"
	"log"
	"net"
	"time"

	"github.com/pkg/errors"
	"go.uber.org/ratelimit"

	"github.com/mrksmt/test-task-7/pkg/helpers"
	"github.com/mrksmt/test-task-7/pkg/pow"
)

type Parameters struct {
	Address      string        `envconfig:"SENTENCE_SERVICE_ADDRESS" default:"localhost:8090"`
	ReadTimeout  time.Duration `envconfig:"CLIENT_READ_TIMEOUT" default:"1s"`
	WriteTimeout time.Duration `envconfig:"CLIENT_WRITE_TIMEOUT" default:"1s"`
}

type SentenceClient struct {
	address string
}

func NewSentenceClient(
	params *Parameters,
) *SentenceClient {
	c := &SentenceClient{address: params.Address}
	return c
}

// Run implements slice.Dispatcher
func (c *SentenceClient) Run(
	ctx context.Context,
) error {

	rl := ratelimit.New(1)

	for {

		rl.Take()

		select {
		case <-ctx.Done():
			return nil
		default:
		}

		fmt.Println()

		sentence, err := c.GetSentence(ctx)
		if err != nil {
			log.Println("client get sentence error:", err)
			continue
		}

		log.Println("client got  sentence:", sentence)
	}
}

// GetSentence get one sentence from server.
// Implements proof of work handshake TCP connection.
func (c *SentenceClient) GetSentence(
	ctx context.Context,
) (string, error) {

	// connect to the server

	conn, err := net.Dial("tcp", c.address)
	if err != nil {
		return "", errors.Wrapf(err, "deal tcp address: %s", c.address)
	}
	defer conn.Close()

	ctx, cancel := context.WithCancel(ctx)
	go func() {
		rl := ratelimit.New(10)
		for {

			for {
				rl.Take()
				select {
				case <-ctx.Done():
					return
				default:
					err := helpers.ConnCheck(conn)
					if err != nil {
						cancel()
						return
					}
				}
			}
		}
	}()

	log.Println("client established connection", conn.LocalAddr().String())

	// get challenge code

	challengeCode, err := helpers.Receive(conn)
	if err != nil {
		return "", errors.Wrap(err, "get challenge code")
	}
	log.Println("client got  challenge code", challengeCode)

	// make proof of work response

	proof, err := pow.GetProofAsync(ctx, challengeCode)
	if err != nil {
		return "", errors.Wrap(err, "get proof async")
	}

	// send proof of work

	err = helpers.ConnCheck(conn)
	if err != nil {
		return "", errors.Wrap(err, "client send challenge resp")
	}

	_, err = conn.Write(proof)
	if err != nil {
		return "", errors.Wrap(err, "client send challenge resp")
	}
	log.Println("client send challenge resp", proof)

	// get server response

	resp, err := helpers.Receive(conn)
	if err != nil && !errors.Is(err, io.EOF) {
		return "", errors.Wrap(err, "receive sentence")
	}

	return string(resp), nil
}
