package main

import (
	"flag"
	sc "github.com/no5stranger/simple-containerd/snapshot"
	"log"
)

func main() {
	var (
		err     error
		ref     string
		snapID  string
		waiTime int64
		action  string
	)
	flag.StringVar(&ref, "ref", "busybox:latest", "image name")
	flag.StringVar(&snapID, "snapid", "", "snapshot id")
	flag.Int64Var(&waiTime, "wait", 60, "wait time to kill container")
	flag.StringVar(&action, "action", "image", "action to exec, image:with image, snap:with snapshot")
	flag.Parse()

	if waiTime <= 0 {
		waiTime = 60
	}

	switch action {
	case "image":
		if len(ref) == 0 {
			log.Fatal("miss image ref")
		}
		err = sc.ContainerExample(ref, waiTime)
	case "snap":
		if len(snapID) == 0 {
			log.Fatal("miss snapshot id")
		}
		err = sc.WithSnapshot(snapID, waiTime)
	}
	if err != nil {
		log.Fatal(err)
	}
	log.Printf("done...")
}
