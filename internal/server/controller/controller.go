package controller

import (
	"context"
	"log"
	"net"
	"sync"
	"time"

	"github.com/google/uuid"
	"github.com/pkg/errors"

	"github.com/mrksmt/test-task-7/internal/server/service"
	"github.com/mrksmt/test-task-7/pkg/helpers"
	"github.com/mrksmt/test-task-7/pkg/pow"
)

type Parameters struct {
	Address          string        `envconfig:"SENTENCE_SERVICE_ADDRESS" default:"localhost:8090"`
	WriteTimeout     time.Duration `envconfig:"SERVICE_WRITE_TIMEOUT" default:"1s"`
	ChallengeTimeout time.Duration `envconfig:"SERVICE_CHALLENGE_TIMEOUT" default:"10s"`
}

type Controller struct {
	params  *Parameters
	service service.Service
}

func NewController(
	params *Parameters,
	service service.Service,
) *Controller {

	s := &Controller{
		params:  params,
		service: service,
	}
	return s
}

// Run implements slice.Dispatcher
func (s *Controller) Run(
	ctx context.Context,
) error {

	listenerWG := new(sync.WaitGroup)
	handlerWG := new(sync.WaitGroup)

	// listen for incoming connections
	listener, err := net.Listen("tcp", s.params.Address)
	if err != nil {
		return errors.Wrap(err, "make tcp listener error")
	}

	log.Println("server is listening on", listener.Addr().String())

	listenerWG.Add(1)
	go func() {
		defer listenerWG.Done()
		for {

			select {
			case <-ctx.Done():
				return
			default:
			}

			// accept incoming connections
			conn, err := listener.Accept()
			if err != nil {
				select {
				case <-ctx.Done():
					return
				default:
					log.Println("listener accept error:", err)
					continue
				}
			}

			// handle client connection in a goroutine
			handlerWG.Add(1)
			go s.handleClient(context.WithoutCancel(ctx), handlerWG, conn)
		}
	}()

	<-ctx.Done()
	handlerWG.Wait()

	err = listener.Close()
	if err != nil {
		log.Println("listener close error:", err)
	}

	listenerWG.Wait()
	return nil
}

func (c *Controller) handleClient(
	ctx context.Context,
	wg *sync.WaitGroup,
	conn net.Conn,
) {
	defer wg.Done()
	defer conn.Close()

	log.Println("server established connection", conn.LocalAddr().String())

	// send challenge code

	challengeCode := c.getChallengeCode()

	conn.SetWriteDeadline(time.Now().Add(c.params.WriteTimeout))
	_, err := conn.Write(challengeCode)
	if err != nil {
		log.Printf("write challenge code error: %s", err)
		return
	}
	log.Println("server send challenge code", challengeCode)

	// get challenge response

	conn.SetReadDeadline(time.Now().Add(c.params.ChallengeTimeout))
	challengeResponse, err := helpers.Receive(conn)
	if err != nil {
		log.Printf("read challenge resp error: %s", err)
		return
	}
	log.Println("server got  challenge resp", challengeResponse)

	// verify challenge response

	ok, err := pow.Verify(challengeCode, challengeResponse)
	if err != nil {
		log.Printf("verify challenge resp error: %s", err)
		return
	}
	log.Println("server verified challenge resp:", ok)

	if !ok {
		return
	}

	// write fake sentence as response

	sentence, err := c.service.GetSentence(ctx)
	if err != nil {
		log.Printf("service get sentence error: %s", err)
		return
	}

	conn.SetWriteDeadline(time.Now().Add(c.params.WriteTimeout))
	_, err = conn.Write([]byte(sentence))
	if err != nil {
		log.Printf("write response error: %s", err)
		return
	}

	log.Println("server send sentence:", sentence)
}

func (s *Controller) getChallengeCode() []byte {
	challengeCode := uuid.New()
	return challengeCode[:]
}
