// Copyright 2026 James Cunningham
// SPDX-License-Identifier: BSD-3-Clause
//
// Use of this source code is governed by a BSD-style license that can be found
// in the LICENSE file or at https://opensource.org/license/BSD-3-clause

// Package procfs implements reading structured text-based information about a
// system through the procfs filesystem.
package procfs

import (
	"bytes"
	"encoding"
	"fmt"
	"io"
	"io/fs"
	"os"
	"unsafe"
)

// Procfs implements reading structured text-based information about a system
// through the procfs filesystem.
type Procfs interface {
	// List all the files and directories at a given path within the procfs
	// filesystem.
	List(path string) ([]fs.DirEntry, error)

	// Read a path from the procfs filesystem, and unmarshal it's contents to
	// a pointer to an object implementing the [encoding.TextUnmarshaler]
	// interface.
	Read(path string, dst encoding.TextUnmarshaler) error

	// ReadBytes reads a path from the procfs filesystem as bytes.
	//
	// Any leading or trailing whitespace will automatically be trimmed.
	ReadBytes(path string) ([]byte, error)

	// ReadString reads a path from the procfs filesystem as a string.
	//
	// Any leading or trailing whitespace will automatically be trimmed.
	ReadString(path string) (string, error)
}

type procfs struct {
	proc fs.FS
}

// New initializes a new [Procfs] where the procfs filesystem is mounted at
// /proc.
//
// To read from a procfs filesystem mounted elsewhere, such as inside a
// container and wanting to read the host procfs, see [From].
func New() Procfs {
	return From(os.DirFS("/proc"))
}

// From initializes a new [Procfs] with a given [fs.FS] pointing to the procfs
// filesystem.
//
// This is useful for situations where the target procfs to read from is not
// mounted at /proc, such as a bind mount within a container pointing to the
// host procfs.
func From(proc fs.FS) Procfs {
	return &procfs{proc: proc}
}

func (p *procfs) List(path string) ([]fs.DirEntry, error) {
	files, err := fs.ReadDir(p.proc, path)
	if err != nil {
		return nil, wrap(path, err)
	}

	return files, nil
}

func (p *procfs) Read(path string, dst encoding.TextUnmarshaler) error {
	if dst == nil {
		return fmt.Errorf("encoding.TextUnmarshaler is nil")
	}

	value, err := p.ReadBytes(path)
	if err != nil {
		return err
	}

	err = dst.UnmarshalText(value)
	if err != nil {
		return wrap(path, err)
	}

	return nil
}

func (p *procfs) ReadBytes(path string) ([]byte, error) {
	// procfs does not report the correct file size for files being read, often
	// either zero or the page size. so we start with a small buffer, fill it,
	// and keep expanding it until the whole file is read.

	// this should be a reasonable value for most files.
	const bufSize = 512

	file, err := p.proc.Open(path)
	if err != nil {
		return nil, wrap(path, err)
	}
	defer file.Close()

	var value []byte

	for {
		buf := make([]byte, bufSize)

		n, err := file.Read(buf)
		if err == io.EOF || n == 0 {
			// finished reading the file on the buffer size boundary.
			break
		} else if err != nil {
			return nil, wrap(path, err)
		}

		if value == nil {
			// first read, no need to append.
			value = buf[:n]
		} else {
			value = append(value, buf[:n]...)
		}

		if n < bufSize {
			// read less than the buffer size, final read.
			break
		}
	}

	return bytes.TrimSpace(value), nil
}

func (p *procfs) ReadString(path string) (string, error) {
	value, err := p.ReadBytes(path)
	if err != nil {
		return "", err
	}

	//nolint:gosec // safe, as the bytes are never returned outside this func.
	return unsafe.String(unsafe.SliceData(value), len(value)), nil
}

func (procfs) String() string {
	return "Procfs(path: /proc)"
}
