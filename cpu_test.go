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

func TestCPUInfo(t *testing.T) {
	tests := []struct {
		path     string
		expected *CPUInfo
	}{
		{
			"testdata/cpuinfo/amd_opteron",
			&CPUInfo{
				Processors: []*Processor{
					{
						ID:       0,
						Vendor:   "AuthenticAMD",
						Model:    "Dual-Core AMD Opteron(tm) Processor 2212",
						Core:     0,
						BogoMIPS: 4003.71,
					},
					{
						ID:       1,
						Vendor:   "AuthenticAMD",
						Model:    "Dual-Core AMD Opteron(tm) Processor 2212",
						Core:     0,
						BogoMIPS: 4000.19,
					},
					{
						ID:       2,
						Vendor:   "AuthenticAMD",
						Model:    "Dual-Core AMD Opteron(tm) Processor 2212",
						Core:     1,
						BogoMIPS: 4000.2,
					},
					{
						ID:       3,
						Vendor:   "AuthenticAMD",
						Model:    "Dual-Core AMD Opteron(tm) Processor 2212",
						Core:     1,
						BogoMIPS: 4000.2,
					},
				},
			},
		},
		{
			"testdata/cpuinfo/amd_ryzen7",
			&CPUInfo{
				Processors: []*Processor{
					{
						ID:       0,
						Vendor:   "AuthenticAMD",
						Model:    "AMD Ryzen 7 2700X Eight-Core Processor",
						Core:     0,
						BogoMIPS: 7386.22,
					},
					{
						ID:       1,
						Vendor:   "AuthenticAMD",
						Model:    "AMD Ryzen 7 2700X Eight-Core Processor",
						Core:     0,
						BogoMIPS: 7386.22,
					},
					{
						ID:       2,
						Vendor:   "AuthenticAMD",
						Model:    "AMD Ryzen 7 2700X Eight-Core Processor",
						Core:     1,
						BogoMIPS: 7386.22,
					},
					{
						ID:       3,
						Vendor:   "AuthenticAMD",
						Model:    "AMD Ryzen 7 2700X Eight-Core Processor",
						Core:     1,
						BogoMIPS: 7386.22,
					},
					{
						ID:       4,
						Vendor:   "AuthenticAMD",
						Model:    "AMD Ryzen 7 2700X Eight-Core Processor",
						Core:     2,
						BogoMIPS: 7386.22,
					},
					{
						ID:       5,
						Vendor:   "AuthenticAMD",
						Model:    "AMD Ryzen 7 2700X Eight-Core Processor",
						Core:     2,
						BogoMIPS: 7386.22,
					},
					{
						ID:       6,
						Vendor:   "AuthenticAMD",
						Model:    "AMD Ryzen 7 2700X Eight-Core Processor",
						Core:     3,
						BogoMIPS: 7386.22,
					},
					{
						ID:       7,
						Vendor:   "AuthenticAMD",
						Model:    "AMD Ryzen 7 2700X Eight-Core Processor",
						Core:     3,
						BogoMIPS: 7386.22,
					},
					{
						ID:       8,
						Vendor:   "AuthenticAMD",
						Model:    "AMD Ryzen 7 2700X Eight-Core Processor",
						Core:     4,
						BogoMIPS: 7386.22,
					},
					{
						ID:       9,
						Vendor:   "AuthenticAMD",
						Model:    "AMD Ryzen 7 2700X Eight-Core Processor",
						Core:     4,
						BogoMIPS: 7386.22,
					},
					{
						ID:       10,
						Vendor:   "AuthenticAMD",
						Model:    "AMD Ryzen 7 2700X Eight-Core Processor",
						Core:     5,
						BogoMIPS: 7386.22,
					},
					{
						ID:       11,
						Vendor:   "AuthenticAMD",
						Model:    "AMD Ryzen 7 2700X Eight-Core Processor",
						Core:     5,
						BogoMIPS: 7386.22,
					},
					{
						ID:       12,
						Vendor:   "AuthenticAMD",
						Model:    "AMD Ryzen 7 2700X Eight-Core Processor",
						Core:     6,
						BogoMIPS: 7386.22,
					},
					{
						ID:       13,
						Vendor:   "AuthenticAMD",
						Model:    "AMD Ryzen 7 2700X Eight-Core Processor",
						Core:     6,
						BogoMIPS: 7386.22,
					},
					{
						ID:       14,
						Vendor:   "AuthenticAMD",
						Model:    "AMD Ryzen 7 2700X Eight-Core Processor",
						Core:     7,
						BogoMIPS: 7386.22,
					},
					{
						ID:       15,
						Vendor:   "AuthenticAMD",
						Model:    "AMD Ryzen 7 2700X Eight-Core Processor",
						Core:     7,
						BogoMIPS: 7386.22,
					},
				},
			},
		},
		{
			"testdata/cpuinfo/apple_m1_vm",
			&CPUInfo{
				Processors: []*Processor{
					{ID: 0, Vendor: "Apple", BogoMIPS: 48},
					{ID: 1, Vendor: "Apple", BogoMIPS: 48},
					{ID: 2, Vendor: "Apple", BogoMIPS: 48},
					{ID: 3, Vendor: "Apple", BogoMIPS: 48},
					{ID: 4, Vendor: "Apple", BogoMIPS: 48},
					{ID: 5, Vendor: "Apple", BogoMIPS: 48},
					{ID: 6, Vendor: "Apple", BogoMIPS: 48},
					{ID: 7, Vendor: "Apple", BogoMIPS: 48},
				},
			},
		},
		{
			"testdata/cpuinfo/intel_xeon",
			&CPUInfo{
				Processors: []*Processor{
					{
						ID:       0,
						Vendor:   "GenuineIntel",
						Model:    "Intel(R) Xeon(R) E-2176G CPU @ 3.70GHz",
						Core:     0,
						BogoMIPS: 7399.7,
					},
					{
						ID:       1,
						Vendor:   "GenuineIntel",
						Model:    "Intel(R) Xeon(R) E-2176G CPU @ 3.70GHz",
						Core:     1,
						BogoMIPS: 7399.7,
					},
					{
						ID:       2,
						Vendor:   "GenuineIntel",
						Model:    "Intel(R) Xeon(R) E-2176G CPU @ 3.70GHz",
						Core:     2,
						BogoMIPS: 7399.7,
					},
					{
						ID:       3,
						Vendor:   "GenuineIntel",
						Model:    "Intel(R) Xeon(R) E-2176G CPU @ 3.70GHz",
						Core:     3,
						BogoMIPS: 7399.7,
					},
					{
						ID:       4,
						Vendor:   "GenuineIntel",
						Model:    "Intel(R) Xeon(R) E-2176G CPU @ 3.70GHz",
						Core:     4,
						BogoMIPS: 7399.7,
					},
					{
						ID:       5,
						Vendor:   "GenuineIntel",
						Model:    "Intel(R) Xeon(R) E-2176G CPU @ 3.70GHz",
						Core:     5,
						BogoMIPS: 7399.7,
					},
					{
						ID:       6,
						Vendor:   "GenuineIntel",
						Model:    "Intel(R) Xeon(R) E-2176G CPU @ 3.70GHz",
						Core:     0,
						BogoMIPS: 7399.7,
					},
					{
						ID:       7,
						Vendor:   "GenuineIntel",
						Model:    "Intel(R) Xeon(R) E-2176G CPU @ 3.70GHz",
						Core:     1,
						BogoMIPS: 7399.7,
					},
					{
						ID:       8,
						Vendor:   "GenuineIntel",
						Model:    "Intel(R) Xeon(R) E-2176G CPU @ 3.70GHz",
						Core:     2,
						BogoMIPS: 7399.7,
					},
					{
						ID:       9,
						Vendor:   "GenuineIntel",
						Model:    "Intel(R) Xeon(R) E-2176G CPU @ 3.70GHz",
						Core:     3,
						BogoMIPS: 7399.7,
					},
					{
						ID:       10,
						Vendor:   "GenuineIntel",
						Model:    "Intel(R) Xeon(R) E-2176G CPU @ 3.70GHz",
						Core:     4,
						BogoMIPS: 7399.7,
					},
					{
						ID:       11,
						Vendor:   "GenuineIntel",
						Model:    "Intel(R) Xeon(R) E-2176G CPU @ 3.70GHz",
						Core:     5,
						BogoMIPS: 7399.7,
					},
				},
			},
		},
		{
			"testdata/cpuinfo/raspi4",
			&CPUInfo{
				Processors: []*Processor{
					{ID: 0, Vendor: "ARM", Model: "Cortex-A72", BogoMIPS: 108},
					{ID: 1, Vendor: "ARM", Model: "Cortex-A72", BogoMIPS: 108},
					{ID: 2, Vendor: "ARM", Model: "Cortex-A72", BogoMIPS: 108},
					{ID: 3, Vendor: "ARM", Model: "Cortex-A72", BogoMIPS: 108},
				},
				Revision: "c03111",
				Serial:   "1234567890123456",
				Model:    "Raspberry Pi 4 Model B Rev 1.1",
				Board: &RaspberryPi{
					Overvolt:     true,
					OTPWrite:     true,
					OTPRead:      true,
					Warranty:     true,
					New:          true,
					Memory:       4096,
					Manufacturer: "Sony UK",
					Processor:    "BCM2711",
					Type:         "4B",
					Revision:     1,
				},
			},
		},
	}

	for _, test := range tests {
		t.Run(test.path, func(t *testing.T) {
			bytes, err := os.ReadFile(test.path)
			if err != nil {
				t.Fatal("could not load testdata:", err)
			}

			target := new(CPUInfo)

			err = target.UnmarshalText(bytes)
			if err != nil {
				t.Fatal("unexpected error:", err)
			}

			if !cmp.Equal(test.expected, target) {
				t.Error(cmp.Diff(test.expected, target))
			}
		})
	}
}

