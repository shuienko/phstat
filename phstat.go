package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"sort"
	"strconv"

	ui "github.com/gizak/termui"
	gohole "github.com/shuienko/go-pihole"
)

const (
	readme = "https://github.com/shuienko/phstat/blob/master/README.md"
	nItems = 10
)

// Show help page
func usage() {
	fmt.Println("Pi-Hole Dashboard. Docs:", readme)
	fmt.Println("Usage:", os.Args[0], "[-h] [-n seconds]")
	flag.PrintDefaults()
}

// Sort and reverse map
func sortReverseMap(m map[string]int) (map[int]string, []int) {
	reverseMap := make(map[int]string)
	var freq []int

	for k, v := range m {
		reverseMap[v] = k
		freq = append(freq, v)
	}

	sort.Ints(freq)

	return reverseMap, freq
}

// Show summary
func getSummary(ph gohole.PiHConnector) []string {
	s := ph.Summary()
	sData := []string{
		"Status: [" + s.Status + "](fg-blue)",
		"Blocked Domains: [" + strconv.Itoa(s.AdsBlocked) + "](fg-blue)",
		"Blocked Percentage: [" + fmt.Sprintf("%f", s.AdsPercentage) + "%](fg-blue)",
		"DNS Queries Today: [" + strconv.Itoa(s.DNSQueries) + "](fg-blue)",
		"Domains Being Blocked: [" + strconv.Itoa(s.DomainsBlocked) + "](fg-blue)",
		"Queries Cached: [" + strconv.Itoa(s.QueriesCached) + "](fg-blue)",
		"Queries Forwarded: [" + strconv.Itoa(s.QueriesForwarded) + "](fg-blue)",
		"Clients Ever Seen: [" + strconv.Itoa(s.ClientsEverSeen) + "](fg-blue)",
		"Unique Clients: [" + strconv.Itoa(s.UniqueClients) + "](fg-blue)",
		"Unique Domains: [" + strconv.Itoa(s.UniqueDomains) + "](fg-blue)",
	}
	return sData
}

// Show top blocked DNS records
func getTopBlocked(ph gohole.PiHConnector) []string {
	b := ph.Top(nItems).Blocked
	var bData []string

	reverseMapBlocked, freqBlocked := sortReverseMap(b)

	for i := len(freqBlocked) - 1; i >= 0; i-- {
		row := fmt.Sprintf("[%5d](fg-blue): %s", freqBlocked[i], reverseMapBlocked[freqBlocked[i]])
		bData = append(bData, row)
	}
	return bData
}

// Show top DNS Queries
func getTopQueries(ph gohole.PiHConnector) []string {
	q := ph.Top(nItems).Queries
	var qData []string

	reverseMapQueries, freqQueries := sortReverseMap(q)

	for i := len(freqQueries) - 1; i >= 0; i-- {
		row := fmt.Sprintf("[%5d](fg-blue): %s", freqQueries[i], reverseMapQueries[freqQueries[i]])
		qData = append(qData, row)
	}

	return qData
}

// Show top clients
func getTopClients(ph gohole.PiHConnector) []string {
	c := ph.Clients(nItems).Clients
	var cData []string

	reverseMapClients, freqClients := sortReverseMap(c)

	for i := len(freqClients) - 1; i >= 0; i-- {
		row := fmt.Sprintf("[%5d](fg-blue): %s", freqClients[i], reverseMapClients[freqClients[i]])
		cData = append(cData, row)
	}

	return cData
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

	// Parse options
	var updateInterval = flag.Uint64("n", 2, "update interval in `seconds`")
	var help = flag.Bool("h", false, "show this page")
	flag.Parse()

	if *help {
		usage()
		os.Exit(0)
	}

	// Start UI
	err := ui.Init()
	if err != nil {
		panic(err)
	}
	defer ui.Close()

	termWidth := ui.TermWidth()
	par0 := ui.NewPar("Pi-Hole Dashboard " + "[http://" + ph.Host + "/admin/index.php](fg-blue)")
	par0.Height = 1
	par0.Width = 82
	par0.Border = false
	par0.TextFgColor = ui.ColorGreen
	par0.PaddingLeft = 1
	par0.Y = 1
	par0.Float = ui.AlignCenterHorizontal

	apiString := fmt.Sprintf("API: [%.1f %s](fg-blue)", ph.Version().Version, ph.Type().Type)
	par1 := ui.NewPar(apiString)
	par1.Height = 1
	par1.Width = 82
	par1.Border = false
	par1.TextFgColor = ui.ColorYellow
	par1.PaddingLeft = 1
	par1.Y = 2
	par1.Float = ui.AlignCenterHorizontal

	par2 := ui.NewPar("Last Blocked: [" + ph.RecentBlocked() + "](fg-blue)")
	par2.Height = 3
	par2.Width = 82
	par2.Border = true
	par2.TextFgColor = ui.ColorYellow
	par2.PaddingLeft = 1
	par2.Y = 3
	par2.Float = ui.AlignCenterHorizontal

	// Summary
	ls1 := ui.NewList()
	ls1.Items = getSummary(ph)
	ls1.ItemFgColor = ui.ColorYellow
	ls1.BorderLabel = "Summary"
	ls1.Height = 12
	ls1.Width = 40
	ls1.PaddingLeft = 1
	ls1.Y = 6
	ls1.X = ui.TermWidth()/2 - 41

	// Top Blocked
	ls2 := ui.NewList()
	ls2.Items = getTopBlocked(ph)
	ls2.ItemFgColor = ui.ColorYellow
	ls2.BorderLabel = "Top Blocked"
	ls2.Height = 12
	ls2.Width = 40
	ls2.PaddingLeft = 1
	ls2.Y = 6
	ls2.X = ui.TermWidth()/2 + 1

	// Top Queries
	ls3 := ui.NewList()
	ls3.Items = getTopQueries(ph)
	ls3.ItemFgColor = ui.ColorYellow
	ls3.BorderLabel = "Top Queries"
	ls3.Height = 12
	ls3.Width = 40
	ls3.PaddingLeft = 1
	ls3.Y = 18
	ls3.X = ui.TermWidth()/2 + 1

	// Top Clients
	ls4 := ui.NewList()
	ls4.Items = getTopClients(ph)
	ls4.ItemFgColor = ui.ColorYellow
	ls4.BorderLabel = "Top Clients"
	ls4.Height = 12
	ls4.Width = 40
	ls4.PaddingLeft = 1
	ls4.Y = 18
	ls4.X = ui.TermWidth()/2 - 41

	// Render
	ui.Render(ls1, ls2, ls3, ls4, par0, par1, par2)
	ui.Handle("/sys/kbd/q", func(ui.Event) {
		ui.StopLoop()
	})

	// Render periodically
	ui.Handle("/timer/1s", func(e ui.Event) {
		t := e.Data.(ui.EvtTimer)
		if t.Count%*updateInterval == 0 {
			par2.Text = "Last Blocked: [" + ph.RecentBlocked() + "](fg-blue)"
			ls1.Items = getSummary(ph)
			ls2.Items = getTopBlocked(ph)
			ls3.Items = getTopQueries(ph)
			ls4.Items = getTopClients(ph)
			if ui.TermWidth() != termWidth {
				ls1.X = ui.TermWidth()/2 - 41
				ls2.X = ui.TermWidth()/2 + 1
				ls3.X = ui.TermWidth()/2 + 1
				ls4.X = ui.TermWidth()/2 - 41

				termWidth = ui.TermWidth()
				ui.Clear()
			}

			ui.Render(par0, par1, par2, ls1, ls2, ls3, ls4)
		}
	})
	ui.Loop()
}
