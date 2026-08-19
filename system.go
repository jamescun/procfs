// Copyright 2026 James Cunningham
// SPDX-License-Identifier: BSD-3-Clause
//
// Use of this source code is governed by a BSD-style license that can be found
// in the LICENSE file or at https://opensource.org/license/BSD-3-clause

package procfs

import (
	"strconv"
	"strings"
	"time"

	"go.jamescun.com/procfs/internal/utils"
)

// Cgroup contains details of a single cgroup subsystem, read from
// /proc/cgroups.
//
// References:
//   - proc_groups
type Cgroup struct {
	Name       string
	Hierarchy  int
	NumCgroups int
	Enabled    bool
}

// UnmarshalText unmarshals a single line from /proc/cgroups.
func (c *Cgroup) UnmarshalText(b []byte) error {
	for i, field := range utils.Fields(b) {
		switch i {
		case 0:
			c.Name = string(field)
		case 1:
			c.Hierarchy, _ = utils.Int[int](field)
		case 2:
			c.NumCgroups, _ = utils.Int[int](field)
		case 3:
			enabled, _ := utils.Int[int](field)
			if enabled != 0 {
				c.Enabled = true
			}
		}
	}

	return nil
}

// Cgroups contains the details of the systems cgroup subsystems, read from
// /proc/cgroups.
//
// References:
//   - proc_groups
type Cgroups []*Cgroup

// GetCgroups reads the systems cgroup subsystems, read from /proc/cgroups in
// the given [Procfs].
func GetCgroups(proc Procfs) (Cgroups, error) {
	cs := Cgroups{}

	err := proc.Read("cgroups", &cs)
	if err != nil {
		return nil, err
	}

	return cs, nil
}

// UnmarshalText unmarshals the lines from /proc/cgroups.
func (cs *Cgroups) UnmarshalText(b []byte) error {
	for i, line := range utils.Lines(b) {
		if i == 0 {
			// skip header.
			continue
		}

		c := new(Cgroup)

		err := c.UnmarshalText(line)
		if err != nil {
			return err
		}

		*cs = append(*cs, c)
	}

	return nil
}

// CPUStat is either the sum processor time spent or time spent for a single
// processor in [CPUStats], read from /proc/stat.
//
// References:
//   - proc_stat(5)
type CPUStat struct {
	CPU       int
	User      uint64
	Nice      uint64
	System    uint64
	Idle      uint64
	IOWait    uint64
	IRQ       uint64
	SoftIRQ   uint64
	Steal     uint64
	Guest     uint64
	GuestNice uint64
}

// UnmarshalText unmarshals a single processor line from /proc/stat.
func (c *CPUStat) UnmarshalText(b []byte) error {
	for i, field := range utils.Fields(b) {
		if i == 0 {
			// total cpu value does not have a CPU number.
			if !utils.Equal(field, "cpu") {
				c.CPU, _ = utils.Int[int](field[3:])
			}

			continue
		}

		n, ok := utils.Uint[uint64](field)
		if !ok {
			continue
		}

		switch i {
		case 1:
			c.User = n
		case 2:
			c.Nice = n
		case 3:
			c.System = n
		case 4:
			c.Idle = n
		case 5:
			c.IOWait = n
		case 6:
			c.IRQ = n
		case 7:
			c.SoftIRQ = n
		case 8:
			c.Steal = n
		case 9:
			c.Guest = n
		case 10:
			c.GuestNice = n
		}
	}

	return nil
}

// CPUStats contains statistics about processor time and running processes,
// read from /proc/stat.
//
// References:
//   - proc_stat(5)
type CPUStats struct {
	Total        CPUStat
	CPU          []*CPUStat
	Ctxt         uint64
	Btime        uint64
	Processes    uint64
	ProcsRunning uint64
	ProcsBlocked uint64
}

// GetCPUStats reads statistics about processor time and running processed,
// read from /proc/stat in the given [Procfs].
func GetCPUStats(proc Procfs) (*CPUStats, error) {
	cs := new(CPUStats)

	err := proc.Read("stat", cs)
	if err != nil {
		return nil, err
	}

	return cs, nil
}

