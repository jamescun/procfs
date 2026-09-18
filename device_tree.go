// Copyright 2026 James Cunningham
// SPDX-License-Identifier: BSD-3-Clause
//
// Use of this source code is governed by a BSD-style license that can be found
// in the LICENSE file or at https://opensource.org/license/BSD-3-clause

package procfs

import (
	"fmt"
	"strings"
)

// DeviceTree contains information about a system that uses a device tree to
// configure devices for the Linux Kernel, read from /proc/device-tree.
//
// Not all values will be populated for every system depending on their
// configuration.
type DeviceTree struct {
	// Name of the device, read from /proc/device-tree/name.
	Name string

	// Model of the device, read from /proc/device-tree/model.
	Model string

	// Serial number of the device.
	Serial string

	// Compatible contains additional system definitions this system is
	// compatible with, read from /proc/device-tree/compatible.
	Compatible []string

	// BootArgs are the command line arguments given to the Linux Kernel by the
	// bootloader, read from /proc/device-tree/chosen/bootargs.
	BootArgs []string
}

// GetDeviceTree reads information about the systems device tree, from the
// given [Procfs].
//
// An error will be returned if the system does not use a device tree, or
// nothing could be read from it.
func GetDeviceTree(proc Procfs) (*DeviceTree, error) {
	ok := false
	dt := &DeviceTree{}

	compat, err := proc.ReadString("device-tree/compatible")
	if err == nil {
		ok = true
		dt.Compatible = strings.Split(compat, ",")
	}

	bootArgs, err := proc.ReadString("device-tree/chosen/bootargs")
	if err == nil {
		ok = true
		dt.BootArgs = strings.Fields(bootArgs)
	}

	dt.Name, err = proc.ReadString("device-tree/name")
	if err == nil {
		ok = true
	}

	dt.Model, err = proc.ReadString("device-tree/model")
	if err == nil {
		ok = true
	}

	dt.Serial, err = proc.ReadString("device-tree/serial")
	if err == nil {
		ok = true
	}

	if !ok {
		return nil, fmt.Errorf("no files found in /proc/device-tree")
	}

	return dt, nil
}
