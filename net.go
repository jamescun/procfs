// Copyright 2026 James Cunningham
// SPDX-License-Identifier: BSD-3-Clause
//
// Use of this source code is governed by a BSD-style license that can be found
// in the LICENSE file or at https://opensource.org/license/BSD-3-clause

package procfs

import (
	"go.jamescun.com/procfs/internal/utils"
)

// RxStats contains the receive statistics for a [LinkStat].
type RxStats struct {
	Bytes      uint64
	Packets    uint64
	Errors     uint64
	Dropped    uint64
	FIFO       uint64
	Frame      uint64
	Compressed uint64
	Multicast  uint64
}

// TxStats contains the transmit statistics for a [LinkStat].
type TxStats struct {
	Bytes      uint64
	Packets    uint64
	Errors     uint64
	Dropped    uint64
	FIFO       uint64
	Collisions uint64
	Carrier    uint64
	Compressed uint64
}

// LinkStat contains the receive and transmit statistics for a single network
// interface, read from /proc/net/dev.
//
// References:
//   - proc_net(5)
type LinkStat struct {
	Name string
	Rx   RxStats
	Tx   TxStats
}

// UnmarshalText unmarshals a single line from /proc/net/dev.
func (l *LinkStat) UnmarshalText(b []byte) error {
	for i, field := range utils.Fields(b) {
		if i == 0 {
			// name is the only non-numerical field, handle as special case.
			l.Name = string(field[:len(field)-1])
			continue
		}

		n, ok := utils.Uint[uint64](field)
		if !ok {
			continue
		}

		switch i {
		case 1:
			l.Rx.Bytes = n
		case 2:
			l.Rx.Packets = n
		case 3:
			l.Rx.Errors = n
		case 4:
			l.Rx.Dropped = n
		case 5:
			l.Rx.FIFO = n
		case 6:
			l.Rx.Frame = n
		case 7:
			l.Rx.Compressed = n
		case 8:
			l.Rx.Multicast = n

		case 9:
			l.Tx.Bytes = n
		case 10:
			l.Tx.Packets = n
		case 11:
			l.Tx.Errors = n
		case 12:
			l.Tx.Dropped = n
		case 13:
			l.Tx.FIFO = n
		case 14:
			l.Tx.Collisions = n
		case 15:
			l.Tx.Carrier = n
		case 16:
			l.Tx.Compressed = n
		}
	}

	return nil
}

// LinkStats contains the receive and transmit statistics for the systems
// network interfaces, read from /proc/net/dev.
//
// References:
//   - proc_net(5)
type LinkStats []*LinkStat

// GetLinkStats reads receive and transmit statistics for the systems network
// interfaces, read from /proc/net/dev for the given [Procfs].
func GetLinkStats(proc Procfs) (LinkStats, error) {
	ls := LinkStats{}

	err := proc.Read("net/dev", &ls)
	if err != nil {
		return nil, err
	}

	return ls, nil
}

// UnmarshalText unmarshals the lines from /proc/net/dev.
func (ls *LinkStats) UnmarshalText(b []byte) error {
	for _, line := range utils.Lines(b) {
		if utils.HasPrefix(line, "Inter-|") || utils.HasPrefix(line, " face |") {
			// skip headers.
			continue
		}

		l := new(LinkStat)

		err := l.UnmarshalText(line)
		if err != nil {
			return err
		}

		*ls = append(*ls, l)
	}

	return nil
}
