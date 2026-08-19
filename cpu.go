// Copyright 2026 James Cunningham
// SPDX-License-Identifier: BSD-3-Clause
//
// Use of this source code is governed by a BSD-style license that can be found
// in the LICENSE file or at https://opensource.org/license/BSD-3-clause

package procfs

import (
	"fmt"
	"strconv"

	"go.jamescun.com/procfs/internal/utils"
)

// Board is available on some single-board computers (SBC) through
// /proc/cpuinfo detailing their configuration, such as the Raspberry Pi.
type Board interface {
	BoardVendor() string
	BoardVersion() string
}

// Processor contains information about one of the systems processors or
// processor cores (depending on platform), read from /proc/cpuinfo.
//
// Due to the platform-specific nature of the contents of /proc/cpuinfo, this
// is not an extensive implementation of all possible fields, merely a best
// effort to capture the primary information.
type Processor struct {
	ID       int
	Vendor   string
	Model    string
	Core     int
	BogoMIPS float64
}

// CPUInfo contains information about the systems processors or processor cores
// (depending on platform), read from /proc/cpuinfo.
//
// Due to the platform-specific nature of the contents of /proc/cpuinfo, this
// is not an extensive implementation of all possible fields, merely a best
// effort to capture the primary information.
type CPUInfo struct {
	Processors []*Processor
	Revision   string
	Serial     string
	Model      string
	Board      Board
}

// GetCPUInfo reads information about the systems processors or processor cores
// (depending on platform), read from /proc/cpuinfo using the given [Procfs].
func GetCPUInfo(proc Procfs) (*CPUInfo, error) {
	c := new(CPUInfo)

	err := proc.Read("cpuinfo", c)
	if err != nil {
		return nil, err
	}

	return c, nil
}

// UnmarshalText unmarshals the lines from /proc/cpuinfo.
func (c *CPUInfo) UnmarshalText(b []byte) error {
	var working *Processor

	for _, line := range utils.Lines(b) {
		key, value, split := utils.Split(line, func(b byte) bool {
			return b == ':'
		})
		if !split {
			continue
		}

		// trim whitespace around delimiter.
		key = utils.TrimRight(key)
		value = utils.TrimLeft(value)

		switch string(key) {
		case "processor":
			// start of a new processor block
			if working == nil {
				working = new(Processor)
			} else {
				c.Processors = append(c.Processors, working)
				working = new(Processor)
			}

			working.ID, _ = utils.Int[int](value)

		case "bogomips", "BogoMIPS":
			working.BogoMIPS, _ = strconv.ParseFloat(string(value), 64)
		case "core id":
			working.Core, _ = utils.Int[int](value)
		case "CPU implementer":
			if len(value) > 2 {
				code, _ := strconv.ParseUint(string(value[2:]), 16, 64)
				working.Vendor = implementers[code]
			}
		case "CPU part":
			if len(value) > 2 && working.Vendor != "" {
				code, _ := strconv.ParseUint(string(value[2:]), 16, 64)
				switch working.Vendor {
				case "Ampere":
					working.Model = ampereParts[code]
				case "Apple":
					working.Model = appleParts[code]
				case "ARM":
					working.Model = armParts[code]
				case "Microsoft":
					working.Model = microsoftParts[code]
				case "NVIDIA":
					working.Model = nvidiaParts[code]
				case "Qualcomm":
					working.Model = qualcommParts[code]
				}
			}
		case "model name":
			working.Model = string(value)
		case "vendor_id":
			working.Vendor = string(value)

		case "Revision":
			c.Revision = string(value)
		case "Serial":
			c.Serial = string(value)
		case "Model":
			c.Model = string(value)
		}
	}

	if working != nil {
		// ensure last working processor is captured.
		c.Processors = append(c.Processors, working)
	}

	if c.Revision != "" {
		// attempt to parse the board revision as a raspberry pi, as more are
		// added this will need a way to detect which board to parse.
		c.Board, _ = parseRaspberryPi(c.Revision)
	}

	return nil
}

var implementers = map[uint64]string{
	0x41: "ARM",
	0x42: "Broadcom",
	0x43: "Cavium",
	0x44: "DEC",
	0x46: "FUJITSU",
	0x48: "HiSilicon",
	0x49: "Infineon",
	0x4d: "Motorola/Freescale",
	0x4e: "NVIDIA",
	0x50: "APM",
	0x51: "Qualcomm",
	0x53: "Samsung",
	0x56: "Marvell",
	0x61: "Apple",
	0x66: "Faraday",
	0x69: "Intel",
	0x6d: "Microsoft",
	0x70: "Phytium",
	0xc0: "Ampere",
}

