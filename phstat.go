package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"sort"

	ui "github.com/gizak/termui"
	"github.com/shuienko/go-pihole"
)

func usage() {
	fmt.Println("Usage:", os.Args[0], "[-n NUMBER] summary|blocked|queries|clients|type|version|enable|disable|recent")
	flag.PrintDefaults()
}

func main() {
	// Get environment variables
	piholeHost, ok := os.LookupEnv("PIHOLE_HOST")
	if !ok {
		log.Fatal("PIHOLE_HOST environment variable in NOT set")
	}

	apiToken, ok := os.LookupEnv("PIHOLE_TOKEN")
	if !ok {
		log.Fatal("PIHOLE_TOKEN environment variable is NOT set")
	}

	// Create connector object
	ph := gohole.PiHConnector{
		Host:  piholeHost,
		Token: apiToken,
	}

	err := ui.Init()
	if err != nil {
		panic(err)
	}
	defer ui.Close()

	// Summary
	s := ph.Summary()
	sData := []string{
		"Status: [" + s.Status + "](fg-blue)",
		"Blocked Domains: [" + s.AdsBlocked + "](fg-blue)",
		"Blocked Percentage: [" + s.AdsPercentage + "](fg-blue)",
		"DNS Queries Today: [" + s.DNSQueries + "](fg-blue)",
		"Domains Being Blocked: [" + s.DomainsBlocked + "](fg-blue)",
		"Queries Cached: [" + s.QueriesCached + "](fg-blue)",
		"Queries Forwarded: [" + s.QueriesForwarded + "](fg-blue)",
		"Clients Ever Seen: [" + s.ClientsEverSeen + "](fg-blue)",
		"Unique Clients: [" + s.UniqueClients + "](fg-blue)",
		"Unique Domains: [" + s.UniqueDomains + "](fg-blue)",
	}

	ls1 := ui.NewList()
	ls1.Items = sData
	ls1.ItemFgColor = ui.ColorYellow
	ls1.BorderLabel = "Summary"
	ls1.Height = 10
	ls1.Width = 35
	ls1.PaddingLeft = 1
	ls1.Y = 0

	// Top Blocked
	b := ph.Top(10).Blocked
	var bData []string

	reverseMapBlocked := make(map[int]string)
	var freqBlocked []int

	for k, v := range b {
		reverseMapBlocked[v] = k
		freqBlocked = append(freqBlocked, v)
	}

	sort.Ints(freqBlocked)

	for i := len(freqBlocked) - 1; i >= 0; i-- {
		row := fmt.Sprintf("[%5d](fg-blue): %s", freqBlocked[i], reverseMapBlocked[freqBlocked[i]])
		bData = append(bData, row)
	}

	ls2 := ui.NewList()
	ls2.Items = bData
	ls2.ItemFgColor = ui.ColorYellow
	ls2.BorderLabel = "Top Blocked"
	ls2.Height = 10
	ls2.Width = 35
	ls2.PaddingLeft = 1
	ls2.Y = 0
	ls2.X = 37

	// Top Queries
	q := ph.Top(10).Queries
	var qData []string

	reverseMapQueries := make(map[int]string)
	var freqQueries []int

	for k, v := range q {
		reverseMapQueries[v] = k
		freqQueries = append(freqQueries, v)
	}

	sort.Ints(freqQueries)

	for i := len(freqQueries) - 1; i >= 0; i-- {
		row := fmt.Sprintf("[%5d](fg-blue): %s", freqQueries[i], reverseMapQueries[freqQueries[i]])
		qData = append(qData, row)
	}

	ls3 := ui.NewList()
	ls3.Items = qData
	ls3.ItemFgColor = ui.ColorYellow
	ls3.BorderLabel = "Top Queries"
	ls3.Height = 10
	ls3.Width = 35
	ls3.PaddingLeft = 1
	ls3.Y = 10
	ls3.X = 37

	ui.Render(ls1, ls2, ls3)
	ui.Handle("/sys/kbd/q", func(ui.Event) {
		ui.StopLoop()
	})
	ui.Loop()

}
