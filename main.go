package main

import (
	"log"
	"os"

	sc "github.com/no5stranger/simple-containerd/snapshot"
)

func main() {
	if len(os.Args) < 2 {
		log.Fatal("miss image ref")
	}
	err := sc.ContainerExample(os.Args[1])
	if err != nil {
		log.Fatal(err)
	}
}
