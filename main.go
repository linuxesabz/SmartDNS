package main

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"os"

	"github.com/miekg/dns"
)

var records = map[string]string{
	"test.service.": "192.168.0.2",
	"test.com.":     "192.168.1.1",
}

func parseQuery(m *dns.Msg) {
	for _, q := range m.Question {
		switch q.Qtype {
		case dns.TypeA:
			log.Printf("Query for %s\n", q.Name)
			ip := records[q.Name]
			if ip != "" {
				rr, err := dns.NewRR(fmt.Sprintf("%s A %s", q.Name, ip))
				if err == nil {
					m.Answer = append(m.Answer, rr)
				}
			}
		}
	}
}

func ReadConfigfile(ConfigName string) map[string]interface{} {
	var configJson map[string]interface{}
	configfile, err := os.Open(ConfigName)

	if err != nil {
		log.Println(err)
	}
	jsonResult, err := io.ReadAll(configfile)
	defer configfile.Close()
	err = json.Unmarshal([]byte(jsonResult), &configJson)
	if err != nil {
		log.Println(err)
	}
	return configJson
}
func handleDnsRequest(w dns.ResponseWriter, r *dns.Msg) {
	m := new(dns.Msg)
	m.SetReply(r)
	m.Compress = false

	switch r.Opcode {
	case dns.OpcodeQuery:
		parseQuery(m)
	}

	w.WriteMsg(m)
}

func main() {
	var ConfigJson map[string]interface{}
	// attach request handler func
	dns.HandleFunc(".", handleDnsRequest)
	ConfigJson = ReadConfigfile("config.conf")

	// start server
	port := ConfigJson["Port"].(string)
	protocol := ConfigJson["Protocol"].(string)
	server := &dns.Server{Addr: ":" + port, Net: protocol}
	log.Printf("Starting at %s\n", port)
	err := server.ListenAndServe()
	defer server.Shutdown()
	if err != nil {
		log.Fatalf("Failed to start server: %s\n ", err.Error())
	}
}
