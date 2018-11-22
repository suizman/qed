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

package main

import (
	"flag"
	"fmt"
	"os"

	c "github.com/bbva/qed/tests/riot/common"
)

func main() {
	var n int
	switch m := os.Getenv("CLUSTER_SIZE"); m {
	case "":
		n = 0
	case "2":
		n = 2
	case "4":
		n = 4
	default:
		fmt.Println("Error: CLUSTER_SIZE env var should have values 2 or 4, or not be defined at all.")
	}

	flag.Parse()

	if c.Offload {
		n = n / 2
		fmt.Printf("Offload: %v | %d\n", c.Offload, n)
	}

	if c.WantAdd {
		fmt.Print("Benchmark ADD")
		c.BenchmarkAdd(n, c.NumRequests, c.ReadConcurrency, c.WriteConcurrency, c.Offset)
	}

	if c.WantMembership {
		fmt.Print("Benchmark MEMBERSHIP")
		c.BenchmarkMembership(n, c.NumRequests, c.ReadConcurrency, c.WriteConcurrency)
	}

	if c.WantIncremental {
		fmt.Print("Benchmark INCREMENTAL")
		c.BenchmarkIncremental(n, c.NumRequests, c.ReadConcurrency, c.WriteConcurrency)
	}
}
