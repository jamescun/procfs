// Copyright 2026 James Cunningham
// SPDX-License-Identifier: BSD-3-Clause
//
// Use of this source code is governed by a BSD-style license that can be found
// in the LICENSE file or at https://opensource.org/license/BSD-3-clause

package procfs

import (
	"strings"

	"go.jamescun.com/procfs/internal/utils"
)

// DiskStat contains the I/O statistics for a single block storage device, read
// from /proc/diskstats.
//
// References:
//   - linux/Documentation/admin-guide/iostats.rst
type DiskStat struct {
	Major          int
	Minor          int
	Name           string
	Reads          uint64
	ReadMerged     uint64
	ReadSectors    uint64
	ReadTicks      uint64
	Writes         uint64
	WriteMerged    uint64
	WriteSectors   uint64
	WriteTicks     uint64
	InFlight       uint64
	IOTicks        uint64
	TimeInQueue    uint64
	Discards       uint64
	DiscardMerged  uint64
	DiscardSectors uint64
	DiscardTicks   uint64
	Flushes        uint64
	FlushTicks     uint64
}

// UnmarshalText unmarshals a single line from /proc/diskstats.
func (d *DiskStat) UnmarshalText(b []byte) error {
	for i, field := range utils.Fields(b) {
		switch i {
		case 0:
			d.Major, _ = utils.Int[int](field)
		case 1:
			d.Minor, _ = utils.Int[int](field)
		case 2:
			d.Name = string(field)

		case 3:
			d.Reads, _ = utils.Uint[uint64](field)
		case 4:
			d.ReadMerged, _ = utils.Uint[uint64](field)
		case 5:
			d.ReadSectors, _ = utils.Uint[uint64](field)
		case 6:
			d.ReadTicks, _ = utils.Uint[uint64](field)

		case 7:
			d.Writes, _ = utils.Uint[uint64](field)
		case 8:
			d.WriteMerged, _ = utils.Uint[uint64](field)
		case 9:
			d.WriteSectors, _ = utils.Uint[uint64](field)
		case 10:
			d.WriteTicks, _ = utils.Uint[uint64](field)

		case 11:
			d.InFlight, _ = utils.Uint[uint64](field)
		case 12:
			d.IOTicks, _ = utils.Uint[uint64](field)
		case 13:
			d.TimeInQueue, _ = utils.Uint[uint64](field)

		case 14:
			d.Discards, _ = utils.Uint[uint64](field)
		case 15:
			d.DiscardMerged, _ = utils.Uint[uint64](field)
		case 16:
			d.DiscardSectors, _ = utils.Uint[uint64](field)
		case 17:
			d.DiscardTicks, _ = utils.Uint[uint64](field)

		case 18:
			d.Flushes, _ = utils.Uint[uint64](field)
		case 19:
			d.FlushTicks, _ = utils.Uint[uint64](field)
		}
	}

	return nil
}

// DiskStats contains I/O statistics for the systems block storage devices,
// read from /proc/diskstats.
//
// References:
//   - linux/Documentation/admin-guide/iostats.rst
type DiskStats []*DiskStat

// GetDiskStats reads I/O statistics for the systems block storage devices,
// read from /proc/diskstats in the given [Procfs].
func GetDiskStats(proc Procfs) (DiskStats, error) {
	ds := DiskStats{}

	err := proc.Read("diskstats", &ds)
	if err != nil {
		return nil, err
	}

	return ds, nil
}

// UnmarshalText unmarshals the lines from /proc/diskstats.
func (ds *DiskStats) UnmarshalText(b []byte) error {
	for _, line := range utils.Lines(b) {
		d := new(DiskStat)

		err := d.UnmarshalText(line)
		if err != nil {
			return err
		}

		*ds = append(*ds, d)
	}

	return nil
}

// Filesystem is one of the filesystems supported by the kernel, read from
// /proc/filesystems.
//
// References:
//   - proc_filesystems(5)
type Filesystem struct {
	NoDev bool
	Name  string
}

// UnmarshalText unmarshals a single line from /proc/filesystems.
func (f *Filesystem) UnmarshalText(b []byte) error {
	const nodev = "nodev"

	for _, field := range utils.Fields(b) {
		if utils.Equal(field, nodev) {
			f.NoDev = true
			continue
		}

		f.Name = string(field)
	}

	return nil
}

// Filesystems are the filesystems supported by the kernel, read from
// /proc/filesystems.
//
// References:
//   - proc_filesystems(5)
type Filesystems []*Filesystem

