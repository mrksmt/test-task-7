package client

import (
	"context"
	"fmt"
	"log"
	"net"

	"github.com/pkg/errors"
	"go.uber.org/ratelimit"

	"github.com/mrksmt/test-task-7/pkg/helpers"
	"github.com/mrksmt/test-task-7/pkg/pow"
)

type Parameters struct {
	Address string `envconfig:"SENTENCE_SERVICE_ADDRESS" default:"localhost:8090"`
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

		sentence, err := c.GetSentence(context.WithoutCancel(ctx))
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

	log.Println("client established connection", conn.LocalAddr().String())

	// get challenge code

	challengeCode, err := helpers.Receive(conn)
	if err != nil {
		return "", errors.Wrap(err, "deal tcp")
	}
	log.Println("client got  challenge code", challengeCode)

	// make proof of work response

	proof, err := pow.GetProofAsync(ctx, challengeCode)
	if err != nil {
		return "", errors.Wrap(err, "deal tcp")
	}

	// send proof of work

	_, err = conn.Write(proof)
	if err != nil {
		return "", errors.Wrap(err, "deal tcp")
	}
	log.Println("client send challenge resp", proof)

	// get server response

	resp, err := helpers.Receive(conn)
	if err != nil {
		return "", errors.Wrap(err, "receive sentence")
	}

	return string(resp), nil
}