func TestRaspberryPi(t *testing.T) {
	tests := []struct {
		desc     string
		revision string
		expected *RaspberryPi
	}{
		{
			"3B/1.2/1GB",
			"a32082",
			&RaspberryPi{
				Overvolt:     true,
				OTPWrite:     true,
				OTPRead:      true,
				Warranty:     true,
				New:          true,
				Memory:       1024,
				Manufacturer: "Sony Japan",
				Processor:    "BCM2837",
				Type:         "3B",
				Revision:     2,
			},
		},
		{
			"Zero/1.3/512MB",
			"920093",
			&RaspberryPi{
				Overvolt:     true,
				OTPWrite:     true,
				OTPRead:      true,
				Warranty:     true,
				New:          true,
				Memory:       512,
				Manufacturer: "Embest",
				Processor:    "BCM2835",
				Type:         "Zero",
				Revision:     3,
			},
		},
		{
			"4B/1.1/4GB",
			"c03111",
			&RaspberryPi{
				Overvolt:     true,
				OTPWrite:     true,
				OTPRead:      true,
				Warranty:     true,
				New:          true,
				Memory:       4096,
				Manufacturer: "Sony UK",
				Processor:    "BCM2711",
				Type:         "4B",
				Revision:     1,
			},
		},
		{
			"5/1.1/16GB",
			"e04171",
			&RaspberryPi{
				Overvolt:     true,
				OTPWrite:     true,
				OTPRead:      true,
				Warranty:     true,
				New:          true,
				Memory:       16386,
				Manufacturer: "Sony UK",
				Processor:    "BCM2712",
				Type:         "5",
				Revision:     1,
			},
		},
		{
			"CM5 Lite/1.0/2GB",
			"b041a0",
			&RaspberryPi{
				Overvolt:     true,
				OTPWrite:     true,
				OTPRead:      true,
				Warranty:     true,
				New:          true,
				Memory:       2048,
				Manufacturer: "Sony UK",
				Processor:    "BCM2712",
				Type:         "CM5 Lite",
				Revision:     0,
			},
		},
	}

	for _, test := range tests {
		t.Run(test.desc, func(t *testing.T) {
			target, err := parseRaspberryPi(test.revision)
			if err != nil {
				t.Fatal("unexpected error:", err)
			}

			if !cmp.Equal(test.expected, target) {
				t.Error(cmp.Diff(test.expected, target))
			}
		})
	}
}
