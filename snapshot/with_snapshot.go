package snapshot

import (
	"context"
	"fmt"
	"log"
	"syscall"
	"time"

	"github.com/containerd/containerd"
	"github.com/containerd/containerd/cio"
	"github.com/containerd/containerd/namespaces"
	"github.com/containerd/containerd/oci"
)

func BusyBoxExample() error {
	startTime := time.Now()
	client, err := containerd.New("/run/containerd/containerd.sock")
	if err != nil {
		return err
	}
	ctx := namespaces.WithNamespace(context.Background(), "default")
	image, err := client.Pull(
		ctx,
		"docker.io/library/busybox:latest",
		containerd.WithPullUnpack,
		//containerd.WithPullSnapshotter("a-overlayfs"),
	)
	if err != nil {
		return err
	}
	log.Printf("Pulled image %s", image.Name())
	rootfs, err := image.RootFS(ctx)
	if err != nil {
		return err
	}
	pullCostTime := time.Now().Unix() - startTime.Unix()
	log.Printf("Image Pulled, Rootfs:%v, cost:%d", rootfs, pullCostTime)

	container, err := client.NewContainer(
		ctx,
		"busybox-server",
		containerd.WithNewSnapshot("busybox-snapshot", image),
		containerd.WithNewSpec(oci.WithImageConfig(image)),
	)
	if err != nil {
		return err
	}
	defer container.Delete(ctx, containerd.WithSnapshotCleanup)
	containerCostTime := time.Now().Unix() - startTime.Unix()
	log.Printf("Containerd created %s, cost:%d", container.ID(), containerCostTime)

	task, err := container.NewTask(ctx, cio.NewCreator(cio.WithStdio))
	if err != nil {
		return err
	}
	log.Printf("Task created, %s", task.ID())
	defer task.Delete(ctx)
	existStatus, err := task.Wait(ctx)
	if err != nil {
		return err
	}
	if err := task.Start(ctx); err != nil {
		return err
	}
	log.Printf("Containerd started...")

	time.Sleep(3 * time.Second)
	if err := task.Kill(ctx, syscall.SIGKILL); err != nil {
		return err
	}
	fmt.Printf("Task killed, %v", task.ID())

	status := <-existStatus
	code, _, err := status.Result()
	if err != nil {
		return err
	}
	log.Printf("Task exited, code:%d", code)
	return nil
}
