package ramble

import (
	"log"

	"github.com/esote/util/splay"
)

func init() {
	var err error

	pub, err = splay.NewSplay("publickeys", 2)

	if err != nil {
		log.Fatal(err)
	}
}

var pub *splay.Splay
