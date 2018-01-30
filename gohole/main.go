// Package gohole provides a client for the Pi-Hole API.
// In order to use this package you will need Pi-Hole's HTTP port 80 to be available.
// Important: only AdminLTE v3.0+
package gohole

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"sort"
	"strconv"
)

// PiHConnector represents base API connector type
// Host: DNS or IP address of your Pi-Hole Token
// Token: API Token (see /etc/pihole/setupVars.conf)
type PiHConnector struct {
	Host  string
	Token string
}

// PiHType coitains backend Type (PHP or FTL)
type PiHType struct {
	Type string `json:"type"`
}

// PiHVersion contains API version
type PiHVersion struct {
	Version float32 `json:"version"`
}

// PiHSummaryRaw contains raw Pi-Hole summary data
type PiHSummaryRaw struct {
	AdsBlocked       int     `json:"ads_blocked_today"`
	AdsPercentage    float64 `json:"ads_percentage_today"`
	ClientsEverSeen  int     `json:"clients_ever_seen"`
	DNSQueries       int     `json:"dns_queries_today"`
	DomainsBlocked   int     `json:"domains_being_blocked"`
	QueriesCached    int     `json:"queries_cached"`
	QueriesForwarded int     `json:"queries_forwarded"`
	Status           string  `json:"status"`
	UniqueClients    int     `json:"unique_clients"`
	UniqueDomains    int     `json:"unique_domains"`
}

// PiHSummary contains Pi-Hole statistics in formatted style. All fields are strings
type PiHSummary struct {
	AdsBlocked       string `json:"ads_blocked_today"`
	AdsPercentage    string `json:"ads_percentage_today"`
	ClientsEverSeen  string `json:"clients_ever_seen"`
	DNSQueries       string `json:"dns_queries_today"`
	DomainsBlocked   string `json:"domains_being_blocked"`
	QueriesCached    string `json:"queries_cached"`
	QueriesForwarded string `json:"queries_forwarded"`
	Status           string `json:"status"`
	UniqueClients    string `json:"unique_clients"`
	UniqueDomains    string `json:"unique_domains"`
}

// PiHTimeData represents statistics over time.
// Each record contains number of queries/blocked ads within 10min timeframe
type PiHTimeData struct {
	AdsOverTime     map[string]int `json:"ads_over_time"`
	DomainsOverTime map[string]int `json:"domains_over_time"`
}

// PiHTopItems contains top queries/blocked ads
// Format: "DNS": Frequency
type PiHTopItems struct {
	Queries map[string]int `json:"top_queries"`
	Blocked map[string]int `json:"top_ads"`
}

// PiHTopClients represents Pi-Hole client IPs with corresponding number of requests
type PiHTopClients struct {
	Clients map[string]int `json:"top_sources"`
}

// PiHForwardDestinations represents number of queries that have been forwarded and the target
type PiHForwardDestinations struct {
	Destinations map[string]float32 `json:"forward_destinations"`
}

// PiHQueryTypes contains DNS query type and number of queries
type PiHQueryTypes struct {
	Types map[string]float32 `json:"querytypes"`
}

// PiHQueries contains all DNS queries
// This is slice of slices of strings.
// Each slice contains: timestamp of query, type of query (IPv4, IPv6), requested DNS, requesting client, answer type
// Answer types: 1 = blocked by gravity.list, 2 = forwarded to upstream server, 3 = answered by local cache, 4 = blocked by wildcard blocking
type PiHQueries struct {
	Data [][]string `json:"data"`
}

