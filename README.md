# procfs

[![Go Reference](https://pkg.go.dev/badge/go.jamescun.com/procfs.svg)](https://pkg.go.dev/go.jamescun.com/procfs)

```sh
$ go get go.jamescun.com/procfs
```

This package implements reading structured text-based information about a Linux (or other compatible) system through the procfs filesystem.

The source for this project is hosted on both [GitHub](https://github.com/jamescun/procfs) and [Codeberg](https://codeberg.org/jamescun/procfs).

## usage

This package is built with containers in mind, specifically the case where the procfs filesystem mounted at `/proc` is not the one you want to read, instead you want to read from the procfs of the host system bind mounted elsewhere inside a container.

Initializing to read from the normal `/proc` directory:

```go
proc := procfs.New()
```

Initializing to read from a different directory:

```go
proc := procfs.From(os.DirFS("/mnt/host/proc"))
```

It can also be used with anything implementing the [fs.FS](https://pkg.go.dev/io/fs) interface.

The procfs filesystem can be read directly as bytes or strings, additionally all the structures contained within this package implement the [encoding.TextUnmarshaler](https://pkg.go.dev/encoding#TextUnmarshaler) interface, which can be used with anything using that interface.

## examples

Below are some examples of how to use the `procfs` package.

> [!NOTE]
> Error handling is omitted in these examples for brevity.


### getting mount information

```go
package main

import (
	"fmt"

	"go.jamescun.com/procfs"
)

func main() {
	// initialize a procfs reader for the root /proc directory.
	proc := procfs.New()

	// read the mounted filesystems from /proc/self/mountinfo.
	mounts, _ := procfs.GetMounts(proc)

	for _, mount := range mounts {
		fmt.Printf("Root: %s, Path: %s, Type: %s\n", mount.Root, mount.Path, mount.Type)
	}
}
```


### getting network interface statistics

```go
package main

import (
	"fmt"

	"go.jamescun.com/procfs"
)

func main() {
	// initialize a procfs reader for the root /proc directory.
	proc := procfs.New()

	// get statistics from all the network links/interfaces on the system.
	stats, _ := procfs.GetLinkStats(proc)

	for _, link := range stats {
		fmt.Printf("Interface: %s\n", link.Name)
		fmt.Printf("  Receive: %d bytes, %d packets\n", link.Rx.Bytes, link.Rx.Packets)
		fmt.Printf("  Transmit: %d bytes, %d packets\n", link.Tx.Bytes, link.Tx.Packets)
	}
}
```
