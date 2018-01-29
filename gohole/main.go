package gohole

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"sort"
)

// ParseJSON conterts raw binary output into object
func ParseJSON(r interface{}, body []byte) error {
	err := json.Unmarshal(body, r)
	if err != nil {
		log.Fatal(err)
	}
	return err
}

// PiConnector type
type PiConnector struct {
	Host  string
	Token string
}

// Get request to API
func (r *PiConnector) Get(endpoint string) []byte {
	var requestString = "http://" + r.Host + "/admin/api.php?" + endpoint
	if r.Token != "" {
		requestString += "&auth=" + r.Token
	}

	resp, err := http.Get(requestString)
	if err != nil {
		log.Fatal(err)
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Fatal(err)
	}

	return body
}

// PiResponseSummary type
type PiResponseSummary struct {
	BlockedDomains   int     `json:"domains_being_blocked"`
	QueriesToday     int     `json:"dns_queries_today"`
	BlockedToday     int     `json:"ads_blocked_today"`
	BlockedPercent   float64 `json:"ads_percentage_today"`
	UniqueDomains    int     `json:"unique_domains"`
	QueriesForwarded int     `json:"queries_forwarded"`
	QueriesCached    int     `json:"queries_cached"`
	ClientsEverSeen  int     `json:"clients_ever_seen"`
	UniqueClients    int     `json:"unique_clients"`
	Status           string  `json:"status"`
}

// PiResponseTop type
type PiResponseTop struct {
	Queries map[string]int `json:"top_queries"`
	Blocked map[string]int `json:"top_ads"`
}

// PiResponseTopClients type
type PiResponseTopClients struct {
	Clients map[string]int `json:"top_sources"`
}

// Show returns 24h Summary of PiHole System
func (r *PiResponseSummary) Show() {
	fmt.Println("=== 24h Summary:")
	fmt.Printf("- Blocked Domains: %d\n", r.BlockedToday)
	fmt.Printf("- Blocked Percentage: %.2f%%\n", r.BlockedPercent)
	fmt.Printf("- Queries: %d\n", r.QueriesToday)
	fmt.Printf("- Clients Ever Seen: %d\n", r.ClientsEverSeen)
}

// ShowBlocked returns sorted top Blocked domains over last 24h
func (r *PiResponseTop) ShowBlocked() {
	reverseMapBlocked := make(map[int]string)
	var freqBlocked []int

	for k, v := range r.Blocked {
		reverseMapBlocked[v] = k
		freqBlocked = append(freqBlocked, v)
	}

	sort.Ints(freqBlocked)

	fmt.Println("=== Blocked domains over last 24h:")
	for i := len(freqBlocked) - 1; i >= 0; i-- {
		fmt.Printf("- %s : %d\n", reverseMapBlocked[freqBlocked[i]], freqBlocked[i])
	}
}

// ShowQueries returns sorted top queries over last 24h
func (r *PiResponseTop) ShowQueries() {
	reverseMapQueries := make(map[int]string)
	var freqQueries []int

	for k, v := range r.Queries {
		reverseMapQueries[v] = k
		freqQueries = append(freqQueries, v)
	}

	sort.Ints(freqQueries)

	fmt.Println("=== Queries over last 24h:")
	for i := len(freqQueries) - 1; i >= 0; i-- {
		fmt.Printf("- %s : %d\n", reverseMapQueries[freqQueries[i]], freqQueries[i])
	}
}

// ShowClients returns sorted top clients over last 24h
func (r *PiResponseTopClients) ShowClients() {
	reverseMapClients := make(map[int]string)
	var freqClients []int

	for k, v := range r.Clients {
		reverseMapClients[v] = k
		freqClients = append(freqClients, v)
	}

	sort.Ints(freqClients)

	fmt.Println("=== Clients over last 24h:")
	for i := len(freqClients) - 1; i >= 0; i-- {
		fmt.Printf("- %s : %d\n", reverseMapClients[freqClients[i]], freqClients[i])
	}
}
