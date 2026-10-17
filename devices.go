// Copyright 2026 James Cunningham
// SPDX-License-Identifier: BSD-3-Clause
//
// Use of this source code is governed by a BSD-style license that can be found
// in the LICENSE file or at https://opensource.org/license/BSD-3-clause

package procfs

import (
	"bytes"
	"strconv"
	"strings"

	"go.jamescun.com/procfs/internal/utils"
)

// InputDevice is one of the devices recognized by the input devices subsystem,
// read from /proc/bus/input/devices.
type InputDevice struct {
	Bus      string
	Vendor   string
	Product  string
	Version  string
	Name     string
	Phys     string
	Sysfs    string
	Uniq     string
	Handlers []string
}

// UnmarshalText unmarshals a single double-newline terminated entry from
// /proc/bus/input/devices.
func (i *InputDevice) UnmarshalText(b []byte) error {
	for _, line := range utils.Lines(b) {
		key, value, ok := utils.Split(string(line), func(b byte) bool {
			return b == ' '
		})
		if !ok {
			continue
		}

		if key == "I:" {
			for _, field := range utils.Fields(value) {
				key, value, ok := utils.Split(field, func(b byte) bool {
					return b == '='
				})
				if !ok {
					continue
				}

				switch key {
				case "Bus":
					i.Bus = value
				case "Vendor":
					i.Vendor = value
				case "Product":
					i.Product = value
				case "Version":
					i.Version = value
				}
			}

			continue
		}

		_, value, ok = utils.Split(value, func(b byte) bool {
			return b == '='
		})
		if !ok {
			continue
		}

		switch key {
		case "N:":
			i.Name, _ = strconv.Unquote(value)

		case "P:":
			i.Phys = value

		case "S:":
			i.Sysfs = value

		case "U:":
			i.Uniq = value

		case "H:":
			i.Handlers = strings.Split(value, " ")
		}
	}

	return nil
}

// InputDevices are the devices recognized by the input devices subsystem, read
// from /proc/bus/input/devices.
type InputDevices []*InputDevice

// GetInputDevices reads all the devices recognized by the input devices
// subsystem, read from /proc/bus/input/devices in the given [Procfs].
func GetInputDevices(proc Procfs) (InputDevices, error) {
	is := InputDevices{}

	err := proc.Read("bus/input/devices", &is)
	if err != nil {
		return nil, err
	}

	return is, nil
}

// UnmarshalText unmarshals all the double-newline terminated entries from
// /proc/bus/input/devices.
func (is *InputDevices) UnmarshalText(b []byte) error {
	var device []byte

	for {
		device, b, _ = bytes.Cut(b, []byte{'\n', '\n'})

		device = bytes.TrimSpace(device)
		if len(device) == 0 {
			break
		}

		i := new(InputDevice)
		err := i.UnmarshalText(device)
		if err != nil {
			return err
		}

		*is = append(*is, i)
	}

	return nil
}
