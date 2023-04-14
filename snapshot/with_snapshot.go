package snapshot

import (
	"context"
	"fmt"
	"github.com/containerd/containerd/snapshots"
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
	labels := map[string]string{"containerd.io/snapshot/image-cache": "imc-busybox-latest"}
	image, err := client.Pull(
		ctx,
		"docker.io/library/busybox:latest",
		containerd.WithPullUnpack,
		containerd.WithPullSnapshotter("a-overlayfs", snapshots.WithLabels(labels)),
	)
	if err != nil {
		return err
	}
	log.Printf("Pulled image, name:%s", image.Name())
	rootfs, err := image.RootFS(ctx)
	if err != nil {
		return err
	}
	pullCostTime := time.Now().UnixMilli() - startTime.UnixMilli()
	log.Printf("Image Pulled, Rootfs:%v, cost:%dms", rootfs, pullCostTime)

	container, err := client.NewContainer(
		ctx,
		"busybox-server",
		containerd.WithSnapshotter("a-overlayfs"),
		containerd.WithNewSnapshot("busybox-snapshot", image),
		containerd.WithNewSpec(oci.WithImageConfig(image)),
	)
	if err != nil {
		return err
	}
	defer container.Delete(ctx, containerd.WithSnapshotCleanup)
	containerCostTime := time.Now().UnixMilli() - startTime.UnixMilli()
	log.Printf("Containerd created, id:%s, cost:%dms", container.ID(), containerCostTime)

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

func NginxExample() error {
	startTime := time.Now()
	client, err := containerd.New("/run/containerd/containerd.sock")
	if err != nil {
		return err
	}
	ctx := namespaces.WithNamespace(context.Background(), "default")
	labels := map[string]string{"containerd.io/snapshot/image-cache": "imc-nginx-stable"}
	image, err := client.Pull(
		ctx,
		"docker.io/library/nginx:stable",
		containerd.WithPullUnpack,
		containerd.WithPullSnapshotter("a-overlayfs", snapshots.WithLabels(labels)),
	)
	if err != nil {
		return err
	}
	log.Printf("Pulled image, name:%s", image.Name())
	rootfs, err := image.RootFS(ctx)
	if err != nil {
		return err
	}
	pullCostTime := time.Now().UnixMilli() - startTime.UnixMilli()
	log.Printf("Image Pulled, Rootfs:%v, cost:%dms", rootfs, pullCostTime)

	container, err := client.NewContainer(
		ctx,
		"nginx-server",
		containerd.WithSnapshotter("a-overlayfs"),
		containerd.WithNewSnapshot("nginx-snapshot", image),
		containerd.WithNewSpec(oci.WithImageConfig(image)),
	)
	if err != nil {
		return err
	}
	defer container.Delete(ctx, containerd.WithSnapshotCleanup)
	containerCostTime := time.Now().UnixMilli() - startTime.UnixMilli()
	log.Printf("Containerd created, id:%s, cost:%dms", container.ID(), containerCostTime)

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
