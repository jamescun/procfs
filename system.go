// Copyright 2026 James Cunningham
// SPDX-License-Identifier: BSD-3-Clause
//
// Use of this source code is governed by a BSD-style license that can be found
// in the LICENSE file or at https://opensource.org/license/BSD-3-clause

package procfs

import (
	"strconv"

	"go.jamescun.com/procfs/internal/utils"
)

// LoadAvg contains information about the systems load and scheduling, read
// from /proc/loadavg.
//
// References:
//   - proc_loadavg(5)
type LoadAvg struct {
	// Load is the 1, 5 and 15 minute load average.
	Load [3]float64

	Runnable  int
	Scheduled int
	LastPID   int
}

// GetLoadAvg reads system load and scheduling information, from /proc/loadavg
// in the given [Procfs].
func GetLoadAvg(proc Procfs) (*LoadAvg, error) {
	l := new(LoadAvg)

	err := proc.Read("loadavg", l)
	if err != nil {
		return nil, err
	}

	return l, nil
}

// UnmarshalText unmarshals the contents of /proc/loadavg.
func (l *LoadAvg) UnmarshalText(b []byte) error {
	for i, field := range utils.Fields(b) {
		if i < 3 {
			value, err := strconv.ParseFloat(string(field), 64)
			if err != nil {
				return err
			}

			l.Load[i] = value
		}

		if i == 3 {
			r, c, split := utils.Split(field, func(b byte) bool {
				return b == '/'
			})
			if split {
				l.Runnable, _ = utils.Int[int](r)
				l.Scheduled, _ = utils.Int[int](c)
			}
		}

		if i == 4 {
			l.LastPID, _ = utils.Int[int](field)
		}
	}

	return nil
}

// Meminfo contains memory, swap and cache usage information, read from
// /proc/meminfo. All values are in kilobytes.
//
// References:
//   - proc_meminfo(5)
type Meminfo struct {
	MemTotal          uint64
	MemFree           uint64
	MemAvailable      uint64
	Buffers           uint64
	Cached            uint64
	SwapCached        uint64
	Active            uint64
	Inactive          uint64
	ActiveAnon        uint64
	InactiveAnon      uint64
	ActiveFile        uint64
	InactiveFile      uint64
	Unevictable       uint64
	Mlocked           uint64
	SwapTotal         uint64
	SwapFree          uint64
	Zswap             uint64
	Zswapped          uint64
	Dirty             uint64
	Writeback         uint64
	AnonPages         uint64
	Mapped            uint64
	Shmem             uint64
	KReclaimable      uint64
	Slab              uint64
	SReclaimable      uint64
	SUnreclaim        uint64
	KernelStack       uint64
	PageTables        uint64
	SecPageTables     uint64
	NFSUnstable       uint64
	Bounce            uint64
	WritebackTmp      uint64
	CommitLimit       uint64
	CommittedAS       uint64
	VmallocTotal      uint64
	VmallocUsed       uint64
	VmallocChunk      uint64
	Percpu            uint64
	HardwareCorrupted uint64
	AnonHugePages     uint64
	ShmemHugePages    uint64
	ShmemPmdMapped    uint64
	FileHugePages     uint64
	FilePmdMapped     uint64
	CmaTotal          uint64
	CmaFree           uint64
	Unaccepted        uint64
	Balloon           uint64
	HugePagesTotal    uint64
	HugePagesFree     uint64
	HugePagesRsvd     uint64
	HugePagesSurp     uint64
	Hugepagesize      uint64
	Hugetlb           uint64
	DirectMap4k       uint64
	DirectMap2M       uint64
	DirectMap1G       uint64
}

// GetMeminfo get memory, swap and cache usage information, read from
// /proc/meminfo in the given [Procfs].
func GetMeminfo(proc Procfs) (*Meminfo, error) {
	m := new(Meminfo)

	err := proc.Read("meminfo", m)
	if err != nil {
		return nil, err
	}

	return m, nil
}

