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

func TestDiskStats(t *testing.T) {
	proc := From(os.DirFS("testdata/proc"))

	expected := DiskStats{
		{
			Major:          259,
			Name:           "nvme1n1",
			Reads:          389347453,
			ReadMerged:     20821683,
			ReadSectors:    52501372594,
			ReadTicks:      436688967,
			Writes:         75798107,
			WriteMerged:    14894112,
			WriteSectors:   748554584,
			WriteTicks:     7724419,
			IOTicks:        9630357,
			TimeInQueue:    444432805,
			Discards:       91084,
			DiscardMerged:  20,
			DiscardSectors: 2991609264,
			DiscardTicks:   19417,
		},
		{
			Major:        259,
			Minor:        5,
			Name:         "nvme1n1p1",
			Reads:        13907394,
			ReadMerged:   758016,
			ReadSectors:  1877158460,
			ReadTicks:    1146875,
			Writes:       476,
			WriteSectors: 476,
			WriteTicks:   40,
			IOTicks:      539648,
			TimeInQueue:  1146915,
		},
		{
			Major:          259,
			Minor:          6,
			Name:           "nvme1n1p2",
			Reads:          424780,
			ReadMerged:     33247,
			ReadSectors:    58609204,
			ReadTicks:      375610,
			Writes:         3310,
			WriteMerged:    733,
			WriteSectors:   2539800,
			WriteTicks:     67968,
			IOTicks:        15114,
			TimeInQueue:    443668,
			Discards:       772,
			DiscardSectors: 3046360,
			DiscardTicks:   89,
		},
		{
			Major:          259,
			Minor:          7,
			Name:           "nvme1n1p3",
			Reads:          375015013,
			ReadMerged:     20030420,
			ReadSectors:    50565592738,
			ReadTicks:      435166430,
			Writes:         75794321,
			WriteMerged:    14893379,
			WriteSectors:   746014308,
			WriteTicks:     7656410,
			IOTicks:        16393743,
			TimeInQueue:    442842169,
			Discards:       90312,
			DiscardMerged:  20,
			DiscardSectors: 2988562904,
			DiscardTicks:   19328,
		},
		{
			Major:          259,
			Minor:          1,
			Name:           "nvme0n1",
			Reads:          389371423,
			ReadMerged:     20822077,
			ReadSectors:    52502294465,
			ReadTicks:      433155021,
			Writes:         75788673,
			WriteMerged:    14903546,
			WriteSectors:   748554584,
			WriteTicks:     6943728,
			IOTicks:        10671484,
			TimeInQueue:    440114407,
			Discards:       91084,
			DiscardMerged:  20,
			DiscardSectors: 2991609264,
			DiscardTicks:   15657,
		},
		{
			Major:        259,
			Minor:        2,
			Name:         "nvme0n1p1",
			Reads:        13907480,
			ReadMerged:   758012,
			ReadSectors:  1877161596,
			ReadTicks:    16667794,
			Writes:       476,
			WriteSectors: 476,
			WriteTicks:   115,
			IOTicks:      583924,
			TimeInQueue:  16667909,
		},
		{
			Major:          259,
			Minor:          3,
			Name:           "nvme0n1p2",
			Reads:          424986,
			ReadMerged:     33338,
			ReadSectors:    58613862,
			ReadTicks:      261619,
			Writes:         3307,
			WriteMerged:    736,
			WriteSectors:   2539800,
			WriteTicks:     58764,
			IOTicks:        14891,
			TimeInQueue:    320464,
			Discards:       772,
			DiscardSectors: 3046360,
			DiscardTicks:   80,
		},
		{
			Major:          259,
			Minor:          4,
			Name:           "nvme0n1p3",
			Reads:          375038600,
			ReadMerged:     20030727,
			ReadSectors:    50566495783,
			ReadTicks:      416225553,
			Writes:         75784890,
			WriteMerged:    14902810,
			WriteSectors:   746014308,
			WriteTicks:     6884849,
			IOTicks:        16235793,
			TimeInQueue:    423125980,
			Discards:       90312,
			DiscardMerged:  20,
			DiscardSectors: 2988562904,
			DiscardTicks:   15577,
		},
		{
			Major:       9,
			Name:        "md0",
			Reads:       100,
			ReadSectors: 4704,
			ReadTicks:   6,
			IOTicks:     5,
			TimeInQueue: 6,
		},
		{
			Major:          9,
			Minor:          2,
			Name:           "md2",
			Reads:          28466,
			ReadSectors:    1087930,
			ReadTicks:      4407,
			Writes:         54031602,
			WriteSectors:   644865632,
			WriteTicks:     12300570,
			IOTicks:        307002453,
			TimeInQueue:    12325056,
			Discards:       90332,
			DiscardSectors: 2988562904,
			DiscardTicks:   20079,
		},
		{
			Major:          9,
			Minor:          1,
			Name:           "md1",
			Reads:          397,
			ReadSectors:    8258,
			ReadTicks:      85,
			Writes:         3291,
			WriteSectors:   2539048,
			WriteTicks:     60684,
			IOTicks:        996,
			TimeInQueue:    60859,
			Discards:       772,
			DiscardSectors: 3046360,
			DiscardTicks:   90,
		},
	}

	target, err := GetDiskStats(proc)
	if err != nil {
		t.Fatal("unexpected error:", err)
	}

	if !cmp.Equal(expected, target) {
		t.Error(cmp.Diff(expected, target))
	}
}

