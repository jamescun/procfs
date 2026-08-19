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
	for i, line := range utils.Lines(b) {
		if i < 2 {
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

// WirelessStat contains statistics about a single wireless network interface,
// read from /proc/net/wireless.
//
// References:
//   - linux/net/wireless/wext-proc.c
type WirelessStat struct {
	Name    string
	Status  uint16
	Quality int
	Level   int
	Noise   int
	Nwid    uint64
	Crypt   uint64
	Frag    uint64
	Retry   uint64
	Misc    uint64
	Beacon  uint64
}

// UnmarshalText unmarshals a single line from /proc/net/wireless.
func (w *WirelessStat) UnmarshalText(b []byte) error {
	for i, field := range utils.Fields(b) {
		if i == 0 {
			w.Name = string(field[:len(field)-1])
			continue
		}

		if i == 1 {
			w.Status, _ = utils.Uint[uint16](field)
			continue
		}

		if i >= 2 && i <= 4 {
			// strip '.' updated marker if present.
			if field[len(field)-1] == '.' {
				field = field[:len(field)-1]
			}

			n, _ := utils.Int[int](field)

			switch i {
			case 2:
				w.Quality = n
			case 3:
				w.Level = n
			case 4:
				w.Noise = n
			}

			continue
		}

		n, ok := utils.Uint[uint64](field)
		if !ok {
			continue
		}

		switch i {
		case 5:
			w.Nwid = n
		case 6:
			w.Crypt = n
		case 7:
			w.Frag = n
		case 8:
			w.Retry = n
		case 9:
			w.Misc = n
		case 10:
			w.Beacon = n
		}
	}

	return nil
}

// WirelessStats contains statistics about the systems wireless network
// interfaces, read from /proc/net/wireless.
//
// References:
//   - linux/net/wireless/wext-proc.c
type WirelessStats []*WirelessStat

// GetWirelessStats reads statistics about the systems wireless network
// interfaces, read from /proc/net/wireless in the given [Procfs].
func GetWirelessStats(proc Procfs) (WirelessStats, error) {
	ws := WirelessStats{}

	err := proc.Read("net/wireless", &ws)
	if err != nil {
		return nil, err
	}

	return ws, nil
}

// UnmarshalText unmarshals the lines from /proc/net/wireless.
func (ws *WirelessStats) UnmarshalText(b []byte) error {
	for i, line := range utils.Lines(b) {
		if i < 2 {
			// skip headers.
			continue
		}

		w := new(WirelessStat)

		err := w.UnmarshalText(line)
		if err != nil {
			return err
		}

		*ws = append(*ws, w)
	}

	return nil
}
