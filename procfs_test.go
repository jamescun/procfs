// Copyright 2026 James Cunningham
// SPDX-License-Identifier: BSD-3-Clause
//
// Use of this source code is governed by a BSD-style license that can be found
// in the LICENSE file or at https://opensource.org/license/BSD-3-clause

package procfs

import (
	"errors"
	"io/fs"
	"os"
	"testing"

	"github.com/google/go-cmp/cmp"
)

// String implements [encoding.TextUnmarshaler] for testing.
type String string

func (s *String) UnmarshalText(b []byte) error {
	*s = String(b)
	return nil
}

func TestProcfs(t *testing.T) {
	proc := From(os.DirFS("testdata/proc"))

	t.Run("Read", func(t *testing.T) {
		expected := String("Linux version 6.12.0-124.35.1.el10_1.x86_64 (mockbuild@x64-builder02.almalinux.org)")

		var target String

		err := proc.Read("version", &target)
		if err != nil {
			t.Fatal("unexpected error:", err)
		}

		if !cmp.Equal(expected, target) {
			t.Error(cmp.Diff(expected, target))
		}
	})

	t.Run("ReadString", func(t *testing.T) {
		// tests both ReadBytes and ReadString with current implementation.

		expected := "Linux version 6.12.0-124.35.1.el10_1.x86_64 (mockbuild@x64-builder02.almalinux.org)"

		target, err := proc.ReadString("version")
		if err != nil {
			t.Fatal("unexpected error:", err)
		}

		if !cmp.Equal(expected, target) {
			t.Error(cmp.Diff(expected, target))
		}
	})

	t.Run("ReadString/NotFound", func(t *testing.T) {
		_, err := proc.ReadString("not_found")
		if err == nil {
			t.Fatal("expected error")
		}

		if !errors.Is(err, fs.ErrNotExist) {
			t.Errorf("expected error:\n%v\ngot:\n%v", fs.ErrInvalid, err)
		}
	})
}