// UnmarshalText unmarshals the lines from /proc/stat.
func (cs *CPUStats) UnmarshalText(b []byte) error {
	for i, line := range utils.Lines(b) {
		if i == 0 {
			// first line is always the total CPUStat.
			err := cs.Total.UnmarshalText(line)
			if err != nil {
				return err
			}

			continue
		}

		key, value, split := utils.Split(line, func(b byte) bool {
			return b == ' '
		})
		if !split {
			continue
		}

		if utils.HasPrefix(key, "cpu") {
			// per-processor CPUStat entry.
			c := new(CPUStat)

			err := c.UnmarshalText(line)
			if err != nil {
				return err
			}

			cs.CPU = append(cs.CPU, c)
			continue
		}

		// everything else is a uint64.
		n, ok := utils.Uint[uint64](value)
		if !ok {
			continue
		}

		switch string(key) {
		case "ctxt":
			cs.Ctxt = n
		case "btime":
			cs.Btime = n
		case "processes":
			cs.Processes = n
		case "procs_running":
			cs.ProcsRunning = n
		case "procs_blocked":
			cs.ProcsBlocked = n
		}
	}

	return nil
}

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

// Module is a currently loaded Linux Kernel module, read from /proc/modules.
//
// References:
//   - proc_modules(5)
type Module struct {
	Name      string
	Size      uint64
	Instances int
	Depends   []string
	State     string
}

// UnmarshalText unmarshals a single line from /proc/modules.
func (m *Module) UnmarshalText(b []byte) error {
	// the contents are mostly string based, cast to a string and cut that.
	text := string(b)

	for i, field := range utils.Fields(text) {
		switch i {
		case 0:
			m.Name = field
		case 1:
			m.Size, _ = utils.Uint[uint64](field)
		case 2:
			m.Instances, _ = utils.Int[int](field)
		case 3:
			if field != "-" {
				m.Depends = strings.Split(field[:len(field)-1], ",")
			}
		case 4:
			m.State = field
		}
	}

	return nil
}

// Modules are the currently loaded Linux Kernel modules, read from
// /proc/modules.
//
// References:
//   - proc_modules(5)
type Modules []*Module

// GetModules reads the currently loaded Linux Kernel modules, from
// /proc/modules in the given [Procfs].
func GetModules(proc Procfs) (Modules, error) {
	ms := Modules{}

	err := proc.Read("modules", &ms)
	if err != nil {
		return nil, err
	}

	return ms, nil
}

// UnmarshalText unmarshals the lines from /proc/modules.
func (ms *Modules) UnmarshalText(b []byte) error {
	for _, line := range utils.Lines(b) {
		m := new(Module)

		err := m.UnmarshalText(line)
		if err != nil {
			return err
		}

		*ms = append(*ms, m)
	}

	return nil
}

// Swap is one of the swap devices configured, read from /proc/swaps.
//
// References:
//   - proc_swaps(5)
//   - linux/mm/swapfile.c swap_show
type Swap struct {
	Name     string
	Type     string
	Size     uint64
	Used     uint64
	Priority int
}

// UnmarshalText unmarshals a single line from /proc/swaps.
func (s *Swap) UnmarshalText(b []byte) error {
	// the contents are mostly string based, cast to a string and cut that.
	text := string(b)

	for i, field := range utils.Fields(text) {
		switch i {
		case 0:
			s.Name = field
		case 1:
			s.Type = field
		case 2:
			s.Size, _ = utils.Uint[uint64](field)
		case 3:
			s.Used, _ = utils.Uint[uint64](field)
		case 4:
			s.Priority, _ = utils.Int[int](field)
		}
	}

	return nil
}

// Swaps are the swap devices configured for the system, read from /proc/swaps.
//
// References:
//   - proc_swaps(5)
//   - linux/mm/swapfile.c swap_show
type Swaps []*Swap

