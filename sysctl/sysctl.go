// Copyright 2026 James Cunningham
// SPDX-License-Identifier: BSD-3-Clause
//
// Use of this source code is governed by a BSD-style license that can be found
// in the LICENSE file or at https://opensource.org/license/BSD-3-clause

// Package sysctl implements reading the sysctl Linux Kernel parameters
// exposed through the procfs filesystem.
//
// This interface does not expose writing to them.
package sysctl

import (
	"fmt"
	"strings"

	"go.jamescun.com/procfs"
	"go.jamescun.com/procfs/internal/utils"
)

// Sysctl implements reading the sysctl Linux Kernel parameters exposed through
// the procfs filesystem.
//
// Names will automatically be converted from dot-notation to the equivalent
// path name.
type Sysctl interface {
	// ReadBool reads the named boolean parameter from sysctl.
	//
	// Boolean parameters should always be either 0 or 1, however this method
	// will tolerate other truthy values.
	ReadBool(name string) (bool, error)

	// ReadInt64 reads the named int64 parameter from sysctl.
	ReadInt64(name string) (int64, error)

	// ReadString reads the named string parameter from sysctl.
	//
	// Any leading or trailing whitespace will be trimmed.
	ReadString(name string) (string, error)

	// ReadUint64 reads the named uint64 parameter from sysctl.
	ReadUint64(name string) (uint64, error)
}

type sysctl struct {
	proc procfs.Procfs
}

// New initializes a new [Sysctl] where the procfs filesystem is mounted at
// /proc.
//
// To read from sysctl where the procfs filesystem mounted elsewhere, such as
// inside a container and wanting to read the host procfs, see [From].
func New() Sysctl {
	return From(procfs.New())
}

// From initializes a new [Sysctl] with a given [procfs.Procfs] pointing to the
// procfs filesystem.
//
// This is useful for situations where the target procfs to read from is not
// mounted at /proc, such as a bind mount within a container pointing to the
// host procfs.
func From(proc procfs.Procfs) Sysctl {
	return &sysctl{proc: proc}
}

func (s *sysctl) ReadBool(name string) (bool, error) {
	value, err := s.proc.ReadBytes(toPath(name))
	if err != nil {
		return false, err
	}

	if len(value) == 0 {
		return false, nil
	}

	switch value[0] {
	case '1', 't', 'T', 'y', 'Y':
		return true, nil

	default:
		return false, nil
	}
}

func (s *sysctl) ReadInt64(name string) (int64, error) {
	value, err := s.proc.ReadBytes(toPath(name))
	if err != nil {
		return 0, err
	}

	n, ok := utils.Int[int64](value)
	if !ok {
		return 0, fmt.Errorf("invalid number %q", value)
	}

	return n, nil
}

func (s *sysctl) ReadString(name string) (string, error) {
	return s.proc.ReadString(toPath(name))
}

func (s *sysctl) ReadUint64(name string) (uint64, error) {
	value, err := s.proc.ReadBytes(toPath(name))
	if err != nil {
		return 0, err
	}

	n, ok := utils.Uint[uint64](value)
	if !ok {
		return 0, fmt.Errorf("invalid number %q", value)
	}

	return n, nil
}

// toPath converts a dot-notation sysctl parameter name to it's equivalent
// path name, including the `sys/` prefix. It does not clean the path.
func toPath(name string) string {
	var s strings.Builder

	s.Grow(len(name) + 4)

	s.WriteString("sys/")

	for i := range name {
		if name[i] == '.' {
			s.WriteByte('/')
		} else {
			s.WriteByte(name[i])
		}
	}

	return s.String()
}