// Get creates API request. Returns slice of bytes
func (r *PiHConnector) Get(endpoint string) []byte {
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

// Type returns Pi-Hole API type as an object
func (r *PiHConnector) Type() PiHType {
	bs := r.Get("type")
	s := &PiHType{}

	err := json.Unmarshal(bs, s)
	if err != nil {
		log.Fatal(err)
	}
	return *s
}

// Version returns Pi-Hole API version as an object
func (r *PiHConnector) Version() PiHVersion {
	bs := r.Get("version")
	s := &PiHVersion{}

	err := json.Unmarshal(bs, s)
	if err != nil {
		log.Fatal(err)
	}
	return *s
}

// SummaryRaw returns Pi-Hole's raw summary statistics
func (r *PiHConnector) SummaryRaw() PiHSummaryRaw {
	bs := r.Get("summaryRaw")
	s := &PiHSummaryRaw{}

	err := json.Unmarshal(bs, s)
	if err != nil {
		log.Fatal(err)
	}
	return *s
}

// Summary returns statistics in formatted style
func (r *PiHConnector) Summary() PiHSummary {
	bs := r.Get("summary")
	s := &PiHSummary{}

	err := json.Unmarshal(bs, s)
	if err != nil {
		log.Fatal(err)
	}
	return *s
}

// TimeData returns PiHTimeData object which contains requests statistics
func (r *PiHConnector) TimeData() PiHTimeData {
	bs := r.Get("overTimeData10mins")
	s := &PiHTimeData{}

	err := json.Unmarshal(bs, s)
	if err != nil {
		log.Fatal(err)
	}
	return *s
}

// Top returns top blocked and requested domains
func (r *PiHConnector) Top(n int) PiHTopItems {
	bs := r.Get("topItems=" + strconv.Itoa(n))
	s := &PiHTopItems{}

	err := json.Unmarshal(bs, s)
	if err != nil {
		log.Fatal(err)
	}
	return *s
}

// Clients returns top clients
func (r *PiHConnector) Clients(n int) PiHTopClients {
	bs := r.Get("topClients=" + strconv.Itoa(n))
	s := &PiHTopClients{}

	err := json.Unmarshal(bs, s)
	if err != nil {
		log.Fatal(err)
	}
	return *s
}

// ForwardDestinations returns forward destinations (DNS servers)
func (r *PiHConnector) ForwardDestinations() PiHForwardDestinations {
	bs := r.Get("getForwardDestinations")
	s := &PiHForwardDestinations{}

	err := json.Unmarshal(bs, s)
	if err != nil {
		log.Fatal(err)
	}
	return *s
}

// QueryTypes returns DNS query type and frequency as a PiHQueryTypes object
func (r *PiHConnector) QueryTypes() PiHQueryTypes {
	bs := r.Get("getQueryTypes")
	s := &PiHQueryTypes{}

	err := json.Unmarshal(bs, s)
	if err != nil {
		log.Fatal(err)
	}
	return *s
}

// Queries returns all DNS queries as a PiHQueries object
func (r *PiHConnector) Queries() PiHQueries {
	bs := r.Get("getAllQueries")
	s := &PiHQueries{}

	err := json.Unmarshal(bs, s)
	if err != nil {
		log.Fatal(err)
	}
	return *s
}

// Enable enables Pi-Hole server
func (r *PiHConnector) Enable() error {
	bs := r.Get("enable")
	resp := make(map[string]string)

	err := json.Unmarshal(bs, &resp)
	if err != nil {
		log.Fatal(err)
	}

	if resp["status"] != "enabled" {
		return errors.New("Something went wrong")
	}
	return nil
}

// Disable disables Pi-Hole server permanently
func (r *PiHConnector) Disable() error {
	bs := r.Get("disable")
	resp := make(map[string]string)

	err := json.Unmarshal(bs, &resp)
	if err != nil {
		log.Fatal(err)
	}

	if resp["status"] != "disabled" {
		return errors.New("Something went wrong")
	}
	return nil
}

// RecentBlocked returns string with the last blocked DNS record
func (r *PiHConnector) RecentBlocked() string {
	bs := r.Get("recentBlocked")
	return string(bs)
}

// Show returns 24h Summary of PiHole System
func (r *PiHSummaryRaw) Show() {
	fmt.Println("=== 24h Summary:")
	fmt.Printf("- Blocked Domains: %d\n", r.AdsBlocked)
	fmt.Printf("- Blocked Percentage: %.2f%%\n", r.AdsPercentage)
	fmt.Printf("- Queries: %d\n", r.DNSQueries)
	fmt.Printf("- Clients Ever Seen: %d\n", r.ClientsEverSeen)
}

// ShowBlocked returns sorted top Blocked domains over last 24h
func (r *PiHTopItems) ShowBlocked() {
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
func (r *PiHTopItems) ShowQueries() {
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

// Show returns sorted top clients over last 24h
func (r *PiHTopClients) Show() {
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