// GetSwaps reads the swap devices configured for the system, read from
// proc/swaps in the given [Procfs].
func GetSwaps(proc Procfs) (Swaps, error) {
	ss := Swaps{}

	err := proc.Read("swaps", &ss)
	if err != nil {
		return nil, err
	}

	return ss, nil
}

// UnmarshalText unmarshals the lines from /proc/swaps.
func (ss *Swaps) UnmarshalText(b []byte) error {
	for i, line := range utils.Lines(b) {
		if i == 0 {
			// skip header.
			continue
		}

		s := new(Swap)
		err := s.UnmarshalText(line)
		if err != nil {
			return err
		}

		*ss = append(*ss, s)
	}

	return nil
}

// Uptime contains system uptime information, read from /proc/uptime.
//
// References:
//   - proc_uptime(5)
type Uptime struct {
	Up   time.Duration
	Idle time.Duration
}

// GetUptime reads system uptime information, read from /proc/uptime in the
// given [Procfs].
func GetUptime(proc Procfs) (*Uptime, error) {
	u := new(Uptime)

	err := proc.Read("uptime", u)
	if err != nil {
		return nil, err
	}

	return u, nil
}

// UnmarshalText unmarshals the line from /proc/uptime.
func (u *Uptime) UnmarshalText(b []byte) error {
	for i, field := range utils.Fields(b) {
		switch i {
		case 0:
			field, _, _ = utils.Split(field, func(b byte) bool {
				return b == '.'
			})

			value, _ := utils.Int[time.Duration](field)
			u.Up = value * time.Second

		case 1:
			field, _, _ = utils.Split(field, func(b byte) bool {
				return b == '.'
			})

			value, _ := utils.Int[time.Duration](field)
			u.Idle = value * time.Second
		}
	}

	return nil
}

