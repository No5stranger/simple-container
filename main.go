package main

import (
	"log"

	sc "github.com/no5stranger/simple-containerd/snapshot"
)

func main() {
	//err := sc.BusyBoxExample()
	err := sc.NginxExample()
	if err != nil {
		log.Fatal(err)
	}
}