// UnmarshalText unmarshals the lines from /proc/meminfo.
func (m *Meminfo) UnmarshalText(b []byte) error {
	for _, line := range utils.Lines(b) {
		key, value, split := utils.Split(line, func(b byte) bool {
			return b == ':'
		})
		if !split {
			continue
		}

		// split value from `kB` suffix.
		value, _, _ = utils.Split(utils.TrimLeft(value), func(b byte) bool {
			return b == ' '
		})

		n, ok := utils.Uint[uint64](value)
		if !ok {
			continue
		}

		switch string(key) {
		case "MemTotal":
			m.MemTotal = n
		case "MemFree":
			m.MemFree = n
		case "MemAvailable":
			m.MemAvailable = n
		case "Buffers":
			m.Buffers = n
		case "Cached":
			m.Cached = n
		case "SwapCached":
			m.SwapCached = n
		case "Active":
			m.Active = n
		case "Inactive":
			m.Inactive = n
		case "Active(anon)":
			m.ActiveAnon = n
		case "Inactive(anon)":
			m.InactiveAnon = n
		case "Active(file)":
			m.ActiveFile = n
		case "Inactive(file)":
			m.InactiveFile = n
		case "Unevictable":
			m.Unevictable = n
		case "Mlocked":
			m.Mlocked = n
		case "SwapTotal":
			m.SwapTotal = n
		case "SwapFree":
			m.SwapFree = n
		case "Zswap":
			m.Zswap = n
		case "Zswapped":
			m.Zswapped = n
		case "Dirty":
			m.Dirty = n
		case "Writeback":
			m.Writeback = n
		case "AnonPages":
			m.AnonPages = n
		case "Mapped":
			m.Mapped = n
		case "Shmem":
			m.Shmem = n
		case "KReclaimable":
			m.KReclaimable = n
		case "Slab":
			m.Slab = n
		case "SReclaimable":
			m.SReclaimable = n
		case "SUnreclaim":
			m.SUnreclaim = n
		case "KernelStack":
			m.KernelStack = n
		case "PageTables":
			m.PageTables = n
		case "SecPageTables":
			m.SecPageTables = n
		case "NFS_Unstable":
			m.NFSUnstable = n
		case "Bounce":
			m.Bounce = n
		case "WritebackTmp":
			m.WritebackTmp = n
		case "CommitLimit":
			m.CommitLimit = n
		case "Committed_AS":
			m.CommittedAS = n
		case "VmallocTotal":
			m.VmallocTotal = n
		case "VmallocUsed":
			m.VmallocUsed = n
		case "VmallocChunk":
			m.VmallocChunk = n
		case "Percpu":
			m.Percpu = n
		case "HardwareCorrupted":
			m.HardwareCorrupted = n
		case "AnonHugePages":
			m.AnonHugePages = n
		case "ShmemHugePages":
			m.ShmemHugePages = n
		case "ShmemPmdMapped":
			m.ShmemPmdMapped = n
		case "FileHugePages":
			m.FileHugePages = n
		case "FilePmdMapped":
			m.FilePmdMapped = n
		case "CmaTotal":
			m.CmaTotal = n
		case "CmaFree":
			m.CmaFree = n
		case "Unaccepted":
			m.Unaccepted = n
		case "Balloon":
			m.Balloon = n
		case "HugePages_Total":
			m.HugePagesTotal = n
		case "HugePages_Free":
			m.HugePagesFree = n
		case "HugePages_Rsvd":
			m.HugePagesRsvd = n
		case "HugePages_Surp":
			m.HugePagesSurp = n
		case "Hugepagesize":
			m.Hugepagesize = n
		case "Hugetlb":
			m.Hugetlb = n
		case "DirectMap4k":
			m.DirectMap4k = n
		case "DirectMap2M":
			m.DirectMap2M = n
		case "DirectMap1G":
			m.DirectMap1G = n
		}
	}

	return nil
}