var ampereParts = map[uint64]string{
	0xac3: "Ampere-1",
	0xac4: "Ampere-1a",
}

var appleParts = map[uint64]string{
	0x001: "Cyclone",
	0x002: "Typhoon",
	0x003: "Typhoon/Capri",
	0x004: "Twister",
	0x005: "Twister/Elba/Malta",
	0x006: "Hurricane",
	0x007: "Hurricane/Myst",
	0x008: "Monsoon",
	0x009: "Mistral",
	0x00b: "Vortex",
	0x00c: "Tempest",
	0x00f: "Tempest-M9",
	0x010: "Vortex/Aruba",
	0x011: "Tempest/Aruba",
	0x012: "Lightning",
	0x013: "Thunder",
	0x020: "Icestorm-A14",
	0x021: "Firestorm-A14",
	0x022: "Icestorm-M1",
	0x023: "Firestorm-M1",
	0x024: "Icestorm-M1-Pro",
	0x025: "Firestorm-M1-Pro",
	0x026: "Thunder-M10",
	0x028: "Icestorm-M1-Max",
	0x029: "Firestorm-M1-Max",
	0x030: "Blizzard-A15",
	0x031: "Avalanche-A15",
	0x032: "Blizzard-M2",
	0x033: "Avalanche-M2",
	0x034: "Blizzard-M2-Pro",
	0x035: "Avalanche-M2-Pro",
	0x036: "Sawtooth-A16",
	0x037: "Everest-A16",
	0x038: "Blizzard-M2-Max",
	0x039: "Avalanche-M2-Max",
}

var armParts = map[uint64]string{
	0x810: "ARM810",
	0x920: "ARM920",
	0x922: "ARM922",
	0x926: "ARM926",
	0x940: "ARM940",
	0x946: "ARM946",
	0x966: "ARM966",
	0xa20: "ARM1020",
	0xa22: "ARM1022",
	0xa26: "ARM1026",
	0xb02: "ARM11-MPCore",
	0xb36: "ARM1136",
	0xb56: "ARM1156",
	0xb76: "ARM1176",
	0xc05: "Cortex-A5",
	0xc07: "Cortex-A7",
	0xc08: "Cortex-A8",
	0xc09: "Cortex-A9",
	0xc0d: "Cortex-A17",
	0xc0f: "Cortex-A15",
	0xc0e: "Cortex-A17",
	0xc14: "Cortex-R4",
	0xc15: "Cortex-R5",
	0xc17: "Cortex-R7",
	0xc18: "Cortex-R8",
	0xc20: "Cortex-M0",
	0xc21: "Cortex-M1",
	0xc23: "Cortex-M3",
	0xc24: "Cortex-M4",
	0xc27: "Cortex-M7",
	0xc60: "Cortex-M0+",
	0xd01: "Cortex-A32",
	0xd02: "Cortex-A34",
	0xd03: "Cortex-A53",
	0xd04: "Cortex-A35",
	0xd05: "Cortex-A55",
	0xd06: "Cortex-A65",
	0xd07: "Cortex-A57",
	0xd08: "Cortex-A72",
	0xd09: "Cortex-A73",
	0xd0a: "Cortex-A75",
	0xd0b: "Cortex-A76",
	0xd0c: "Neoverse-N1",
	0xd0d: "Cortex-A77",
	0xd0e: "Cortex-A76AE",
	0xd13: "Cortex-R52",
	0xd14: "Cortex-R82AE",
	0xd15: "Cortex-R82",
	0xd16: "Cortex-R52+",
	0xd20: "Cortex-M23",
	0xd21: "Cortex-M33",
	0xd24: "Cortex-M52",
	0xd22: "Cortex-M55",
	0xd23: "Cortex-M85",
	0xd40: "Neoverse-V1",
	0xd41: "Cortex-A78",
	0xd42: "Cortex-A78AE",
	0xd43: "Cortex-A65AE",
	0xd44: "Cortex-X1",
	0xd46: "Cortex-A510",
	0xd47: "Cortex-A710",
	0xd48: "Cortex-X2",
	0xd49: "Neoverse-N2",
	0xd4a: "Neoverse-E1",
	0xd4b: "Cortex-A78C",
	0xd4c: "Cortex-X1C",
	0xd4d: "Cortex-A715",
	0xd4e: "Cortex-X3",
	0xd4f: "Neoverse-V2",
	0xd80: "Cortex-A520",
	0xd81: "Cortex-A720",
	0xd82: "Cortex-X4",
	0xd83: "Neoverse-V3AE",
	0xd84: "Neoverse-V3",
	0xd85: "Cortex-X925",
	0xd87: "Cortex-A725",
	0xd88: "Cortex-A520AE",
	0xd89: "Cortex-A720AE",
	0xd8a: "C1-Nano",
	0xd8b: "C1-Pro",
	0xd8c: "C1-Ultra",
	0xd8e: "Neoverse-N3",
	0xd8f: "Cortex-A320",
	0xd90: "C1-Premium",
}

