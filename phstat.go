package main

import (
	"flag"
	"fmt"
	"log"
	"os"

	"github.com/shuienko/phstat/gohole"
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
	c := gohole.PiConnector{
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
		summary := c.Summary()
		summary.Show()
	case "blocked":
		topItems := c.Top(*maxNum)
		topItems.ShowBlocked()
	case "queries":
		topItems := c.Top(*maxNum)
		topItems.ShowQueries()
	case "clients":
		clients := c.Clients(*maxNum)
		clients.Show()
	default:
		usage()
	}
}
