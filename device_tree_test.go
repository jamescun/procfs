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

func TestDeviceTree(t *testing.T) {
	proc := From(os.DirFS("testdata/proc"))

	expected := &DeviceTree{
		Model:      "Raspberry Pi 4 Model B Rev 1.1",
		Serial:     "1234567890123456",
		Compatible: []string{"raspberrypi", "4-model-bbrcm", "bcm2711"},
		BootArgs: []string{
			"coherent_pool=1M", "8250.nr_uarts=1", "snd_bcm2835.enable_headphones=0",
			"cgroup_disable=memory", "numa_policy=interleave",
			"nvme.max_host_mem_size_mb=0", "bcm2708_fb.fbwidth=0", "bcm2708_fb.fbheight=0",
			"bcm2708_fb.fbswap=1", "smsc95xx.macaddr=12:34:56:78:90:12",
			"vc_mem.mem_base=0x3eb00000", "vc_mem.mem_size=0x3ff00000",
			"console=ttyS0,115200", "console=tty1",
			"root=PARTUUID=3939c54e-e3fe-45cd-9c75-cf4ba65823d2", "rootfstype=ext4",
			"rootwait",
		},
	}

	target, err := GetDeviceTree(proc)
	if err != nil {
		t.Fatal("unexpected error:", err)
	}

	if !cmp.Equal(expected, target) {
		t.Error(cmp.Diff(expected, target))
	}
}
