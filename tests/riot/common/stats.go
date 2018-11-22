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
	"fmt"
	"os"
	"time"

	chart "github.com/wcharczuk/go-chart"
)

type axis struct {
	x, y []float64
}

func drawChart(m string, a *axis) {
	graph := chart.Chart{
		XAxis: chart.XAxis{
			Name:      "Time",
			NameStyle: chart.StyleShow(),
			Style:     chart.StyleShow(),
		},
		YAxis: chart.YAxis{
			Name:      "Reqests",
			NameStyle: chart.StyleShow(),
			Style:     chart.StyleShow(),
		},
		Series: []chart.Series{
			chart.ContinuousSeries{
				Style: chart.Style{
					Show:        true,
					StrokeColor: chart.GetDefaultColor(0).WithAlpha(64),
					FillColor:   chart.GetDefaultColor(0).WithAlpha(64),
				},

				XValues: a.x,
				YValues: a.y,
			},
		},
	}

	req := fmt.Sprint(NumRequests)
	file, _ := os.Create("results/graph-" + m + "-" + req + ".png")
	defer file.Close()
	_ = graph.Render(chart.PNG, file)

}

func chartsData(a *axis, elapsed, reqs float64) *axis {
	a.x = append(a.x, float64(elapsed))
	a.y = append(a.y, float64(reqs))

	return a
}

func summary(message string, numRequestsf, elapsed float64, c *Config) {

	fmt.Printf(
		"%s throughput: %.0f req/s: (%v reqs in %.3f seconds) | Concurrency: %d\n",
		message,
		numRequestsf/elapsed,
		c.numRequests,
		elapsed,
		c.maxGoRoutines,
	)
}

func summaryPerDuration(message string, numRequestsf, elapsed float64, c *Config) {

	fmt.Printf(
		"%s throughput: %.0f req/s | Concurrency: %d | Elapsed time: %.3f seconds\n",
		message,
		c.counter/elapsed,
		c.maxGoRoutines,
		elapsed,
	)
}

func stats(c *Config, t Task, message string) {
	graph := &axis{}
	ticker := time.NewTicker(1 * time.Second)
	numRequestsf := float64(c.numRequests)
	start := time.Now()
	defer ticker.Stop()
	done := make(chan bool)
	go func() {
		SpawnerOfEvil(c, t)
		elapsed := time.Now().Sub(start).Seconds()
		fmt.Println("Task done.")
		summary(message, numRequestsf, elapsed, c)
		done <- true
	}()
	for {
		select {
		case <-done:
			return
		case t := <-ticker.C:
			_ = t
			elapsed := time.Now().Sub(start).Seconds()
			if charts {
				os.Mkdir("results", 0755)
				go drawChart(message, chartsData(graph, elapsed, c.counter/elapsed))
			}
			summaryPerDuration(message, numRequestsf, elapsed, c)
		}
	}
}
