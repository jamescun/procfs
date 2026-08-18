// Copyright 2026 James Cunningham
// SPDX-License-Identifier: BSD-3-Clause
//
// Use of this source code is governed by a BSD-style license that can be found
// in the LICENSE file or at https://opensource.org/license/BSD-3-clause

// Package utils contains a set of utilities that make it easier to deal with
// files within the procfs filesystem, but don't make sense to expose as the
// general purpose API.
package utils

import (
	"iter"
)

// Equal returns true if the given string or []byte are equal.
func Equal[A, B ~string | ~[]byte](a A, b B) bool {
	if len(a) != len(b) {
		// cannot be equal.
		return false
	}

	for i := range len(a) {
		if a[i] != b[i] {
			return false
		}
	}

	return true
}

// Fields returns an [iter.Seq2] iterator that will yield for each whitespace
// delimited field in the given string or []byte body, as well as the logical
// field number.
func Fields[T ~string | ~[]byte](body T) iter.Seq2[int, T] {
	return func(yield func(int, T) bool) {
		var field T

		i := 0
		for len(body) > 0 {
			body = TrimLeft(body)

			field, body, _ = Split(body, isSpace)
			if len(field) == 0 {
				// empty field, end of body.
				break
			}

			if !yield(i, field) {
				break
			}

			i++
		}
	}
}

// Lines returns an [iter.Seq2] iterator that will yield for each non-empty
// line in the given string or []byte body, as well as the logical line number.
func Lines[T ~string | ~[]byte](body T) iter.Seq2[int, T] {
	return func(yield func(int, T) bool) {
		var line T

		i := 0
		for len(body) > 0 {
			line, body, _ = Split(body, isLine)
			if len(line) == 0 {
				// empty line, nothing left to yield.
				continue
			}

			if !yield(i, line) {
				break
			}

			i++
		}
	}
}

// HasPrefix returns true is a string or []byte body has the given prefix.
func HasPrefix[A, B ~string | ~[]byte](body A, prefix B) bool {
	if len(body) < len(prefix) {
		// cannot have prefix.
		return false
	}

	for i := range len(prefix) {
		if body[i] != prefix[i] {
			return false
		}
	}

	return true
}

// IndexByte returns the index in string or []byte of a byte matched by the
// given function. If not found, -1 is returned.
func IndexByte[T ~string | ~[]byte](body T, match func(byte) bool) int {
	if len(body) == 0 {
		// nothing to index.
		return -1
	}

	for i := range len(body) {
		if match(body[i]) {
			return i
		}
	}

	return -1
}

// Int parses a string or []byte as a signed integer from an ASCII-encoded
// body. It does not tolerate whitespace or any non-numerical characters.
func Int[E ~int8 | ~int16 | ~int32 | ~int64 | ~int, T ~string | ~[]byte](body T) (value E, ok bool) {
	if len(body) == 0 {
		// no number to read.
		return
	}

	// strip negative numbers, negate after parsing.
	neg := body[0] == '-'
	if neg {
		body = body[1:]
	}

	for i := range len(body) {
		if body[i] < '0' || body[i] > '9' {
			// not a number.
			return
		}

		value = value*10 + E(body[i]-'0')
	}

	if neg {
		value = -value
	}

	ok = true
	return
}

// Split a string or []byte by a delimiting byte function, returning before and
// after the split, and if a split took place.
//
// If no split takes place, the whole body will be returned as before.
func Split[T ~string | ~[]byte](body T, delim func(byte) bool) (before, after T, split bool) {
	if len(body) == 0 {
		// no split can take place.
		return
	}

	// ensure before is returned when there is no split.
	before = body

	for i := range len(before) {
		if delim(before[i]) {
			after = before[i+1:]
			before = before[:i]
			split = true
			return
		}
	}

	return
}

// TrimSpace trims any leading or trailing whitespace from a string or []byte.
// An empty value is returned if the body only contains whitespace.
func TrimSpace[T ~string | ~[]byte](body T) T {
	return TrimLeft(TrimRight(body))
}

// TrimLeft trims any leading whitespace from a string or []byte. An empty
// value is returned if the body only contains whitespace.
func TrimLeft[T ~string | ~[]byte](body T) (value T) {
	if len(body) == 0 {
		// nothing to trim.
		return
	}

	for i := range len(body) {
		if notSpace(body[i]) {
			value = body[i:]
			return
		}
	}

	return
}

// TrimRight trims any trailing whitespace from a string or []byte. An empty
// value is returned if the body only contains whitespace.
func TrimRight[T ~string | ~[]byte](body T) (value T) {
	if len(body) == 0 {
		// nothing to trim
		return
	}

	for i := len(body) - 1; i >= 0; i-- {
		if notSpace(body[i]) {
			value = body[:i+1]
			return
		}
	}

	return
}

// isLine returns true if byte b is a newline.
func isLine(b byte) bool {
	return b == '\n'
}

// isSpace returns true if byte b is ASCII whitespace.
func isSpace(b byte) bool {
	switch b {
	case '\t', '\n', '\v', '\f', '\r', ' ', 0x85, 0xA0: // NEL, NBSP
		return true

	default:
		return false
	}
}

// notSpace returns true if byte b is not ASCII whitespace.
func notSpace(b byte) bool {
	return !isSpace(b)
}

// Uint parses a string or []byte as an unsigned integer from an ASCII-encoded
// body. It does not tolerate whitespace or any non-numerical characters.
func Uint[E ~uint8 | ~uint16 | ~uint32 | ~uint64 | ~uint, T ~string | ~[]byte](body T) (value E, ok bool) {
	if len(body) == 0 {
		// no number to read.
		return
	}

	for i := range len(body) {
		if body[i] < '0' || body[i] > '9' {
			// not a number.
			return
		}

		value = value*10 + E(body[i]-'0')
	}

	ok = true
	return
}
