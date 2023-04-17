package main

import (
	"flag"
	sc "github.com/no5stranger/simple-containerd/snapshot"
	"log"
)

func main() {
	var (
		ref     string
		waiTime int64
	)
	flag.StringVar(&ref, "ref", "busybox:latest", "image name")
	flag.Int64Var(&waiTime, "wait", 60, "wait time to kill container")
	if len(ref) == 0 {
		log.Fatal("miss image ref")
	}
	if waiTime <= 0 {
		waiTime = 60
	}
	err := sc.ContainerExample(ref, waiTime)
	if err != nil {
		log.Fatal(err)
	}
}
