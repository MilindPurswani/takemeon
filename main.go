package main

import (
	"bufio"
	"flag"
	"fmt"
	"net"
	"os"
	"sync"

	"github.com/miekg/dns"
)

type domainList []string

func main() {

	var concurrency int
	flag.IntVar(&concurrency, "c", 20, "set the concurrency level")
	flag.Parse()
	config, _ := dns.ClientConfigFromFile("/etc/resolv.conf")
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
					r, _, err2 := c.Exchange(m, config.Servers[0]+":"+config.Port)
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