// VMStat contains virtual memory statistics, read from /proc/vmstat.
//
// Some values are absolute measures, and some are incrementing counters.
//
// References:
//   - proc_vmstat(5)
type VMStat struct {
	NrFreePages                 uint64
	NrFreePagesBlocks           uint64
	NrZoneInactiveAnon          uint64
	NrZoneActiveAnon            uint64
	NrZoneInactiveFile          uint64
	NrZoneActiveFile            uint64
	NrZoneUnevictable           uint64
	NrZoneWritePending          uint64
	NrMlock                     uint64
	NrZsPages                   uint64
	NrFreeCma                   uint64
	NrInactiveAnon              uint64
	NrActiveAnon                uint64
	NrInactiveFile              uint64
	NrActiveFile                uint64
	NrUnevictable               uint64
	NrSlabReclaimable           uint64
	NrSlabUnreclaimable         uint64
	NrIsolatedAnon              uint64
	NrIsolatedFile              uint64
	WorkingSetNodes             uint64
	WorkingSetRefaultAnon       uint64
	WorkingSetRefaultFile       uint64
	WorkingSetActivateAnon      uint64
	WorkingSetActivateFile      uint64
	WorkingSetRestoreAnon       uint64
	WorkingSetRestoreFile       uint64
	WorkingSetNodeReclaim       uint64
	NrAnonPages                 uint64
	NrMapped                    uint64
	NrFilePages                 uint64
	NrDirty                     uint64
	NrWriteback                 uint64
	NrShmem                     uint64
	NrShmemHugePages            uint64
	NrShmemPmdMapped            uint64
	NrFileHugePages             uint64
	NrFilePmdMapped             uint64
	NrAnonTransparentHugePages  uint64
	NrVmscanWrite               uint64
	NrVmscanImmediateReclaim    uint64
	NrDirtied                   uint64
	NrWritten                   uint64
	NrThrottledWritten          uint64
	NrKernelMiscReclaimable     uint64
	NrFollPinAcquired           uint64
	NrFollPinReleased           uint64
	NrKernelStack               uint64
	NrPageTablePages            uint64
	NrSecPageTablePages         uint64
	NrSwapCached                uint64
	PgDemoteKswapd              uint64
	PgDemoteDirect              uint64
	PgDemoteKHugePaged          uint64
	PgDemoteProactive           uint64
	NrBalloonPages              uint64
	NrKernelFilePages           uint64
	NrDirtyThreshold            uint64
	NrDirtyBackgroundThreshold  uint64
	NrMemmapPages               uint64
	NrMemmapBootPages           uint64
	Pgpgin                      uint64
	Pgpgout                     uint64
	Pswpin                      uint64
	Pswpout                     uint64
	PgAllocDMA                  uint64
	PgAllocDMA32                uint64
	PgAllocNormal               uint64
	PgAllocMovable              uint64
	AllocStallDMA               uint64
	AllocStallDMA32             uint64
	AllocStallNormal            uint64
	AllocStallMovable           uint64
	PgSkipDMA                   uint64
	PgSkipDMA32                 uint64
	PgSkipNormal                uint64
	PgSkipMovable               uint64
	PgFree                      uint64
	PgActivate                  uint64
	PgDeactivate                uint64
	PgLazyFree                  uint64
	PgFault                     uint64
	PgMajFault                  uint64
	PgLazyFreed                 uint64
	PgRefill                    uint64
	PgReuse                     uint64
	PgStealKswapd               uint64
	PgStealDirect               uint64
	PgStealKHugePaged           uint64
	PgStealProactive            uint64
	PgScanKswapd                uint64
	PgScanDirect                uint64
	PgScanKHugePaged            uint64
	PgScanProactive             uint64
	PgScanDirectThrottle        uint64
	PgScanAnon                  uint64
	PgScanFile                  uint64
	PgStealAnon                 uint64
	PgStealFile                 uint64
	PgInodeSteal                uint64
	SlabsScanned                uint64
	KswapdInodeSteal            uint64
	KswapdLowWmarkHitQuickly    uint64
	KswapdHighWmarkHitQuickly   uint64
	PageOutRun                  uint64
	PgRotated                   uint64
	DropPageCache               uint64
	DropSlab                    uint64
	OomKill                     uint64
	PgMigrateSuccess            uint64
	PgMigrateFail               uint64
	ThpMigrationSuccess         uint64
	ThpMigrationFail            uint64
	ThpMigrationSplit           uint64
	CompactMigrateScanned       uint64
	CompactFreeScanned          uint64
	CompactIsolated             uint64
	CompactStall                uint64
	CompactFail                 uint64
	CompactSuccess              uint64
	CompactDaemonWake           uint64
	CompactDaemonMigrateScanned uint64
	CompactDaemonFreeScanned    uint64
	UnevictablePgsCulled        uint64
	UnevictablePgsScanned       uint64
	UnevictablePgsRescued       uint64
	UnevictablePgsMlocked       uint64
	UnevictablePgsMunlocked     uint64
	UnevictablePgsCleared       uint64
	UnevictablePgsStranded      uint64
	ThpFaultAlloc               uint64
	ThpFaultFallback            uint64
	ThpFaultFallbackCharge      uint64
	ThpCollapseAlloc            uint64
	ThpCollapseAllocFailed      uint64
	ThpFileAlloc                uint64
	ThpFileFallback             uint64
	ThpFileFallbackCharge       uint64
	ThpFileMapped               uint64
	ThpSplitPage                uint64
	ThpSplitPageFailed          uint64
	ThpDeferredSplitPage        uint64
	ThpUnderusedSplitPage       uint64
	ThpSplitPmd                 uint64
	ThpScanExceedNonePte        uint64
	ThpScanExceedSwapPte        uint64
	ThpScanExceedSharePte       uint64
	ThpZeroPageAlloc            uint64
	ThpZeroPageAllocFailed      uint64
	ThpSwpout                   uint64
	ThpSwpoutFallback           uint64
	BalloonInflate              uint64
	BalloonDeflate              uint64
	BalloonMigrate              uint64
	SwapRa                      uint64
	SwapRaHit                   uint64
	SwpinZero                   uint64
	SwpoutZero                  uint64
	NrUnstable                  uint64
}

