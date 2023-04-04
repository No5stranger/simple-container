package snapshot

import (
	"context"
	"log"

	"github.com/containerd/containerd"
	"github.com/containerd/containerd/cio"
	"github.com/containerd/containerd/namespaces"
	"github.com/containerd/containerd/oci"
)

func BusyBoxExample() error {
	client, err := containerd.New("/run/containerd/containerd.sock")
	if err != nil {
		return err
	}
	ctx := namespaces.WithNamespace(context.Background(), "default")
	image, err := client.Pull(
		ctx,
		"docker.io/library/busybox:latest",
		containerd.WithPullUnpack,
		containerd.WithPullSnapshotter("a-overlayfs"),
	)
	if err != nil {
		return err
	}
	log.Printf("Pulled image %s", image.Name())
	rootfs, err := image.RootFS(ctx)
	if err != nil {
		return err
	}
	log.Printf("Image Rootfs %v", rootfs)

	container, err := client.NewContainer(
		ctx,
		"busy-server",
		containerd.WithNewSnapshot("busybox-snapshot", image),
		containerd.WithNewSpec(oci.WithImageConfig(image)),
	)
	if err != nil {
		return err
	}
	//defer container.Delete(ctx, container.WithSnapshotCleanup)
	log.Printf("Containerd created %s", container.ID())

	task, err := container.NewTask(ctx, cio.NewCreator(cio.WithStdio))
	if err != nil {
		return err
	}
	if err := task.Start(ctx); err != nil {
		return err
	}
	log.Printf("Containerd started...")
	return nil
}
