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

package e2e

import (
	"fmt"
	"net/http"
	"os"
	"os/user"
	"strings"
	"testing"
	"time"

	"github.com/bbva/qed/client"
	"github.com/bbva/qed/gossip"
	"github.com/bbva/qed/gossip/auditor"
	"github.com/bbva/qed/gossip/member"
	"github.com/bbva/qed/gossip/monitor"
	"github.com/bbva/qed/gossip/publisher"
	"github.com/bbva/qed/server"
	"github.com/bbva/qed/testutils/scope"
)

const (
	QEDUrl       = "http://127.0.0.1:8800"
	QEDTLS       = "https://localhost:8800"
	QEDGossip    = "127.0.0.1:8400"
	QEDTamperURL = "http://127.0.0.1:18800/tamper"
	StoreURL     = "http://127.0.0.1:8888"
	APIKey       = "my-key"
)

// merge function is a helper function that execute all the variadic parameters
// inside a score.TestF function
func merge(list ...scope.TestF) scope.TestF {
	return func(t *testing.T) {
		for _, elem := range list {
			elem(t)
		}
	}
}

func delay(duration time.Duration) scope.TestF {
	return func(t *testing.T) {
		time.Sleep(duration)
	}
}

func doReq(method string, url, apiKey string, payload *strings.Reader) (*http.Response, error) {
	var err error
	if payload == nil {
		payload = strings.NewReader("")
	}
	req, err := http.NewRequest(method, url, payload)
	if err != nil {
		return nil, err
	}

	req.Header.Set("Api-Key", apiKey)
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	return resp, err
}

func newAgent(id int, name string, role member.Type, p gossip.Processor, t *testing.T) *gossip.Agent {
	agentConf := gossip.DefaultConfig()
	agentConf.NodeName = fmt.Sprintf("%s%d", name, id)

	switch role {
	case member.Auditor:
		agentConf.BindAddr = fmt.Sprintf("127.0.0.1:810%d", id)
	case member.Monitor:
		agentConf.BindAddr = fmt.Sprintf("127.0.0.1:820%d", id)
	case member.Publisher:
		agentConf.BindAddr = fmt.Sprintf("127.0.0.1:830%d", id)
	}

	agentConf.StartJoin = []string{QEDGossip}
	agentConf.EnableCompression = true
	agentConf.AlertsUrls = []string{StoreURL}
	agentConf.Role = role

	agent, err := gossip.NewAgent(agentConf, []gossip.Processor{p})
	if err != nil {
		t.Fatalf("Failed to start AGENT %s: %v", name, err)
	}
	_, _ = agent.Join([]string{QEDGossip})
	return agent
}

func setupAuditor(id int, t *testing.T) (scope.TestF, scope.TestF) {
	var au *auditor.Auditor
	var agent *gossip.Agent
	var err error

	before := func(t *testing.T) {
		auditorConf := auditor.DefaultConfig()
		auditorConf.MetricsAddr = fmt.Sprintf("127.0.0.1:710%d", id)
		auditorConf.QEDUrls = []string{QEDUrl}
		auditorConf.PubUrls = []string{StoreURL}
		auditorConf.APIKey = APIKey

		au, err = auditor.NewAuditor(*auditorConf)
		if err != nil {
			t.Fatalf("Unable to create a new auditor: %v", err)
		}

		agent = newAgent(id, "auditor", member.Auditor, au, t)
	}

	after := func(t *testing.T) {
		if au != nil {
			au.Shutdown()
		}
		err := agent.Leave()
		if err != nil {
			t.Fatalf("Unable to shutdown the auditor: %v", err)
		}
		err = agent.Shutdown()
		if err != nil {
			t.Fatalf("Unable to shutdown the auditor: %v", err)
		}
	}
	return before, after
}

func setupMonitor(id int, t *testing.T) (scope.TestF, scope.TestF) {
	var mn *monitor.Monitor
	var agent *gossip.Agent
	var err error

	before := func(t *testing.T) {
		monitorConf := monitor.DefaultConfig()
		monitorConf.MetricsAddr = fmt.Sprintf("127.0.0.1:720%d", id)
		monitorConf.QEDUrls = []string{QEDUrl}
		monitorConf.PubUrls = []string{StoreURL}
		monitorConf.APIKey = APIKey

		mn, err = monitor.NewMonitor(*monitorConf)
		if err != nil {
			t.Fatalf("Unable to create a new monitor: %v", err)
		}

		agent = newAgent(id, "monitor", member.Monitor, mn, t)
	}

	after := func(t *testing.T) {
		if mn != nil {
			mn.Shutdown()
		}
		err := agent.Leave()
		if err != nil {
			t.Fatalf("Unable to shutdown the monitor: %v", err)
		}
		err = agent.Shutdown()
		if err != nil {
			t.Fatalf("Unable to shutdown the monitor: %v", err)
		}
	}
	return before, after
}

