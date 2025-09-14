package main

import (
	"log"

	"github.com/aaaxpel/album/internal/cmd"
)

func main() {
	log.SetFlags(log.Ldate | log.Ltime | log.Lshortfile)

	cmd.Router()
}
