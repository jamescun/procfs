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

func TestLoadAvg(t *testing.T) {
	proc := From(os.DirFS("testdata/proc"))

	expected := &LoadAvg{
		Load: [3]float64{
			0.10, 0.07, 0.01,
		},
		Runnable:  1,
		Scheduled: 458,
		LastPID:   1558584,
	}

	target, err := GetLoadAvg(proc)
	if err != nil {
		t.Fatal("unexpected error:", err)
	}

	if !cmp.Equal(expected, target) {
		t.Error(cmp.Diff(expected, target))
	}
}

func TestMeminfo(t *testing.T) {
	proc := From(os.DirFS("testdata/proc"))

	expected := &Meminfo{
		MemTotal:      65429972,
		MemFree:       31539768,
		MemAvailable:  63242600,
		Buffers:       1434928,
		Cached:        29712264,
		Active:        3412976,
		Inactive:      28580052,
		ActiveAnon:    923536,
		ActiveFile:    2489440,
		InactiveFile:  28580052,
		Unevictable:   124,
		SwapTotal:     33520636,
		SwapFree:      33520636,
		Dirty:         40,
		AnonPages:     840448,
		Mapped:        649544,
		Shmem:         77836,
		KReclaimable:  1361156,
		Slab:          1662684,
		SReclaimable:  1361156,
		SUnreclaim:    301528,
		KernelStack:   7248,
		PageTables:    12192,
		SecPageTables: 2056,
		CommitLimit:   66235620,
		CommittedAS:   2098140,
		VmallocTotal:  34359738367,
		VmallocUsed:   57684,
		Percpu:        14144,
		AnonHugePages: 684032,
		Hugepagesize:  2048,
		DirectMap4k:   359964,
		DirectMap2M:   11026432,
		DirectMap1G:   55574528,
	}

	target, err := GetMeminfo(proc)
	if err != nil {
		t.Fatal("unexpected error:", err)
	}

	if !cmp.Equal(expected, target) {
		t.Error(cmp.Diff(expected, target))
	}
}

func TestModules(t *testing.T) {
	proc := From(os.DirFS("testdata/proc"))

	expected := Modules{
		{
			Name:  "wireguard",
			Size:  69632,
			State: "Live",
		},
		{
			Name:      "ip6_udp_tunnel",
			Size:      12288,
			Instances: 1,
			Depends:   []string{"wireguard"},
			State:     "Live",
		},
		{
			Name:      "udp_tunnel",
			Size:      24576,
			Instances: 1,
			Depends:   []string{"wireguard"},
			State:     "Live",
		},
		{
			Name:      "libchacha20poly1305",
			Size:      12288,
			Instances: 1,
			Depends:   []string{"wireguard"},
			State:     "Live",
		},
		{
			Name:      "libcurve25519",
			Size:      32768,
			Instances: 1,
			Depends:   []string{"wireguard"},
			State:     "Live",
		},
		{
			Name:      "libpoly1305",
			Size:      16384,
			Instances: 1,
			Depends:   []string{"libchacha20poly1305"},
			State:     "Live",
		},
	}

	target, err := GetModules(proc)
	if err != nil {
		t.Fatal("unexpected error:", err)
	}

	if !cmp.Equal(expected, target) {
		t.Error(cmp.Diff(expected, target))
	}
}

func TestSwaps(t *testing.T) {
	proc := From(os.DirFS("testdata/proc"))

	expected := Swaps{
		{Name: "/dev/md0", Type: "partition", Size: 33520636, Priority: -2},
	}

	target, err := GetSwaps(proc)
	if err != nil {
		t.Fatal("unexpected error:", err)
	}

	if !cmp.Equal(expected, target) {
		t.Error(cmp.Diff(expected, target))
	}
}

func TestVMStat(t *testing.T) {
	proc := From(os.DirFS("testdata/proc"))

	expected := &VMStat{
		NrFreePages:                1413775,
		NrFreePagesBlocks:          414208,
		NrZoneInactiveAnon:         20804,
		NrZoneActiveAnon:           4249,
		NrZoneInactiveFile:         5485,
		NrZoneActiveFile:           40673,
		NrZoneWritePending:         32,
		NrZsPages:                  5,
		NrInactiveAnon:             20804,
		NrActiveAnon:               4249,
		NrInactiveFile:             5485,
		NrActiveFile:               40673,
		NrSlabReclaimable:          457732,
		NrSlabUnreclaimable:        75180,
		WorkingSetNodes:            92812,
		WorkingSetRefaultFile:      270758664,
		WorkingSetActivateFile:     78372347,
		WorkingSetRestoreFile:      8358030,
		NrAnonPages:                18419,
		NrMapped:                   24121,
		NrFilePages:                51884,
		NrDirty:                    32,
		NrShmem:                    5653,
		NrVmscanImmediateReclaim:   6148,
		NrDirtied:                  39696406,
		NrWritten:                  39151166,
		NrFollPinAcquired:          25048,
		NrFollPinReleased:          25048,
		NrKernelStack:              4016,
		NrPageTablePages:           660,
		NrKernelFilePages:          887,
		NrDirtyThreshold:           287340,
		NrDirtyBackgroundThreshold: 143494,
		NrMemmapBootPages:          32768,
		Pgpgin:                     896393037,
		Pgpgout:                    156742252,
		PgAllocDMA:                 178367,
		PgAllocNormal:              593108903,
		PgFree:                     594710212,
		PgActivate:                 10007239,
		PgLazyFree:                 40061,
		PgFault:                    247274650,
		PgMajFault:                 932938,
		PgLazyFreed:                12029,
		PgReuse:                    54652322,
		PgScanFile:                 306103876,
		PgStealFile:                285814730,
		PgRotated:                  9335272,
		PgMigrateSuccess:           498,
		CompactMigrateScanned:      2365,
		CompactFreeScanned:         771,
		CompactIsolated:            1266,
		CompactStall:               2,
		CompactSuccess:             2,
		UnevictablePgsCulled:       3648,
		UnevictablePgsRescued:      3648,
		UnevictablePgsMlocked:      3648,
		UnevictablePgsMunlocked:    3648,
		ThpCollapseAlloc:           7,
		ThpScanExceedNonePte:       10621,
		ThpZeroPageAlloc:           1,
	}

	target, err := GetVMStat(proc)
	if err != nil {
		t.Fatal("unexpected error:", err)
	}

	if !cmp.Equal(expected, target) {
		t.Error(cmp.Diff(expected, target))
	}
}
