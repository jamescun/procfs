// Copyright 2026 James Cunningham
// SPDX-License-Identifier: BSD-3-Clause
//
// Use of this source code is governed by a BSD-style license that can be found
// in the LICENSE file or at https://opensource.org/license/BSD-3-clause

package procfs

import (
	"os"
	"testing"

	"github.com/google/go-cmp/cmp"
)

func TestInputDevices(t *testing.T) {
	proc := From(os.DirFS("testdata/proc"))

	expected := InputDevices{
		{
			Bus:      "0019",
			Vendor:   "0000",
			Product:  "0001",
			Version:  "0000",
			Name:     "Power Button",
			Phys:     "PNP0C0C/button/input0",
			Sysfs:    "/devices/platform/PNP0C0C:00/input/input0",
			Handlers: []string{"kbd", "event0"},
		},
		{
			Bus:      "0019",
			Vendor:   "0000",
			Product:  "0001",
			Version:  "0000",
			Name:     "Power Button",
			Phys:     "LNXPWRBN/button/input0",
			Sysfs:    "/devices/platform/LNXPWRBN:00/input/input1",
			Handlers: []string{"kbd", "event1"},
		},
		{
			Bus:      "0019",
			Vendor:   "0000",
			Product:  "0006",
			Version:  "0000",
			Name:     "Video Bus",
			Phys:     "LNXVIDEO/video/input0",
			Sysfs:    "/devices/pci0000:00/0000:00:08.1/acpi.video_bus.0/input/input5",
			Handlers: []string{"kbd", "event5"},
		},
		{
			Bus:      "0010",
			Vendor:   "001f",
			Product:  "0001",
			Version:  "0100",
			Name:     "PC Speaker",
			Phys:     "isa0061/input0",
			Sysfs:    "/devices/platform/pcspkr/input/input7",
			Handlers: []string{"kbd", "event7"},
		},
		{
			Bus:      "0000",
			Vendor:   "0000",
			Product:  "0000",
			Version:  "0000",
			Name:     "HD-Audio Generic HDMI/DP,pcm=3",
			Phys:     "ALSA",
			Sysfs:    "/devices/pci0000:00/0000:00:08.1/0000:11:00.1/sound/card1/input9",
			Handlers: []string{"event9"},
		},
		{
			Bus:      "0005",
			Vendor:   "046d",
			Product:  "b019",
			Version:  "0006",
			Name:     "Logitech Wireless Mouse MX Master 2S",
			Phys:     "00:11:22:33:44:55",
			Sysfs:    "/devices/virtual/misc/uhid/0005:046D:B019.0007/input/input24",
			Uniq:     "00:11:22:33:44:55",
			Handlers: []string{"sysrq", "kbd", "event21", "mouse1"},
		},
	}

	target, err := GetInputDevices(proc)
	if err != nil {
		t.Fatal("unexpected error:", err)
	}

	if !cmp.Equal(expected, target) {
		t.Error(cmp.Diff(expected, target))
	}
}
