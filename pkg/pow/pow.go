package pow

import (
	"context"
	"crypto/sha256"
	"encoding/binary"
	"errors"
	"log"
	"math"
	"runtime"
	"sync"
)

var difficult = 3 // number of zero bytes hash started
var rangeSize = 1024 * 16

var workerPoolSize = runtime.NumCPU() / 2

func GetProof(
	ctx context.Context,
	challengeCode []byte,
) ([]byte, error) {

	bs := make([]byte, 8)

	for cnt := range math.MaxUint32 {

		select {
		case <-ctx.Done():
			return nil, ctx.Err()
		default:
		}

		ok, err := checkCount(challengeCode, cnt, bs)
		if err != nil {
			return nil, err
		}

		if !ok {
			continue
		}

		binary.LittleEndian.PutUint64(bs, uint64(cnt))

		return bs, nil
	}

	return nil, errors.New("proof not founded")
}

func Verify(
	challengeCode []byte,
	proof []byte,
) (bool, error) {

	h := sha256.New()
	h.Write(challengeCode)
	h.Write(proof)

	result := h.Sum(nil)

	if len(result) < difficult {
		return false, errors.New("wrong hash size")
	}

	for idx := range difficult {
		if result[idx] != 0 {
			return false, nil
		}
	}

	return true, nil
}

func GetProofAsync(
	ctx context.Context,
	challengeCode []byte,
) ([]byte, error) {

	ctx, cancel := context.WithCancel(ctx)
	defer cancel()
	wg := new(sync.WaitGroup)

	in := make(chan int)
	out := make(chan int, 1)

	// run workers

	for range workerPoolSize {

		wg.Add(1)
		go func() {
			defer wg.Done()

			bs := make([]byte, 8)

			for cnt := range in {
				for c2 := range rangeSize {
					value := cnt + c2
					ok, err := checkCount(challengeCode, value, bs)
					if err != nil {
						log.Println("check count error:", err)
						continue
					}

					if !ok {
						continue
					}

					select {
					case out <- value:
						cancel()
					default:
					}
				}
			}
		}()
	}

	// run producer

	wg.Add(1)
	go func() {
		defer wg.Done()
		for cnt := 0; cnt < math.MaxUint32; cnt += rangeSize {
			select {
			case <-ctx.Done():
				close(in)
				return
			case in <- cnt:
			}
		}
	}()

	// wait for stop and get result

	wg.Wait()

	select {
	case result := <-out:
		bs := make([]byte, 8)
		binary.LittleEndian.PutUint64(bs, uint64(result))
		return bs, nil
	default:
		return nil, errors.New("proof not founded")
	}

}

func checkCount(
	challengeCode []byte,
	cnt int,
	bs []byte,
) (bool, error) {

	binary.LittleEndian.PutUint64(bs, uint64(cnt))

	h := sha256.New()
	h.Write(challengeCode)
	h.Write(bs)
	result := h.Sum(nil)

	if len(result) < difficult {
		return false, errors.New("wrong hash size")
	}

	for idx := range difficult {
		if result[idx] != 0 {
			return false, nil
		}
	}

	return true, nil
}
