package main

import (
	"flag"
	"fmt"
	"log"
	"os"

	"github.com/shuienko/go-pihole"
)

func usage() {
	fmt.Println("Usage:", os.Args[0], "[-n NUMBER] summary|blocked|queries|clients")
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

	// Get command line arguments
	var maxNum = flag.Int("n", 10, "`number` of returned entries")
	flag.Parse()

	var arg string
	if len(flag.Args()) > 0 {
		arg = flag.Args()[0]
	} else {
		usage()
		os.Exit(1)
	}

	// Show output based on arguments and options
	switch arg {
	case "summary":
		summary := ph.SummaryRaw()
		summary.Show()
	case "blocked":
		topItems := ph.Top(*maxNum)
		topItems.ShowBlocked()
	case "queries":
		topItems := ph.Top(*maxNum)
		topItems.ShowQueries()
	case "clients":
		clients := ph.Clients(*maxNum)
		clients.Show()
	case "type":
		phtype := ph.Type()
		fmt.Println(phtype.Type)
	case "version":
		phversion := ph.Version()
		fmt.Println(phversion.Version)
	case "summaryRaw":
		summary := ph.Summary()
		fmt.Println(summary)
	case "timedata":
		data := ph.TimeData()
		fmt.Println(data)
	case "fd":
		fd := ph.ForwardDestinations()
		fmt.Println(fd)
	case "qt":
		qt := ph.QueryTypes()
		fmt.Println(qt)
	case "allqueries":
		aq := ph.Queries()
		fmt.Println(aq.Data[0], len(aq.Data))
	case "enable":
		err := ph.Enable()
		if err != nil {
			fmt.Println(err)
		} else {
			fmt.Println("Enabled")
		}
	case "disable":
		err := ph.Disable()
		if err != nil {
			fmt.Println(err)
		} else {
			fmt.Println("Disabled")
		}
	case "recent":
		fmt.Println(ph.RecentBlocked())
	default:
		usage()
	}
}
