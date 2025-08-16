package shadowid

import (
	"math/rand"
	"reflect"
	"testing"
	"testing/quick"

	"github.com/google/uuid"
)

type Pair struct {
	UUID uuid.UUID
	Incr int64
}

func (p Pair) Generate(rand *rand.Rand, _ int) reflect.Value {
	uid, err := uuid.NewRandomFromReader(rand)
	if err != nil {
		panic(err)
	}

	n := rand.Int63()
	return reflect.ValueOf(Pair{Incr: n, UUID: uid})
}

func TestPropEncoding(t *testing.T) {
	f := func(p Pair) bool {
		x := NewShadowID(p.Incr, p.UUID)
		xt, err := x.MarshalText()
		if err != nil {
			panic(err)
		}
		var y ShadowID
		err = y.UnmarshalText(xt)
		if err != nil {
			panic(err)
		}
		return x == y && y.RandomID() == p.UUID && y.Autoincr() == p.Incr
	}
	err := quick.Check(f, nil)
	if err != nil {
		t.Error(err)
	}
}