// GetVMStat gets virtual memory statistics about the system, read from
// /proc/vmstat in the given [Procfs].
func GetVMStat(proc Procfs) (*VMStat, error) {
	vm := &VMStat{}

	err := proc.Read("vmstat", vm)
	if err != nil {
		return nil, err
	}

	return vm, nil
}

// UnmarshalText unmarshals the lines from /proc/vmstat.
func (v *VMStat) UnmarshalText(b []byte) error {
	for _, line := range utils.Lines(b) {
		key, value, split := utils.Split(line, func(b byte) bool {
			return b == ' '
		})
		if !split {
			continue
		}

		n, ok := utils.Uint[uint64](value)
		if !ok {
			continue
		}

		switch string(key) {
		case "nr_free_pages":
			v.NrFreePages = n
		case "nr_free_pages_blocks":
			v.NrFreePagesBlocks = n
		case "nr_zone_inactive_anon":
			v.NrZoneInactiveAnon = n
		case "nr_zone_active_anon":
			v.NrZoneActiveAnon = n
		case "nr_zone_inactive_file":
			v.NrZoneInactiveFile = n
		case "nr_zone_active_file":
			v.NrZoneActiveFile = n
		case "nr_zone_unevictable":
			v.NrZoneUnevictable = n
		case "nr_zone_write_pending":
			v.NrZoneWritePending = n
		case "nr_mlock":
			v.NrMlock = n
		case "nr_zspages":
			v.NrZsPages = n
		case "nr_free_cma":
			v.NrFreeCma = n
		case "nr_inactive_anon":
			v.NrInactiveAnon = n
		case "nr_active_anon":
			v.NrActiveAnon = n
		case "nr_inactive_file":
			v.NrInactiveFile = n
		case "nr_active_file":
			v.NrActiveFile = n
		case "nr_unevictable":
			v.NrUnevictable = n
		case "nr_slab_reclaimable":
			v.NrSlabReclaimable = n
		case "nr_slab_unreclaimable":
			v.NrSlabUnreclaimable = n
		case "nr_isolated_anon":
			v.NrIsolatedAnon = n
		case "nr_isolated_file":
			v.NrIsolatedFile = n
		case "workingset_nodes":
			v.WorkingSetNodes = n
		case "workingset_refault_anon":
			v.WorkingSetRefaultAnon = n
		case "workingset_refault_file":
			v.WorkingSetRefaultFile = n
		case "workingset_activate_anon":
			v.WorkingSetActivateAnon = n
		case "workingset_activate_file":
			v.WorkingSetActivateFile = n
		case "workingset_restore_anon":
			v.WorkingSetRestoreAnon = n
		case "workingset_restore_file":
			v.WorkingSetRestoreFile = n
		case "workingset_nodereclaim":
			v.WorkingSetNodeReclaim = n
		case "nr_anon_pages":
			v.NrAnonPages = n
		case "nr_mapped":
			v.NrMapped = n
		case "nr_file_pages":
			v.NrFilePages = n
		case "nr_dirty":
			v.NrDirty = n
		case "nr_writeback":
			v.NrWriteback = n
		case "nr_shmem":
			v.NrShmem = n
		case "nr_shmem_hugepages":
			v.NrShmemHugePages = n
		case "nr_shmem_pmdmapped":
			v.NrShmemPmdMapped = n
		case "nr_file_hugepages":
			v.NrFileHugePages = n
		case "nr_file_pmdmapped":
			v.NrFilePmdMapped = n
		case "nr_anon_transparent_hugepages":
			v.NrAnonTransparentHugePages = n
		case "nr_vmscan_write":
			v.NrVmscanWrite = n
		case "nr_vmscan_immediate_reclaim":
			v.NrVmscanImmediateReclaim = n
		case "nr_dirtied":
			v.NrDirtied = n
		case "nr_written":
			v.NrWritten = n
		case "nr_throttled_written":
			v.NrThrottledWritten = n
		case "nr_kernel_misc_reclaimable":
			v.NrKernelMiscReclaimable = n
		case "nr_foll_pin_acquired":
			v.NrFollPinAcquired = n
		case "nr_foll_pin_released":
			v.NrFollPinReleased = n
		case "nr_kernel_stack":
			v.NrKernelStack = n
		case "nr_page_table_pages":
			v.NrPageTablePages = n
		case "nr_sec_page_table_pages":
			v.NrSecPageTablePages = n
		case "nr_swapcached":
			v.NrSwapCached = n
		case "pgdemote_kswapd":
			v.PgDemoteKswapd = n
		case "pgdemote_direct":
			v.PgDemoteDirect = n
		case "pgdemote_khugepaged":
			v.PgDemoteKHugePaged = n
		case "pgdemote_proactive":
			v.PgDemoteProactive = n
		case "nr_balloon_pages":
			v.NrBalloonPages = n
		case "nr_kernel_file_pages":
			v.NrKernelFilePages = n
		case "nr_dirty_threshold":
			v.NrDirtyThreshold = n
		case "nr_dirty_background_threshold":
			v.NrDirtyBackgroundThreshold = n
		case "nr_memmap_pages":
			v.NrMemmapPages = n
		case "nr_memmap_boot_pages":
			v.NrMemmapBootPages = n
		case "pgpgin":
			v.Pgpgin = n
		case "pgpgout":
			v.Pgpgout = n
		case "pswpin":
			v.Pswpin = n
		case "pswpout":
			v.Pswpout = n
		case "pgalloc_dma":
			v.PgAllocDMA = n
		case "pgalloc_dma32":
			v.PgAllocDMA32 = n
		case "pgalloc_normal":
			v.PgAllocNormal = n
		case "pgalloc_movable":
			v.PgAllocMovable = n
		case "allocstall_dma":
			v.AllocStallDMA = n
		case "allocstall_dma32":
			v.AllocStallDMA32 = n
		case "allocstall_normal":
			v.AllocStallNormal = n
		case "allocstall_movable":
			v.AllocStallMovable = n
		case "pgskip_dma":
			v.PgSkipDMA = n
		case "pgskip_dma32":
			v.PgSkipDMA32 = n
		case "pgskip_normal":
			v.PgSkipNormal = n
		case "pgskip_movable":
			v.PgSkipMovable = n
		case "pgfree":
			v.PgFree = n
		case "pgactivate":
			v.PgActivate = n
		case "pgdeactivate":
			v.PgDeactivate = n
		case "pglazyfree":
			v.PgLazyFree = n
		case "pgfault":
			v.PgFault = n
		case "pgmajfault":
			v.PgMajFault = n
		case "pglazyfreed":
			v.PgLazyFreed = n
		case "pgrefill":
			v.PgRefill = n
		case "pgreuse":
			v.PgReuse = n
		case "pgsteal_kswapd":
			v.PgStealKswapd = n
		case "pgsteal_direct":
			v.PgStealDirect = n
		case "pgsteal_khugepaged":
			v.PgStealKHugePaged = n
		case "pgsteal_proactive":
			v.PgStealProactive = n
		case "pgscan_kswapd":
			v.PgScanKswapd = n
		case "pgscan_direct":
			v.PgScanDirect = n
		case "pgscan_khugepaged":
			v.PgScanKHugePaged = n
		case "pgscan_proactive":
			v.PgScanProactive = n
		case "pgscan_direct_throttle":
			v.PgScanDirectThrottle = n
		case "pgscan_anon":
			v.PgScanAnon = n
		case "pgscan_file":
			v.PgScanFile = n
		case "pgsteal_anon":
			v.PgStealAnon = n
		case "pgsteal_file":
			v.PgStealFile = n
		case "pginodesteal":
			v.PgInodeSteal = n
		case "slabs_scanned":
			v.SlabsScanned = n
		case "kswapd_inodesteal":
			v.KswapdInodeSteal = n
		case "kswapd_low_wmark_hit_quickly":
			v.KswapdLowWmarkHitQuickly = n
		case "kswapd_high_wmark_hit_quickly":
			v.KswapdHighWmarkHitQuickly = n
		case "pageoutrun":
			v.PageOutRun = n
		case "pgrotated":
			v.PgRotated = n
		case "drop_pagecache":
			v.DropPageCache = n
		case "drop_slab":
			v.DropSlab = n
		case "oom_kill":
			v.OomKill = n
		case "pgmigrate_success":
			v.PgMigrateSuccess = n
		case "pgmigrate_fail":
			v.PgMigrateFail = n
		case "thp_migration_success":
			v.ThpMigrationSuccess = n
		case "thp_migration_fail":
			v.ThpMigrationFail = n
		case "thp_migration_split":
			v.ThpMigrationSplit = n
		case "compact_migrate_scanned":
			v.CompactMigrateScanned = n
		case "compact_free_scanned":
			v.CompactFreeScanned = n
		case "compact_isolated":
			v.CompactIsolated = n
		case "compact_stall":
			v.CompactStall = n
		case "compact_fail":
			v.CompactFail = n
		case "compact_success":
			v.CompactSuccess = n
		case "compact_daemon_wake":
			v.CompactDaemonWake = n
		case "compact_daemon_migrate_scanned":
			v.CompactDaemonMigrateScanned = n
		case "compact_daemon_free_scanned":
			v.CompactDaemonFreeScanned = n
		case "unevictable_pgs_culled":
			v.UnevictablePgsCulled = n
		case "unevictable_pgs_scanned":
			v.UnevictablePgsScanned = n
		case "unevictable_pgs_rescued":
			v.UnevictablePgsRescued = n
		case "unevictable_pgs_mlocked":
			v.UnevictablePgsMlocked = n
		case "unevictable_pgs_munlocked":
			v.UnevictablePgsMunlocked = n
		case "unevictable_pgs_cleared":
			v.UnevictablePgsCleared = n
		case "unevictable_pgs_stranded":
			v.UnevictablePgsStranded = n
		case "thp_fault_alloc":
			v.ThpFaultAlloc = n
		case "thp_fault_fallback":
			v.ThpFaultFallback = n
		case "thp_fault_fallback_charge":
			v.ThpFaultFallbackCharge = n
		case "thp_collapse_alloc":
			v.ThpCollapseAlloc = n
		case "thp_collapse_alloc_failed":
			v.ThpCollapseAllocFailed = n
		case "thp_file_alloc":
			v.ThpFileAlloc = n
		case "thp_file_fallback":
			v.ThpFileFallback = n
		case "thp_file_fallback_charge":
			v.ThpFileFallbackCharge = n
		case "thp_file_mapped":
			v.ThpFileMapped = n
		case "thp_split_page":
			v.ThpSplitPage = n
		case "thp_split_page_failed":
			v.ThpSplitPageFailed = n
		case "thp_deferred_split_page":
			v.ThpDeferredSplitPage = n
		case "thp_underused_split_page":
			v.ThpUnderusedSplitPage = n
		case "thp_split_pmd":
			v.ThpSplitPmd = n
		case "thp_scan_exceed_none_pte":
			v.ThpScanExceedNonePte = n
		case "thp_scan_exceed_swap_pte":
			v.ThpScanExceedSwapPte = n
		case "thp_scan_exceed_share_pte":
			v.ThpScanExceedSharePte = n
		case "thp_zero_page_alloc":
			v.ThpZeroPageAlloc = n
		case "thp_zero_page_alloc_failed":
			v.ThpZeroPageAllocFailed = n
		case "thp_swpout":
			v.ThpSwpout = n
		case "thp_swpout_fallback":
			v.ThpSwpoutFallback = n
		case "balloon_inflate":
			v.BalloonInflate = n
		case "balloon_deflate":
			v.BalloonDeflate = n
		case "balloon_migrate":
			v.BalloonMigrate = n
		case "swap_ra":
			v.SwapRa = n
		case "swap_ra_hit":
			v.SwapRaHit = n
		case "swpin_zero":
			v.SwpinZero = n
		case "swpout_zero":
			v.SwpoutZero = n
		case "nr_unstable":
			v.NrUnstable = n
		}
	}

	return nil
}
