package main

import (
	"log"
	"os"
)

func main() {
	// run
	if err := WriteCurrentTime(os.Stdout); err != nil {
		log.Fatalln(err)
	}

	os.Exit(0)
}