// GetFilesystems reads the filesystems supported by the kernel, read from
// /proc/filesystems in the given [Procfs].
func GetFilesystems(proc Procfs) (Filesystems, error) {
	fs := Filesystems{}

	err := proc.Read("filesystems", &fs)
	if err != nil {
		return nil, err
	}

	return fs, nil
}

// UnmarshalText unmarshals the lines from /proc/filesystems.
func (fs *Filesystems) UnmarshalText(b []byte) error {
	for _, line := range utils.Lines(b) {
		f := new(Filesystem)

		err := f.UnmarshalText(line)
		if err != nil {
			return err
		}

		*fs = append(*fs, f)
	}

	return nil
}

// Mount is a single mounted filesystem visible in a process' namespace, read
// from /proc/<pid>/mountinfo.
//
// References:
//   - proc_pid_mountinfo(5)
type Mount struct {
	ID           int
	ParentID     int
	Device       string
	Root         string
	Path         string
	Options      []string
	Tags         []string
	Type         string
	Source       string
	SuperOptions []string
}

// UnmarshalText unmarshals a single line from /proc/<pid>/mountinfo.
func (m *Mount) UnmarshalText(b []byte) error {
	// the contents are mostly string based, cast to a string and cut that.
	text := string(b)

	// cut before and after extended options.
	before, after, _ := strings.Cut(text, " - ")

	for i, field := range utils.Fields(before) {
		switch i {
		case 0:
			m.ID, _ = utils.Int[int](field)
		case 1:
			m.ParentID, _ = utils.Int[int](field)
		case 2:
			m.Device = field
		case 3:
			m.Root = field
		case 4:
			m.Path = field
		case 5:
			m.Options = strings.Split(field, ",")
		default:
			m.Tags = append(m.Tags, field)
		}
	}

	for i, field := range utils.Fields(after) {
		switch i {
		case 0:
			m.Type = field
		case 1:
			m.Source = field
		case 2:
			m.SuperOptions = strings.Split(field, ",")
		}
	}

	return nil
}

// Mounts are the mounted filesystems visible in a process' namespace, read
// from /proc/<pid>/mountinfo.
//
// References:
//   - proc_pid_mountinfo(5)
type Mounts []*Mount

// GetMounts reads the mounted filesystems visible in the current process'
// namespace, read from /proc/self/mountinfo in the given [Procfs].
func GetMounts(proc Procfs) (Mounts, error) {
	ms := Mounts{}

	err := proc.Read("self/mountinfo", &ms)
	if err != nil {
		return nil, err
	}

	return ms, nil
}

// UnmarshalText unmarshals the lines from /proc/<pid>/mountinfo.
func (ms *Mounts) UnmarshalText(b []byte) error {
	for _, line := range utils.Lines(b) {
		m := &Mount{}
		err := m.UnmarshalText(line)
		if err != nil {
			return err
		}

		*ms = append(*ms, m)
	}

	return nil
}

// Partition is one of the partitions read from /proc/partitions.
//
// References:
//   - proc_partitions(5)
type Partition struct {
	Major  int
	Minor  int
	Blocks uint64
	Name   string
}

// UnmarshalText unmarshals a single line from /proc/partitions.
func (p *Partition) UnmarshalText(b []byte) error {
	for i, field := range utils.Fields(b) {
		switch i {
		case 0:
			p.Major, _ = utils.Int[int](field)
		case 1:
			p.Minor, _ = utils.Int[int](field)
		case 2:
			p.Blocks, _ = utils.Uint[uint64](field)
		case 3:
			p.Name = string(field)
		}
	}

	return nil
}

// Partitions are the partitions read from /proc/partitions.
//
// References:
//   - proc_partitions(5)
type Partitions []*Partition

// GetPartitions reads the systems partitions from /proc/partitions in the
// given [Procfs].
func GetPartitions(proc Procfs) (Partitions, error) {
	ps := Partitions{}

	err := proc.Read("partitions", &ps)
	if err != nil {
		return nil, err
	}

	return ps, nil
}

// UnmarshalText unmarshals the lines from /proc/partitions.
func (ps *Partitions) UnmarshalText(b []byte) error {
	for _, line := range utils.Lines(b) {
		if line[0] == 'm' {
			// skip header.
			continue
		}

		p := new(Partition)
		err := p.UnmarshalText(line)
		if err != nil {
			return err
		}

		*ps = append(*ps, p)
	}

	return nil
}