func setupPublisher(id int, t *testing.T) (scope.TestF, scope.TestF) {
	var pu *publisher.Publisher
	var agent *gossip.Agent
	var err error

	before := func(t *testing.T) {
		conf := publisher.DefaultConfig()
		conf.MetricsAddr = fmt.Sprintf("127.0.0.1:730%d", id)
		conf.PubUrls = []string{StoreURL}

		pu, err = publisher.NewPublisher(*conf)
		if err != nil {
			t.Fatalf("Unable to create a new publisher: %v", err)
		}

		agent = newAgent(id, "publisher", member.Publisher, pu, t)
	}

	after := func(t *testing.T) {
		if pu != nil {
			pu.Shutdown()
		}
		err := agent.Leave()
		if err != nil {
			t.Fatalf("Unable to shutdown the publisher: %v", err)
		}
		err = agent.Shutdown()
		if err != nil {
			t.Fatalf("Unable to shutdown the publisher: %v", err)
		}
	}
	return before, after
}

func setupStore(t *testing.T) (scope.TestF, scope.TestF) {
	var s *Service
	before := func(t *testing.T) {
		s = NewService()
		foreground := false
		s.Start(foreground)
	}

	after := func(t *testing.T) {
		s.Shutdown()
	}
	return before, after
}

func setupServer(id int, joinAddr string, tls bool, t *testing.T) (scope.TestF, scope.TestF) {
	var srv *server.Server
	var err error
	path := fmt.Sprintf("/var/tmp/e2e-qed%d/", id)

	usr, _ := user.Current()

	before := func(t *testing.T) {
		os.RemoveAll(path)
		err = os.MkdirAll(path, os.FileMode(0755))
		if err != nil {
			t.Fatalf("Unable to create a path: %v", err)
		}

		hostname, _ := os.Hostname()
		conf := server.DefaultConfig()
		conf.APIKey = APIKey
		conf.NodeID = fmt.Sprintf("%s-%d", hostname, id)
		conf.HTTPAddr = fmt.Sprintf("127.0.0.1:880%d", id)
		conf.MgmtAddr = fmt.Sprintf("127.0.0.1:870%d", id)
		conf.MetricsAddr = fmt.Sprintf("127.0.0.1:860%d", id)
		conf.RaftAddr = fmt.Sprintf("127.0.0.1:850%d", id)
		conf.GossipAddr = fmt.Sprintf("127.0.0.1:840%d", id)
		if id > 0 {
			conf.RaftJoinAddr = []string{"127.0.0.1:8700"}
			conf.GossipJoinAddr = []string{"127.0.0.1:8400"}
		}
		conf.DBPath = path + "data"
		conf.RaftPath = path + "raft"
		conf.PrivateKeyPath = fmt.Sprintf("%s/.ssh/id_ed25519", usr.HomeDir)
		if tls {
			conf.SSLCertificate = fmt.Sprintf("%s/.ssh/server.crt", usr.HomeDir)
			conf.SSLCertificateKey = fmt.Sprintf("%s/.ssh/server.key", usr.HomeDir)
		}
		conf.EnableProfiling = true
		conf.EnableTampering = true
		conf.EnableTLS = tls

		//fmt.Printf("Server config: %+v\n", conf)

		srv, err = server.NewServer(conf)
		if err != nil {
			t.Fatalf("Unable to create a new server: %v", err)
		}

		go (func() {
			err := srv.Start()
			if err != nil {
				t.Log(err)
			}
		})()
		time.Sleep(2 * time.Second)
	}

	after := func(t *testing.T) {
		if srv != nil {
			err := srv.Stop()
			if err != nil {
				t.Fatalf("Unable to shutdown the server! %v", err)
			}
		} else {
			t.Fatalf("Unable to shutdown the server!")
		}
	}
	return before, after
}

func endPoint(id int) string {
	return fmt.Sprintf("http://127.0.0.1:880%d", id)
}

func getClient(id int) *client.HTTPClient {
	return client.NewHTTPClient(client.Config{
		Endpoints: []string{endPoint(id)},
		APIKey:    APIKey,
		Insecure:  false,
	})
}
