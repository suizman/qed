/*
   Copyright 2018 Banco Bilbao Vizcaya Argentaria, S.A.

   Licensed under the Apache License, Version 2.0 (the "License");
   you may not use this file except in compliance with the License.
   You may obtain a copy of the License at

       http://www.apache.org/licenses/LICENSE-2.0

   Unless required by applicable law or agreed to in writing, software
   distributed under the License is distributed on an "AS IS" BASIS,
   WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
   See the License for the specific language governing permissions and
   limitations under the License.
*/

package common

import (
	"bufio"
	"flag"
	"fmt"
	"os"
	"strconv"
	"time"
)

var (
	WantAdd          bool
	WantIncremental  bool
	WantMembership   bool
	Offload          bool
	charts           bool
	incrementalDelta int
	Offset           int
	NumRequests      int
	ReadConcurrency  int
	WriteConcurrency int
)

func init() {
	const (
		// Default values
		defaultWantAdd          = false
		defaultWantIncremental  = false
		defaultWantMembership   = false
		defaultOffset           = 0
		defaultIncrementalDelta = 1000
		defaultNumRequests      = 100000

		// Usage
		usage                 = "Benchmark MembershipProof"
		usageDelta            = "Specify delta for the IncrementalProof"
		usageNumRequests      = "Number of requests for the attack"
		usageReadConcurrency  = "Set read concurrency value"
		usageWriteConcurrency = "Set write concurrency value"
		usageOffload          = "Perform reads only on %50 of the cluster size (With cluster size 2 reads will be perofmed only on follower1)"
		usageCharts           = "Create charts while executing the benchmarks. Output: graph-$testname.png"
		usageOffset           = "The starting version from which we start the load"
		usageWantAdd          = "Execute add benchmark"
		usageWantIncremental  = "Execute Incremental benchmark"
	)

	// Create a default config to use as default values in flags
	config := NewDefaultConfig()

	flag.BoolVar(&WantAdd, "add", defaultWantAdd, usageWantAdd)
	flag.BoolVar(&WantMembership, "membership", defaultWantMembership, usage)
	flag.BoolVar(&WantMembership, "m", defaultWantMembership, usage+" (shorthand)")
	flag.BoolVar(&WantIncremental, "incremental", defaultWantIncremental, usageWantIncremental)
	flag.BoolVar(&Offload, "offload", false, usageOffload)
	flag.BoolVar(&charts, "charts", false, usageCharts)
	flag.IntVar(&incrementalDelta, "delta", defaultIncrementalDelta, usageDelta)
	flag.IntVar(&incrementalDelta, "d", defaultIncrementalDelta, usageDelta+" (shorthand)")
	flag.IntVar(&NumRequests, "n", defaultNumRequests, usageNumRequests)
	flag.IntVar(&ReadConcurrency, "r", config.maxGoRoutines, usageReadConcurrency)
	flag.IntVar(&WriteConcurrency, "w", config.maxGoRoutines, usageWriteConcurrency)
	flag.IntVar(&Offset, "offset", defaultOffset, usageOffset)

}

func hotParams(config []*Config) {
	scanner := bufio.NewScanner(os.Stdin)

	for scanner.Scan() {
		value := scanner.Text()

		switch t := value[0:2]; t {
		case "mr":
			i, _ := strconv.ParseInt(value[2:], 10, 64)
			d := time.Duration(i)
			for _, c := range config {
				c.delay_ms = d
			}
			fmt.Printf("Read throughtput set to: %d\n", i)
		case "ir":
			i, _ := strconv.ParseInt(value[2:], 10, 64)
			d := time.Duration(i)
			for _, c := range config {
				c.delay_ms = d
			}
			fmt.Printf("Read throughtput set to: %d\n", i)
		default:
			fmt.Println("Invalid command - Valid commands: mr100|ir200")
		}

	}
}
