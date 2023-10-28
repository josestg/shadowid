package shadowid

import (
	"github.com/google/uuid"
	"testing"
)

func init() {
	SetSalt(9602524670323041146)
}

const (
	autoincr        = int64(237502)
	shadowIDExample = "bbf4f504f8db4aa20495ed5692d111b51b0e6d4dc82e0b54"
)

var randomid = uuid.MustParse("bbf4f504-f8db-4aa2-92d1-11b51b0e6d4d")

func TestNewShadowID(t *testing.T) {
	id := NewShadowID(autoincr, randomid)
	if id.String() != shadowIDExample {
		t.Error("invalid shadow id")
	}

	if id.RandomID() != randomid {
		t.Error("invalid random id")
	}

	if id.Autoincr() != autoincr {
		t.Error("invalid autoincr")
	}
}

func TestParse(t *testing.T) {
	id, err := Parse([]byte(shadowIDExample))
	if err != nil {
		t.Error(err)
	}

	if id.String() != shadowIDExample {
		t.Error("invalid shadow id")
	}

	if id.RandomID() != randomid {
		t.Error("invalid random id")
	}

	if id.Autoincr() != autoincr {
		t.Error("invalid autoincr")
	}
}

func TestShadowID_UnmarshalText(t *testing.T) {
	var id ShadowID

	if err := id.UnmarshalText([]byte(shadowIDExample)); err != nil {
		t.Error(err)
	}

	if id.String() != shadowIDExample {
		t.Error("invalid shadow id")
	}
	if id.RandomID() != randomid {
		t.Error("invalid random id")
	}

	if id.Autoincr() != autoincr {
		t.Error("invalid autoincr")
	}
}

func TestShadowID_MarshalText(t *testing.T) {

	id := NewShadowID(autoincr, uuid.MustParse("bbf4f504-f8db-4aa2-92d1-11b51b0e6d4d"))
	text, err := id.MarshalText()
	if err != nil {
		t.Error(err)
	}

	if string(text) != shadowIDExample {
		t.Error("invalid shadow id")
	}
}
