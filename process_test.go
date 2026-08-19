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

func TestStat(t *testing.T) {
	proc := From(os.DirFS("testdata/proc"))

	expected := &Stat{
		Pid:        1,
		Command:    "(systemd)",
		State:      Sleeping,
		Pgrp:       1,
		Session:    1,
		Tpgid:      -1,
		Flags:      4194560,
		Minflt:     3830258,
		Cminflt:    948560345,
		Majflt:     108,
		Cmajflt:    4733,
		Utime:      10479,
		Stime:      7170,
		Cutime:     4910029,
		Cstime:     1289825,
		Priority:   20,
		NumThreads: 1,
		StartTime:  5,
		Vsize:      34623488,
		Rss:        5106,
		Rsslim:     18446744073709551615,
		StartCode:  1,
		EndCode:    1,
		Blocked:    671173123,
		SigIgnore:  4096,
		SigCatch:   1260,
		ExitSignal: 17,
		Processor:  2,
	}

	target, err := GetStat(proc, 1)
	if err != nil {
		t.Fatal("unexpected error:", err)
	}

	if !cmp.Equal(expected, target) {
		t.Error(cmp.Diff(expected, target))
	}
}
