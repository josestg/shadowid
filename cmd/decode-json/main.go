package main

import (
	"encoding/json"
	"fmt"
	"github.com/josestg/shadowid"
)

func main() {
	shadowid.SetSalt(9602524670323041146)

	const raw = `{"id":"bbf4f504f8db4aa20495ed5692d111b51b0e6d4dc82e0b54"}`

	var target struct {
		ID shadowid.ShadowID `json:"id"`
	}

	if err := json.Unmarshal([]byte(raw), &target); err != nil {
		panic(err)
	}

	fmt.Println("ShadowID:", target.ID)
	fmt.Println("RandomID:", target.ID.RandomID())
	fmt.Println("Autoincr:", target.ID.Autoincr())
}
