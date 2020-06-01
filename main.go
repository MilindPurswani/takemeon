package main

import (
	"bufio"
	"fmt"
	"net"
	"os"

	"github.com/miekg/dns"
)

type domainList []string

func main() {

	// TODO: Implement Concurrency
	// var concurrency int
	// flag.IntVar(&concurrency, "c", 20, "set the concurrency level")
	// flag.Parse()
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
	for _, d := range domains {
		_, err := net.LookupCNAME(d)
		if err != nil {
			//fmt.Fprintf(os.Stderr, "Illegal domain %s\n", err)
			m.SetQuestion(d+".", dns.TypeCNAME)
			m.RecursionDesired = true
			r, _, err2 := c.Exchange(m, config.Servers[0]+":"+config.Port)
			if err2 != nil {
				fmt.Print(d + " | ")
				fmt.Println(r.Answer[0].(*dns.CNAME).Target[:len(r.Answer[0].(*dns.CNAME).Target)-1])
			}

		}

	}

}
