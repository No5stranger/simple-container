package snapshot

import (
	"context"
	"fmt"
	"log"
	"strings"
	"syscall"
	"time"

	"github.com/containerd/containerd/snapshots"

	"github.com/containerd/containerd"
	"github.com/containerd/containerd/cio"
	"github.com/containerd/containerd/namespaces"
	"github.com/containerd/containerd/oci"
)

func ContainerExample(ref string, waiTime int64) error {
	startTime := time.Now()
	client, err := containerd.New("/run/containerd/containerd.sock")
	if err != nil {
		return err
	}
	ctx := namespaces.WithNamespace(context.Background(), "default")
	refName := strings.ReplaceAll(ref, ":", "-")
	labels := map[string]string{"containerd.io/snapshot/image-cache": fmt.Sprintf("imc-%s", refName)}
	image, err := client.Pull(
		ctx,
		fmt.Sprintf("docker.io/library/%s", ref),
		containerd.WithPullUnpack,
		containerd.WithPullSnapshotter("a-overlayfs", snapshots.WithLabels(labels)),
		//containerd.WithPullSnapshotter("a-overlayfs"),
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
		fmt.Sprintf("%s-server", refName),
		containerd.WithSnapshotter("a-overlayfs"),
		containerd.WithNewSnapshot(fmt.Sprintf("%s-snapshot", refName), image),
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

	time.Sleep(time.Duration(waiTime) * time.Second)
	if err := task.Kill(ctx, syscall.SIGKILL); err != nil {
		return err
	}
	fmt.Printf("Task killed, %v\n", task.ID())

	status := <-existStatus
	code, _, err := status.Result()
	if err != nil {
		return err
	}
	log.Printf("Task exited, code:%d", code)
	return nil
}

func WithSnapshot(id string, waiTime int64) error {
	starTime := time.Now()
	client, err := containerd.New("/run/containerd/containerd.sock")
	if err != nil {
		return err
	}
	ctx := namespaces.WithNamespace(context.Background(), "default")
	container, err := client.NewContainer(
		ctx,
		"test-custom-snapshot",
		containerd.WithSnapshot(id),
		containerd.WithNewSpec(oci.WithDefaultSpec()),
	)
	if err != nil {
		return err
	}
	defer container.Delete(ctx)
	cosTime := time.Now().UnixMilli() - starTime.UnixMilli()
	log.Printf("Containerd created, id:%s, cost:%sms", container.ID(), cosTime)

	task, err := container.NewTask(ctx, cio.NewCreator(cio.WithStdio))
	if err != nil {
		return err
	}
	log.Printf("Task Created, %s", task.ID())
	defer task.Delete(ctx)
	existStatus, err := task.Wait(ctx)
	if err != nil {
		return err
	}
	if err := task.Start(ctx); err != nil {
		return err
	}
	log.Printf("containerd started...")

	time.Sleep(time.Duration(waiTime) * time.Second)
	if err := task.Kill(ctx, syscall.SIGKILL); err != nil {
		log.Printf("Task killed, id:%s", task.ID())
	}
	status := <-existStatus
	code, _, err := status.Result()
	if err != nil {
		return err
	}
	log.Printf("Task existed, code:%d", code)
	return nil
}
