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

func TestLinkStats(t *testing.T) {
	proc := From(os.DirFS("testdata/proc"))

	expected := LinkStats{
		{
			Name: "lo",
			Rx:   RxStats{Bytes: 14299578, Packets: 41068},
			Tx:   TxStats{Bytes: 14299578, Packets: 41068},
		},
		{
			Name: "eno1",
			Rx:   RxStats{Bytes: 43967995126, Packets: 62841560, Multicast: 21},
			Tx:   TxStats{Bytes: 14031633927, Packets: 64291149, Dropped: 6},
		},
		{
			Name: "wg0",
			Rx:   RxStats{Bytes: 15057317128, Packets: 15027831, Errors: 430, Frame: 430},
			Tx:   TxStats{Bytes: 2495423504, Packets: 9526411, Errors: 790, Dropped: 1981},
		},
		{
			Name: "podman0",
			Rx:   RxStats{Bytes: 7170041319, Packets: 27941108, Multicast: 14210},
			Tx:   TxStats{Bytes: 6204523344, Packets: 32919695},
		},
		{
			Name: "veth1",
			Rx:   RxStats{Bytes: 4214485363, Packets: 23032706},
			Tx:   TxStats{Bytes: 4418363939, Packets: 32296239},
		},
		{
			Name: "compute0",
			Rx:   RxStats{Bytes: 8348370420, Packets: 46930077, Multicast: 11284},
			Tx:   TxStats{Bytes: 23647084827, Packets: 42840140},
		},
		{
			Name: "veth0",
			Rx:   RxStats{Bytes: 370778814, Packets: 3042671},
			Tx:   TxStats{Bytes: 240879351, Packets: 3560573},
		},
		{
			Name: "veth2",
			Rx:   RxStats{Bytes: 1502740130, Packets: 5161595},
			Tx:   TxStats{Bytes: 645118244, Packets: 5540691},
		},
		{
			Name: "veth5",
			Rx:   RxStats{Bytes: 12817157, Packets: 139076},
			Tx:   TxStats{Bytes: 1013610102, Packets: 180345},
		},
		{
			Name: "veth7",
			Rx:   RxStats{Bytes: 547558117, Packets: 5134986},
			Tx:   TxStats{Bytes: 22916539157, Packets: 3332835},
		},
		{
			Name: "veth3",
			Rx:   RxStats{Bytes: 269551733, Packets: 1879005},
			Tx:   TxStats{Bytes: 313507546, Packets: 2872080},
		},
		{
			Name: "veth4",
			Rx:   RxStats{Bytes: 409222982, Packets: 1504284},
			Tx:   TxStats{Bytes: 273936754, Packets: 2157367},
		},
	}

	target, err := GetLinkStats(proc)
	if err != nil {
		t.Fatal("unexpected error:", err)
	}

	if !cmp.Equal(expected, target) {
		t.Error(cmp.Diff(expected, target))
	}
}