func TestFilesystems(t *testing.T) {
	proc := From(os.DirFS("testdata/proc"))

	expected := Filesystems{
		{NoDev: true, Name: "sysfs"},
		{NoDev: true, Name: "tmpfs"},
		{NoDev: true, Name: "bdev"},
		{NoDev: true, Name: "proc"},
		{NoDev: true, Name: "cgroup"},
		{NoDev: true, Name: "cgroup2"},
		{NoDev: true, Name: "cpuset"},
		{NoDev: true, Name: "devtmpfs"},
		{NoDev: true, Name: "configfs"},
		{NoDev: true, Name: "debugfs"},
		{NoDev: true, Name: "tracefs"},
		{NoDev: true, Name: "securityfs"},
		{NoDev: true, Name: "sockfs"},
		{NoDev: true, Name: "bpf"},
		{NoDev: true, Name: "pipefs"},
		{NoDev: true, Name: "ramfs"},
		{NoDev: true, Name: "hugetlbfs"},
		{NoDev: true, Name: "devpts"},
		{NoDev: true, Name: "autofs"},
		{NoDev: true, Name: "efivarfs"},
		{NoDev: true, Name: "mqueue"},
		{NoDev: true, Name: "selinuxfs"},
		{NoDev: true, Name: "pstore"},
		{NoDev: false, Name: "btrfs"},
		{NoDev: false, Name: "ext3"},
		{NoDev: false, Name: "ext2"},
		{NoDev: false, Name: "ext4"},
		{NoDev: false, Name: "fuseblk"},
		{NoDev: true, Name: "fuse"},
		{NoDev: true, Name: "fusectl"},
		{NoDev: true, Name: "overlay"},
		{NoDev: true, Name: "binfmt_misc"},
	}

	target, err := GetFilesystems(proc)
	if err != nil {
		t.Fatal("unexpected error:", err)
	}

	if !cmp.Equal(expected, target) {
		t.Error(cmp.Diff(expected, target))
	}
}

