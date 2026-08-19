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

// Stat contains status information about a process, read from
// /proc/<pid>/stat.
//
// References:
//   - proc_pid_stat(5)
//   - linux/fs/proc/array.c do_task_stat
type Stat struct {
	Pid                 int
	Command             string
	State               State
	Ppid                int
	Pgrp                int
	Session             int
	TTY                 int
	Tpgid               int
	Flags               uint
	Minflt              uint
	Cminflt             uint
	Majflt              uint
	Cmajflt             uint
	Utime               uint
	Stime               uint
	Cutime              int
	Cstime              int
	Priority            int
	Nice                int
	NumThreads          int
	Itrealvalue         int
	StartTime           uint
	Vsize               uint
	Rss                 int
	Rsslim              uint
	StartCode           uint
	EndCode             uint
	StartStack          uint
	Kstkesp             uint
	Kstkeip             uint
	Signal              uint
	Blocked             uint
	SigIgnore           uint
	SigCatch            uint
	Wchan               uint
	Nswap               uint
	Cnswap              uint
	ExitSignal          int
	Processor           int
	RtPriority          uint
	Policy              uint
	DelayAcctBlkioTicks uint
	GuestTime           uint
	CguestTime          int
	StartData           uint
	EndData             uint
	StartBrk            uint
	ArgStart            uint
	ArgEnd              uint
	EnvStart            uint
	EnvEnd              uint
	ExitCode            int
}

// GetStat reads the specified process id status information, from
// /proc/<pid>/stat in the given [Procfs].
func GetStat(proc Procfs, pid int) (*Stat, error) {
	s := new(Stat)

	err := proc.Read(strconv.Itoa(pid)+"/stat", s)
	if err != nil {
		return nil, err
	}

	return s, nil
}

// UnmarshalText unmarshals the fields from /proc/<pid>/stat.
func (s *Stat) UnmarshalText(b []byte) error {
	for i, field := range utils.Fields(b) {
		// most fields are int/uint, handle special case.
		if i == 1 {
			s.Command = string(field)
			continue
		} else if i == 2 {
			s.State = State(field[0])
			continue
		}

		// parse field as int, allow overflows into uint per-field.
		n, ok := utils.Int[int](field)
		if !ok {
			continue
		}

		switch i {
		case 0:
			s.Pid = n
		case 3:
			s.Ppid = n
		case 4:
			s.Pgrp = n
		case 5:
			s.Session = n
		case 6:
			s.TTY = n
		case 7:
			s.Tpgid = n
		case 8:
			s.Flags, _ = utils.Uint[uint](field)
		case 9:
			s.Minflt = uint(n)
		case 10:
			s.Cminflt = uint(n)
		case 11:
			s.Majflt = uint(n)
		case 12:
			s.Cmajflt = uint(n)
		case 13:
			s.Utime = uint(n)
		case 14:
			s.Stime = uint(n)
		case 15:
			s.Cutime = n
		case 16:
			s.Cstime = n
		case 17:
			s.Priority = n
		case 18:
			s.Nice = n
		case 19:
			s.NumThreads = n
		case 20:
			s.Itrealvalue = n
		case 21:
			s.StartTime = uint(n)
		case 22:
			s.Vsize = uint(n)
		case 23:
			s.Rss = n
		case 24:
			s.Rsslim = uint(n)
		case 25:
			s.StartCode = uint(n)
		case 26:
			s.EndCode = uint(n)
		case 27:
			s.StartStack = uint(n)
		case 28:
			s.Kstkesp = uint(n)
		case 29:
			s.Kstkeip = uint(n)
		case 30:
			s.Signal, _ = utils.Uint[uint](field)
		case 31:
			s.Blocked, _ = utils.Uint[uint](field)
		case 32:
			s.SigIgnore, _ = utils.Uint[uint](field)
		case 33:
			s.SigCatch, _ = utils.Uint[uint](field)
		case 34:
			s.Wchan = uint(n)
		case 35:
			s.Nswap = uint(n)
		case 36:
			s.Cnswap = uint(n)
		case 37:
			s.ExitSignal = n
		case 38:
			s.Processor = n
		case 39:
			s.RtPriority = uint(n)
		case 40:
			s.Policy = uint(n)
		case 41:
			s.DelayAcctBlkioTicks = uint(n)
		case 42:
			s.GuestTime = uint(n)
		case 43:
			s.CguestTime = n
		case 44:
			s.StartData = uint(n)
		case 45:
			s.EndData = uint(n)
		case 46:
			s.StartBrk = uint(n)
		case 47:
			s.ArgStart = uint(n)
		case 48:
			s.ArgEnd = uint(n)
		case 49:
			s.EnvStart = uint(n)
		case 50:
			s.EnvEnd = uint(n)
		case 51:
			s.ExitCode = n
		}
	}

	return nil
}

// State is the current state of a process in [Stat].
type State byte

// Constants for [State].
const (
	Running     State = 'R'
	Sleeping    State = 'S'
	Waiting     State = 'D'
	Zombie      State = 'Z'
	Stopped     State = 'T'
	TracingStop State = 't'
	Dead        State = 'X'
	Idle        State = 'I'
)

func (s State) String() string {
	switch s {
	case Running:
		return "RUNNING"
	case Sleeping:
		return "SLEEPING"
	case Waiting:
		return "WAITING"
	case Zombie:
		return "ZOMBIE"
	case Stopped:
		return "STOPPED"
	case TracingStop:
		return "TRACING_STOP"
	case Dead:
		return "DEAD"
	case Idle:
		return "IDLE"

	default:
		return "UNKNOWN"
	}
}
