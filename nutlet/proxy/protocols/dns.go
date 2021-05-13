package protocols

import (
	"encoding/json"
	"fmt"
	log "github.com/sirupsen/logrus"
	"regexp"
	"strconv"
	"strings"

	"github.com/miekg/dns"
)

var records = map[string]string {
	"*.local.nutshell": "127.0.0.1",
	"dev.local.nutshell": "127.0.0.1",
}

func match(req string) (string, bool) {
	for pattern, target := range records {
		pattern = strings.ReplaceAll(pattern, "*", ".*")
		if matched, _:= regexp.MatchString(pattern, req); matched {
			return target, true
		}
	}
	return "", false
}

func parseQuery(m *dns.Msg) {
	for _, q := range m.Question {
		name := strings.ReplaceAll(q.Name,  "nutshell.", "nutshell")
		switch q.Qtype {
		case dns.TypeA:
			ip, matched := match(name)
			if matched {
				rr, err := dns.NewRR(fmt.Sprintf("%s A %s", q.Name, ip))
				if err == nil {
					m.Answer = append(m.Answer, rr)
				}
			}
		case dns.TypeAAAA:
			ip, matched := match(name)
			if matched {
				rr, err := dns.NewRR(fmt.Sprintf("%s AAAA %s", q.Name, ip))
				if err == nil {
					m.Answer = append(m.Answer, rr)
				}
			}
		}
		log.Printf("Query for %s\n target:%s", q.Name, m.Answer)
	}
}

func handleDnsRequest(w dns.ResponseWriter, r *dns.Msg) {
	m := new(dns.Msg)
	m.SetReply(r)
	m.Compress = false

	d , _ := json.Marshal(m)
	log.Infof("handle dns request:%s", d)
	switch r.Opcode {
	case dns.OpcodeQuery:
		parseQuery(m)
	}

	w.WriteMsg(m)
}

func ServeDns(port int, rules map[string]string) func(host, ip string, del bool ) {
	// attach request handler func
	dns.HandleFunc("nutshell", handleDnsRequest)
	// start server
	server := &dns.Server{Addr: ":" + strconv.Itoa(port), Net: "udp"}
	log.Printf("Starting at %d\n", port)
	go func() {
		defer server.Shutdown()
		err := server.ListenAndServe()
		if err != nil {
			log.Fatalf("Failed to start server: %s\n ", err.Error())
		}
	}()
	if rules != nil {
		records = rules
	}
	return func(host, ip string, del bool) {
		if del {
			delete(records, host)
		} else {
			records[host] = ip
		}
	}
}
