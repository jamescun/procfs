// Copyright 2026 James Cunningham
// SPDX-License-Identifier: BSD-3-Clause
//
// Use of this source code is governed by a BSD-style license that can be found
// in the LICENSE file or at https://opensource.org/license/BSD-3-clause

package procfs

import (
	"io/fs"
	"os"
	"path/filepath"
)

// Error wraps errors returned by [Procfs] to contextualize the path within the
// procfs filesystem that was being read.
type Error struct {
	// Name is the basename of the file requested.
	Name string

	// Path is the relative path within the procfs filesystem attempted.
	Path string

	// Inner is the actual error, returned by [Error.Unwrap].
	Inner error
}

// wrap an error with the [Error] type, contextualized with path information,
// as well as unifying not found errors to the [fs.ErrNoExist] error.
func wrap(path string, inner error) error {
	if inner == nil {
		// given a nil error for some reason.
		return nil
	}

	if err, ok := inner.(*Error); ok {
		// already a wrapper error.
		return err
	}

	err := &Error{
		Name: filepath.Base(path),
		Path: path,
	}

	if os.IsNotExist(err) {
		err.Inner = fs.ErrNotExist
	} else if pe, ok := inner.(*fs.PathError); ok {
		err.Inner = pe.Err
	} else {
		err.Inner = inner
	}

	return err
}

func (e *Error) Error() string {
	return "procfs " + e.Name + " (" + e.Path + "): " + e.Inner.Error()
}

func (e *Error) Unwrap() error {
	return e.Inner
}
