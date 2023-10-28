package main

import (
	"fmt"
	"github.com/google/uuid"
	"github.com/josestg/shadowid"
)

func main() {
	shadowid.SetSalt(9602524670323041146)

	autoincr := int64(237502)
	randomid, _ := uuid.Parse("bbf4f504f8db4aa292d111b51b0e6d4d")
	id := shadowid.NewShadowID(autoincr, randomid)

	fmt.Println(id)
}
