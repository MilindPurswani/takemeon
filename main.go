package main

import (
	"bufio"
	"flag"
	"fmt"
	"net"
	"os"
	"strings"
	"sync"

	"github.com/miekg/dns"
)

type domainList []string

var dnserver string

func main() {
	// cli for setting concurrency
	var concurrency int
	ds := dnserver
	flag.IntVar(&concurrency, "c", 1, "set the concurrency level")

	// cli for specifying dns server manually
	var mdns string
	flag.StringVar(&mdns, "mdns", "/etc/resolv.conf", "Manually Specify dns server IP address only. (Just a little faster)")
	flag.Parse()
	// Use the right dns server
	if strings.Compare(mdns, "/etc/resolv.conf") == 0 {
		config, _ := dns.ClientConfigFromFile("/etc/resolv.conf")
		ds = config.Servers[0] + ":" + config.Port
	} else {
		ds = mdns + ":53"
	}

	c := new(dns.Client)
	m := new(dns.Msg)
	sc := bufio.NewScanner(os.Stdin)
	domains := domainList{}
	for sc.Scan() {
		domains = append(domains, sc.Text())
	}
	if err := sc.Err(); err != nil {
		fmt.Fprintf(os.Stderr, "failed to read input from file %s\n", err)
	}

	wg := new(sync.WaitGroup)
	for i := 0; i < concurrency; i++ {
		wg.Add(1)
		go func() {

			for _, d := range domains {
				_, err := net.LookupCNAME(d)
				if err != nil {
					m.SetQuestion(dns.Fqdn(d), dns.TypeCNAME)
					m.RecursionDesired = true
					r, _, err2 := c.Exchange(m, ds)
					// Check if the domain is actually not existing
					if err2 != nil {
						continue
					}
					// Check to see if the Answer's length is 0
					if len(r.Answer) == 0 {
						continue
					}
					if r, ok := r.Answer[0].(*dns.CNAME); ok {
						fmt.Print(d + " | ")
						fmt.Println(r.Target[:len(r.Target)-1])
					}

				}
			}

			wg.Done()
		}()
	}
	wg.Wait()
}