func TestMounts(t *testing.T) {
	proc := From(os.DirFS("testdata/proc"))

	expected := Mounts{
		{
			ID:           67,
			ParentID:     1,
			Device:       "9:2",
			Root:         "/",
			Path:         "/",
			Options:      []string{"rw", "relatime"},
			Tags:         []string{"shared:1"},
			Type:         "ext4",
			Source:       "/dev/md2",
			SuperOptions: []string{"rw", "seclabel"},
		},
		{
			ID:       34,
			ParentID: 67,
			Device:   "0:6",
			Root:     "/",
			Path:     "/dev",
			Options:  []string{"rw", "nosuid"},
			Tags:     []string{"shared:2"},
			Type:     "devtmpfs",
			Source:   "devtmpfs",
			SuperOptions: []string{
				"rw", "seclabel", "size=32677776k", "nr_inodes=8169444",
				"mode=755", "inode64",
			},
		},
		{
			ID:           36,
			ParentID:     34,
			Device:       "0:24",
			Root:         "/",
			Path:         "/dev/shm",
			Options:      []string{"rw", "nosuid", "nodev"},
			Tags:         []string{"shared:3"},
			Type:         "tmpfs",
			Source:       "tmpfs",
			SuperOptions: []string{"rw", "seclabel", "inode64"},
		},
		{
			ID:       37,
			ParentID: 34,
			Device:   "0:25",
			Root:     "/",
			Path:     "/dev/pts",
			Options:  []string{"rw", "nosuid", "noexec", "relatime"},
			Tags:     []string{"shared:4"},
			Type:     "devpts",
			Source:   "devpts",
			SuperOptions: []string{
				"rw", "seclabel", "gid=5", "mode=620", "ptmxmode=000",
			},
		},
		{
			ID:       38,
			ParentID: 67,
			Device:   "0:23",
			Root:     "/",
			Path:     "/sys",
			Options: []string{
				"rw", "nosuid", "nodev", "noexec", "relatime",
			},
			Tags:         []string{"shared:5"},
			Type:         "sysfs",
			Source:       "sysfs",
			SuperOptions: []string{"rw", "seclabel"},
		},
		{
			ID:       39,
			ParentID: 38,
			Device:   "0:7",
			Root:     "/",
			Path:     "/sys/kernel/security",
			Options: []string{
				"rw", "nosuid", "nodev", "noexec", "relatime",
			},
			Tags:         []string{"shared:6"},
			Type:         "securityfs",
			Source:       "securityfs",
			SuperOptions: []string{"rw"},
		},
		{
			ID:       40,
			ParentID: 38,
			Device:   "0:27",
			Root:     "/",
			Path:     "/sys/fs/cgroup",
			Options: []string{
				"rw", "nosuid", "nodev", "noexec", "relatime",
			},
			Tags:   []string{"shared:7"},
			Type:   "cgroup2",
			Source: "cgroup2",
			SuperOptions: []string{
				"rw", "seclabel", "nsdelegate", "memory_recursiveprot",
			},
		},
		{
			ID:       41,
			ParentID: 38,
			Device:   "0:28",
			Root:     "/",
			Path:     "/sys/fs/pstore",
			Options: []string{
				"rw", "nosuid", "nodev", "noexec", "relatime",
			},
			Tags:         []string{"shared:8"},
			Type:         "pstore",
			Source:       "pstore",
			SuperOptions: []string{"rw", "seclabel"},
		},
		{
			ID:       42,
			ParentID: 38,
			Device:   "0:29",
			Root:     "/",
			Path:     "/sys/fs/bpf",
			Options: []string{
				"rw", "nosuid", "nodev", "noexec", "relatime",
			},
			Tags:         []string{"shared:9"},
			Type:         "bpf",
			Source:       "bpf",
			SuperOptions: []string{"rw", "mode=700"},
		},
		{
			ID:       43,
			ParentID: 38,
			Device:   "0:30",
			Root:     "/",
			Path:     "/sys/kernel/config",
			Options: []string{
				"rw", "nosuid", "nodev", "noexec", "relatime",
			},
			Tags:         []string{"shared:10"},
			Type:         "configfs",
			Source:       "configfs",
			SuperOptions: []string{"rw"},
		},
		{
			ID:           44,
			ParentID:     67,
			Device:       "0:22",
			Root:         "/",
			Path:         "/proc",
			Options:      []string{"rw", "relatime"},
			Tags:         []string{"shared:12"},
			Type:         "proc",
			Source:       "proc",
			SuperOptions: []string{"rw"},
		},
		{
			ID:       45,
			ParentID: 67,
			Device:   "0:26",
			Root:     "/",
			Path:     "/run",
			Options:  []string{"rw", "nosuid", "nodev"},
			Tags:     []string{"shared:13"},
			Type:     "tmpfs",
			Source:   "tmpfs",
			SuperOptions: []string{
				"rw", "seclabel", "size=13085996k", "nr_inodes=819200",
				"mode=755", "inode64",
			},
		},
		{
			ID:           24,
			ParentID:     38,
			Device:       "0:21",
			Root:         "/",
			Path:         "/sys/fs/selinux",
			Options:      []string{"rw", "nosuid", "noexec", "relatime"},
			Tags:         []string{"shared:11"},
			Type:         "selinuxfs",
			Source:       "selinuxfs",
			SuperOptions: []string{"rw"},
		},
		{
			ID:       23,
			ParentID: 44,
			Device:   "0:31",
			Root:     "/",
			Path:     "/proc/sys/fs/binfmt_misc",
			Options:  []string{"rw", "relatime"},
			Tags:     []string{"shared:14"},
			Type:     "autofs",
			Source:   "systemd-1",
			SuperOptions: []string{
				"rw", "fd=35", "pgrp=1", "timeout=0", "minproto=5",
				"maxproto=5", "direct", "pipe_ino=305",
			},
		},
		{
			ID:       25,
			ParentID: 38,
			Device:   "0:13",
			Root:     "/",
			Path:     "/sys/kernel/tracing",
			Options: []string{
				"rw", "nosuid", "nodev", "noexec", "relatime",
			},
			Tags:         []string{"shared:15"},
			Type:         "tracefs",
			Source:       "tracefs",
			SuperOptions: []string{"rw", "seclabel"},
		},
		{
			ID:       26,
			ParentID: 34,
			Device:   "0:20",
			Root:     "/",
			Path:     "/dev/mqueue",
			Options: []string{
				"rw", "nosuid", "nodev", "noexec", "relatime",
			},
			Tags:         []string{"shared:16"},
			Type:         "mqueue",
			Source:       "mqueue",
			SuperOptions: []string{"rw", "seclabel"},
		},
		{
			ID:           27,
			ParentID:     34,
			Device:       "0:32",
			Root:         "/",
			Path:         "/dev/hugepages",
			Options:      []string{"rw", "nosuid", "nodev", "relatime"},
			Tags:         []string{"shared:17"},
			Type:         "hugetlbfs",
			Source:       "hugetlbfs",
			SuperOptions: []string{"rw", "seclabel", "pagesize=2M"},
		},
		{
			ID:       28,
			ParentID: 38,
			Device:   "0:8",
			Root:     "/",
			Path:     "/sys/kernel/debug",
			Options: []string{
				"rw", "nosuid", "nodev", "noexec", "relatime",
			},
			Tags:         []string{"shared:18"},
			Type:         "debugfs",
			Source:       "debugfs",
			SuperOptions: []string{"rw", "seclabel"},
		},
		{
			ID:       33,
			ParentID: 38,
			Device:   "0:37",
			Root:     "/",
			Path:     "/sys/fs/fuse/connections",
			Options: []string{
				"rw", "nosuid", "nodev", "noexec", "relatime",
			},
			Tags:         []string{"shared:67"},
			Type:         "fusectl",
			Source:       "fusectl",
			SuperOptions: []string{"rw"},
		},
		{
			ID:           47,
			ParentID:     67,
			Device:       "9:1",
			Root:         "/",
			Path:         "/boot",
			Options:      []string{"rw", "relatime"},
			Tags:         []string{"shared:70"},
			Type:         "ext3",
			Source:       "/dev/md1",
			SuperOptions: []string{"rw", "seclabel"},
		},
		{
			ID:       342,
			ParentID: 45,
			Device:   "0:45",
			Root:     "/",
			Path:     "/run/credentials/getty@tty1.service",
			Options: []string{
				"ro", "nosuid", "nodev", "noexec", "relatime", "nosymfollow",
			},
			Tags:   []string{"shared:247"},
			Type:   "tmpfs",
			Source: "tmpfs",
			SuperOptions: []string{
				"rw", "seclabel", "size=1024k", "nr_inodes=1024", "mode=700",
				"inode64", "noswap",
			},
		},
		{
			ID:           174,
			ParentID:     67,
			Device:       "9:2",
			Root:         "/var/lib/containers/storage/overlay",
			Path:         "/var/lib/containers/storage/overlay",
			Options:      []string{"rw", "relatime"},
			Type:         "ext4",
			Source:       "/dev/md2",
			SuperOptions: []string{"rw", "seclabel"},
		},
		{
			ID:       278,
			ParentID: 23,
			Device:   "0:76",
			Root:     "/",
			Path:     "/proc/sys/fs/binfmt_misc",
			Options: []string{
				"rw", "nosuid", "nodev", "noexec", "relatime",
			},
			Tags:         []string{"shared:279"},
			Type:         "binfmt_misc",
			Source:       "binfmt_misc",
			SuperOptions: []string{"rw"},
		},
		{
			ID:       495,
			ParentID: 45,
			Device:   "0:26",
			Root:     "/netns",
			Path:     "/run/netns",
			Options:  []string{"rw", "nosuid", "nodev"},
			Tags:     []string{"shared:13"},
			Type:     "tmpfs",
			Source:   "tmpfs",
			SuperOptions: []string{
				"rw", "seclabel", "size=13085996k", "nr_inodes=819200",
				"mode=755", "inode64",
			},
		},
	}

	target, err := GetMounts(proc)
	if err != nil {
		t.Fatal("unexpected error:", err)
	}

	if !cmp.Equal(expected, target) {
		t.Error(cmp.Diff(expected, target))
	}
}

