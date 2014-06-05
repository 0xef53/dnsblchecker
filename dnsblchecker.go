package main

import (
	"bufio"
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"net"
	"os"
	"strconv"
	"strings"
)

var (
	JSONOUT   = false
	MAXPROCS  = 5
	DNSBLLIST = "dnsbl.txt"
)

type DnsblStatus struct {
	Dnsbl    string `json:"dnsbl"`
	IsListed bool   `json:"islisted"`
}

func getReverseIP(ipStr string) []string {
	ip := strings.Split(ipStr, ".")
	for i, j := 0, len(ip)-1; i < j; i, j = i+1, j-1 {
		ip[i], ip[j] = ip[j], ip[i]
	}
	return ip
}

func checkRecord(ipStr, dnsblAddr string) *DnsblStatus {
	var isListed bool
	d := strings.Join(append(getReverseIP(ipStr), dnsblAddr), ".")
	if _, ok := net.LookupHost(d); ok == nil {
		isListed = true
	} else {
		isListed = false
	}
	return &DnsblStatus{dnsblAddr, isListed}
}

func worker(jobs <-chan [2]string, results chan<- *DnsblStatus) {
	for j := range jobs {
		results <- checkRecord(j[0], j[1])
	}
}

func startJob(ipStr string, results chan<- *DnsblStatus) {
	f, err := os.Open(DNSBLLIST)
	if err != nil {
		log.Fatalln(err)
	}
	defer f.Close()
	jobs := make(chan [2]string)
	for i := 0; i < MAXPROCS; i++ {
		go worker(jobs, results)
	}
	scanner := bufio.NewScanner(f)
	for scanner.Scan() {
		if strings.Trim(scanner.Text(), " ") != "" {
			jobs <- [2]string{ipStr, scanner.Text()}
		}
	}
	if err := scanner.Err(); err != nil {
		log.Fatalln(err)
	}
}

func countRows() int {
	f, err := os.Open(DNSBLLIST)
	if err != nil {
		log.Fatalln(err)
	}
	defer f.Close()
	c := 0
	scanner := bufio.NewScanner(f)
	for scanner.Scan() {
		if strings.Trim(scanner.Text(), " ") != "" {
			c += 1
		}
	}
	if err := scanner.Err(); err != nil {
		log.Fatalln(err)
	}
	return c
}

func init() {
	flag.BoolVar(&JSONOUT, "json", JSONOUT, "show output in the JSON format")
	flag.IntVar(&MAXPROCS, "P", MAXPROCS, "number of concurrent processes")
	flag.StringVar(&DNSBLLIST, "c", DNSBLLIST, "path to the dnsbl servers list")
}

func main() {
	flag.Parse()
	if flag.NArg() != 1 {
		os.Exit(1)
	}
	ipStr := flag.Arg(0)
	results := make(chan *DnsblStatus)
	numRows := countRows()

	go startJob(ipStr, results)

	stList := make([]*DnsblStatus, 0, 100)
	for i := 0; i < numRows; i++ {
		stList = append(stList, <-results)
	}

	if JSONOUT {
		jOut, err := json.Marshal(stList)
		if err != nil {
			log.Fatalln(err)
		}
		fmt.Printf(string(jOut))
	} else {
		for k, v := range stList {
			fmt.Printf("%d:%s:%s\n", k, v.Dnsbl, strconv.FormatBool(v.IsListed))
		}
	}
}
