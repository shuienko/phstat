package main

import (
	"flag"
	"fmt"
	"log"
	"os"

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
		fmt.Println("API type:", phtype.Type)
	case "version":
		phversion := ph.Version()
		fmt.Println("API version:", phversion.Version)
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