func TestPartitions(t *testing.T) {
	proc := From(os.DirFS("testdata/proc"))

	expected := Partitions{
		{Major: 259, Blocks: 937692504, Name: "nvme1n1"},
		{Major: 259, Minor: 5, Blocks: 33554432, Name: "nvme1n1p1"},
		{Major: 259, Minor: 6, Blocks: 1048576, Name: "nvme1n1p2"},
		{Major: 259, Minor: 7, Blocks: 903087448, Name: "nvme1n1p3"},
		{Major: 259, Minor: 1, Blocks: 937692504, Name: "nvme0n1"},
		{Major: 259, Minor: 2, Blocks: 33554432, Name: "nvme0n1p1"},
		{Major: 259, Minor: 3, Blocks: 1048576, Name: "nvme0n1p2"},
		{Major: 259, Minor: 4, Blocks: 903087448, Name: "nvme0n1p3"},
		{Major: 9, Blocks: 33520640, Name: "md0"},
		{Major: 9, Minor: 2, Blocks: 902955328, Name: "md2"},
		{Major: 9, Minor: 1, Blocks: 1046528, Name: "md1"},
	}

	target, err := GetPartitions(proc)
	if err != nil {
		t.Fatal("unexpected error:", err)
	}

	if !cmp.Equal(expected, target) {
		t.Error(cmp.Diff(expected, target))
	}
}
