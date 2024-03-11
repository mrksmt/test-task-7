package pow

import (
	"bytes"
	"context"
	"encoding/binary"
	"fmt"
	"math"
	"math/rand"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/require"
)

func Test_GetProof(t *testing.T) {

	code := uuid.New()

	fmt.Println(code)

	proof, err := GetProof(context.TODO(), code[:])
	require.NoError(t, err)
	require.NotEmpty(t, proof)
}

func Test_Verify(t *testing.T) {

	code := uuid.New()

	proof, err := GetProof(context.TODO(), code[:])
	require.NoError(t, err)
	require.NotEmpty(t, proof)

	// positive
	// valid proof should be positive verified
	{
		ok, err := Verify(code[:], proof)
		require.NoError(t, err)
		require.True(t, ok)
	}

	// negative
	// random proof should be negative verified in most cases
	total := 10000
	positive := 0

	for cnt := 0; cnt < 10000; cnt++ {

		rnd := rand.Int31n(math.MaxInt16)

		bs := make([]byte, 4)
		binary.LittleEndian.PutUint16(bs, uint16(rnd))

		if bytes.Equal(proof, bs) {
			continue
		}

		ok, err := Verify(code[:], bs)
		require.NoError(t, err)

		if ok {
			positive++
		}
	}

	require.Less(t, positive, total/100)
}

func Benchmark_GetProof(b *testing.B) {
	difficult = 2
	for i := 0; i < b.N; i++ {
		challengeCode := uuid.New()
		_, err := GetProof(context.TODO(), challengeCode[:])
		require.NoError(b, err)
	}
}

func Benchmark_GetProofAsync(b *testing.B) {
	difficult = 2
	workerPoolSize = 4
	for i := 0; i < b.N; i++ {
		challengeCode := uuid.New()
		_, err := GetProofAsync(context.TODO(), challengeCode[:])
		require.NoError(b, err)
	}
}
