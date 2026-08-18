// Copyright 2026 James Cunningham
// SPDX-License-Identifier: BSD-3-Clause
//
// Use of this source code is governed by a BSD-style license that can be found
// in the LICENSE file or at https://opensource.org/license/BSD-3-clause

package sysctl

import (
	"os"
	"testing"

	"go.jamescun.com/procfs"
)

func TestSysctl(t *testing.T) {
	sys := From(procfs.From(os.DirFS("testdata/proc")))

	t.Run("ReadBool", func(t *testing.T) {
		expected := true

		target, err := sys.ReadBool("net.ipv4.ip_forward")
		if err != nil {
			t.Fatal("unexpected error:", err)
		}

		if expected != target {
			t.Errorf("expected %t got %t", expected, target)
		}
	})

	t.Run("ReadInt64", func(t *testing.T) {
		expected := int64(-1)

		target, err := sys.ReadInt64("kernel.msg_next_id")
		if err != nil {
			t.Fatal("unexpected error:", err)
		}

		if expected != target {
			t.Errorf("expected %d got %d", expected, target)
		}
	})

	t.Run("ReadInt64/Invalid", func(t *testing.T) {
		_, err := sys.ReadInt64("kernel.osrelease")
		if err == nil {
			t.Fatal("expected error")
		}
	})

	t.Run("ReadString", func(t *testing.T) {
		expected := "6.12.0-124.35.1.el10_1.x86_64"

		target, err := sys.ReadString("kernel.osrelease")
		if err != nil {
			t.Fatal("unexpected error:", err)
		}

		if expected != target {
			t.Errorf("expected:\n%s\ngot:\n%s", expected, target)
		}
	})

	t.Run("ReadUint64", func(t *testing.T) {
		expected := uint64(262144)

		target, err := sys.ReadUint64("net.nf_conntrack_max")
		if err != nil {
			t.Fatal("unexpected error:", err)
		}

		if expected != target {
			t.Errorf("expected %d got %d", expected, target)
		}
	})

	t.Run("ReadUint64/Invalid", func(t *testing.T) {
		_, err := sys.ReadUint64("kernel.osrelease")
		if err == nil {
			t.Fatal("expected error")
		}
	})
}

func TestToPath(t *testing.T) {
	tests := []struct {
		name string
		path string
	}{
		{"", "sys/"},
		{"hello", "sys/hello"},
		{"hello.world", "sys/hello/world"},
		{"hello/world", "sys/hello/world"},
	}

	for _, test := range tests {
		path := toPath(test.name)

		if test.path != path {
			t.Errorf("expected:\n%s\ngot:\n%s", test.path, test.name)
		}
	}
}
