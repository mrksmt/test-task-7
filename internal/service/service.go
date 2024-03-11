package service

import (
	"context"
	"log"
	"net"
	"sync"

	"github.com/go-faker/faker/v4"
	"github.com/google/uuid"
	"github.com/pkg/errors"

	"github.com/mrksmt/test-task-7/pkg/helpers"
	"github.com/mrksmt/test-task-7/pkg/pow"
)

type Parameters struct {
	Address string `envconfig:"SENTENCE_SERVICE_ADDRESS" default:"localhost:8090"`
}

type SentenceService struct {
	address string
}

func NewSentenceService(
	params *Parameters,
) *SentenceService {
	s := &SentenceService{address: params.Address}
	return s
}

// Run implements slice.Dispatcher
func (s *SentenceService) Run(
	ctx context.Context,
) error {

	listenerWG := new(sync.WaitGroup)
	handlerWG := new(sync.WaitGroup)

	// listen for incoming connections
	listener, err := net.Listen("tcp", s.address)
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

func (s *SentenceService) handleClient(
	_ context.Context,
	wg *sync.WaitGroup,
	conn net.Conn,
) {
	defer wg.Done()
	defer conn.Close()

	log.Println("server established connection", conn.LocalAddr().String())

	// send challenge code

	challengeCode := s.getChallengeCode()

	_, err := conn.Write(challengeCode)
	if err != nil {
		log.Printf("write challenge code error: %s", err)
		return
	}
	log.Println("server send challenge code", challengeCode)

	// get challenge response

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

	sentence := faker.Sentence()
	_, err = conn.Write([]byte(sentence))
	if err != nil {
		log.Printf("write response error: %s", err)
		return
	}

	log.Println("server send sentence:", sentence)
}

func (s *SentenceService) getChallengeCode() []byte {
	challengeCode := uuid.New()
	return challengeCode[:]
}
