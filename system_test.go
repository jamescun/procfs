// Copyright 2026 James Cunningham
// SPDX-License-Identifier: BSD-3-Clause
//
// Use of this source code is governed by a BSD-style license that can be found
// in the LICENSE file or at https://opensource.org/license/BSD-3-clause

package procfs

import (
	"os"
	"testing"
	"time"

	"github.com/google/go-cmp/cmp"
)

func TestCGroups(t *testing.T) {
	proc := From(os.DirFS("testdata/proc"))

	expected := Cgroups{
		{Name: "cpuset", NumCgroups: 1759, Enabled: true},
		{Name: "cpu", NumCgroups: 1759, Enabled: true},
		{Name: "cpuacct", NumCgroups: 1759, Enabled: true},
		{Name: "blkio", NumCgroups: 1759, Enabled: true},
		{Name: "memory", NumCgroups: 1759, Enabled: true},
		{Name: "devices", NumCgroups: 1759, Enabled: true},
		{Name: "freezer", NumCgroups: 1759, Enabled: true},
		{Name: "net_cls", NumCgroups: 1759, Enabled: true},
		{Name: "perf_event", NumCgroups: 1759, Enabled: true},
		{Name: "net_prio", NumCgroups: 1759, Enabled: true},
		{Name: "hugetlb", NumCgroups: 1759, Enabled: true},
		{Name: "pids", NumCgroups: 1759, Enabled: true},
		{Name: "rdma", NumCgroups: 1759, Enabled: true},
		{Name: "misc", NumCgroups: 1759, Enabled: true},
		{Name: "dmem", NumCgroups: 1759, Enabled: true},
	}

	target, err := GetCgroups(proc)
	if err != nil {
		t.Fatal("unexpected error:", err)
	}

	if !cmp.Equal(expected, target) {
		t.Error(cmp.Diff(expected, target))
	}
}

func TestCPUStats(t *testing.T) {
	proc := From(os.DirFS("testdata/proc"))

	expected := &CPUStats{
		Total: CPUStat{
			User:    7305537,
			Nice:    121612,
			System:  8966416,
			Idle:    19247232615,
			IOWait:  50564,
			IRQ:     2277989,
			SoftIRQ: 1522575,
		},
		CPU: []*CPUStat{
			{
				User:    637717,
				Nice:    9545,
				System:  712619,
				Idle:    1603991825,
				IOWait:  2979,
				IRQ:     171486,
				SoftIRQ: 150508,
			},
			{
				CPU:     1,
				User:    722067,
				Nice:    10165,
				System:  989545,
				Idle:    1603492829,
				IOWait:  5334,
				IRQ:     212102,
				SoftIRQ: 135687,
			},
			{
				CPU:     2,
				User:    683438,
				Nice:    10755,
				System:  903539,
				Idle:    1603621445,
				IOWait:  4119,
				IRQ:     217079,
				SoftIRQ: 142470,
			},
			{
				CPU:     3,
				User:    650600,
				Nice:    11209,
				System:  938359,
				Idle:    1603638622,
				IOWait:  4875,
				IRQ:     203242,
				SoftIRQ: 127266,
			},
			{
				CPU:     4,
				User:    723856,
				Nice:    13536,
				System:  982340,
				Idle:    1603426754,
				IOWait:  5909,
				IRQ:     228740,
				SoftIRQ: 177710,
			},
			{
				CPU:     5,
				User:    644415,
				Nice:    9319,
				System:  1022556,
				Idle:    1603505719,
				IOWait:  4909,
				IRQ:     212329,
				SoftIRQ: 136258,
			},
			{
				CPU:     6,
				User:    689617,
				Nice:    9515,
				System:  910987,
				Idle:    1603666256,
				IOWait:  4263,
				IRQ:     205287,
				SoftIRQ: 124174,
			},
			{
				CPU:     7,
				User:    593386,
				Nice:    10222,
				System:  647838,
				Idle:    1604151382,
				IOWait:  4265,
				IRQ:     161602,
				SoftIRQ: 88420,
			},
			{
				CPU:     8,
				User:    590274,
				Nice:    10283,
				System:  545032,
				Idle:    1604238210,
				IOWait:  3809,
				IRQ:     174132,
				SoftIRQ: 111149,
			},
			{
				CPU:     9,
				User:    579424,
				Nice:    9826,
				System:  529822,
				Idle:    1604285692,
				IOWait:  3716,
				IRQ:     168862,
				SoftIRQ: 101781,
			},
			{
				CPU:     10,
				User:    266490,
				Nice:    8765,
				System:  284367,
				Idle:    1604919890,
				IOWait:  2823,
				IRQ:     126592,
				SoftIRQ: 67011,
			},
			{
				CPU:     11,
				User:    524247,
				Nice:    8468,
				System:  499406,
				Idle:    1604293987,
				IOWait:  3557,
				IRQ:     196528,
				SoftIRQ: 160134,
			},
		},
		Ctxt:         10939017265,
		Btime:        1771085005,
		Processes:    5773423,
		ProcsRunning: 1,
	}

	target, err := GetCPUStats(proc)
	if err != nil {
		t.Fatal("unexpected error:", err)
	}

	if !cmp.Equal(expected, target) {
		t.Error(cmp.Diff(expected, target))
	}
}

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

func TestUptime(t *testing.T) {
	proc := From(os.DirFS("testdata/proc"))

	expected := &Uptime{
		Up:   16072645 * time.Second,
		Idle: 192637915 * time.Second,
	}

	target, err := GetUptime(proc)
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