var microsoftParts = map[uint64]string{
	0xd49: "Azure-Cobalt-100",
}

var nvidiaParts = map[uint64]string{
	0x000: "Denver",
	0x003: "Denver-2",
	0x004: "Carmel",
	0x010: "Olympus",
}

var qualcommParts = map[uint64]string{
	0x001: "Oryon",
	0x00f: "Scorpion",
	0x02d: "Scorpion",
	0x04d: "Krait",
	0x06f: "Krait",
	0x201: "Kryo",
	0x205: "Kryo",
	0x211: "Kryo",
	0x800: "Falkor-V1/Kryo",
	0x801: "Kryo-V2",
	0x802: "Kryo-3XX-Gold",
	0x803: "Kryo-3XX-Silver",
	0x804: "Kryo-4XX-Gold",
	0x805: "Kryo-4XX-Silver",
	0xc00: "Falkor",
	0xc01: "Saphira",
}

// RaspberryPi is a [Board] that contains structured information decoded from a
// "new-style" Raspberry Pi "Revision" within /proc/cpuinfo.
//
// References:
//   - https://www.raspberrypi.com/documentation/computers/raspberry-pi.html#raspberry-pi-revision-codes
type RaspberryPi struct {
	Overvolt     bool
	OTPWrite     bool
	OTPRead      bool
	Warranty     bool
	New          bool
	Memory       int
	Manufacturer string
	Processor    string
	Type         string
	Revision     int
}

// BoardVendor returns "Raspberry Pi".
func (RaspberryPi) BoardVendor() string { return "Raspberry Pi" }

// BoardVersion returns the type of Raspberry Pi.
func (r RaspberryPi) BoardVersion() string { return r.Type }

func parseRaspberryPi(revision string) (*RaspberryPi, error) {
	code, err := strconv.ParseUint(revision, 16, 32)
	if err != nil {
		return nil, fmt.Errorf("invalid revision: %w", err)
	}

	info := &RaspberryPi{
		Overvolt: (code >> 31) == 0,
		OTPWrite: (code>>30)&0x01 == 0,
		OTPRead:  (code>>29)&0x01 == 0,
		Warranty: (code>>25)&0x01 == 0,
		New:      (code>>23)&0x01 == 1,
		Revision: int(code & 0x0F),
	}

	if value := int((code >> 20) & 0x07); value < len(raspPiMemory) {
		info.Memory = raspPiMemory[value]
	}

	if value := int((code >> 16) & 0x0F); value < len(raspPiManufacturer) {
		info.Manufacturer = raspPiManufacturer[value]
	}

	if value := int((code >> 12) & 0x0F); value < len(raspPiProcessor) {
		info.Processor = raspPiProcessor[value]
	}

	if value := int((code >> 4) & 0xFF); value < len(raspPiType) {
		info.Type = raspPiType[value]
	}

	return info, nil
}

var raspPiMemory = []int{
	256, 512, 1024, 2048, 4096, 8192, 16386,
}

var raspPiManufacturer = []string{
	"Sony UK", "Egoman", "Embest", "Sony Japan", "Embest", "Stadium",
}

var raspPiProcessor = []string{
	"BCM2835", "BCM2836", "BCM2837", "BCM2711", "BCM2712",
}

var raspPiType = []string{
	"A", "B", "A+", "B+", "2B", "Alpha", "CM1", "", "3B", "Zero", "CM3", "",
	"Zero W", "3B+", "3A+", "Internal", "CM3+", "4B", "Zero 2 W", "400", "CM4",
	"CM4S", "Internal", "5", "CM5", "500/500+", "CM5 Lite", "CM0",
}
