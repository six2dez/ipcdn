package main

import (
	"bufio"
	"crypto/tls"
	"flag"
	"fmt"
	"log"
	"net"
	"net/http"
	"os"
	"strings"
	"sync"
	"time"

	jsoniter "github.com/json-iterator/go"
	"github.com/yl2chen/cidranger"
)

// Client checks for CDN based IPs which should be excluded
// during scans since they belong to third party firewalls.
type Client struct {
	Options *Options
	ranges  map[string][]string
	rangers map[string]cidranger.Ranger
}

type Options struct {
	Cache bool
}

// New creates a new firewall IP checking client.
func New() (*Client, error) {
	return new(&Options{})
}

// NewWithCache creates a new firewall IP with cached data from project discovery (faster)
func NewWithCache() (*Client, error) {
	return new(&Options{Cache: true})
}

func new(options *Options) (*Client, error) {
	httpClient := &http.Client{
		Transport: &http.Transport{
			MaxIdleConns:        100,
			MaxIdleConnsPerHost: 100,
			TLSClientConfig: &tls.Config{
				Renegotiation:      tls.RenegotiateOnceAsClient,
				InsecureSkipVerify: true,
			},
		},
		Timeout: time.Duration(30) * time.Second,
	}
	client := &Client{Options: options}

	err := client.getCDNDataFromCache(httpClient)
	if err != nil {
		return nil, err
	}
	return client, nil
}

func scrapeProjectDiscovery(httpClient *http.Client) (map[string][]string, error) {

	var (
		err          error
		response     *http.Response
		retries      int = 2
		url          string
		url_fallback string
		data         map[string][]string
	)

	url = "https://cdn.nuclei.sh"
	url_fallback = "https://raw.githubusercontent.com/six2dez/ipcdn/main/ranges.txt"

	for retries > 0 {
		response, err = http.Get(url)
		if err != nil {
			url = url_fallback
			retries -= 1
		} else {
			break
		}
	}
	if response != nil {
		defer response.Body.Close()

		if err := jsoniter.NewDecoder(response.Body).Decode(&data); err != nil {
			return nil, err
		}
	} else {
		return nil, err
	}
	return data, nil
}

func (c *Client) getCDNDataFromCache(httpClient *http.Client) error {
	var err error
	c.ranges, err = scrapeProjectDiscovery(httpClient)
	if err != nil {
		return err
	}

	c.rangers = make(map[string]cidranger.Ranger)
	for provider, ranges := range c.ranges {
		ranger := cidranger.NewPCTrieRanger()

		for _, cidr := range ranges {
			_, network, err := net.ParseCIDR(cidr)
			if err != nil {
				continue
			}
			_ = ranger.Insert(cidranger.NewBasicRangerEntry(*network))
		}
		c.rangers[provider] = ranger
	}
	return nil
}

// Check checks if an IP is contained in the blacklist
func (c *Client) Check(ip net.IP) (bool, string, error) {
	for provider, ranger := range c.rangers {
		if contains, err := ranger.Contains(ip); contains {
			return true, provider, err
		}
	}
	return false, "", nil
}

// Ranges returns the providers and ranges for the cdn client
func (c *Client) Ranges() map[string][]string {
	return c.ranges
}

func main() {

	var verbose bool
	flag.BoolVar(&verbose, "v", false, "verbose mode")

	var method string
	flag.StringVar(&method, "m", "cdn", "Output method:\n    cdn - Only prints IPs on CDN\n    not - only prints NOT on CDN\n    all - prints all with verbose\n")

	flag.Parse()

	var rangesWG sync.WaitGroup
	ips := make(chan string)
	output := make(chan string)
	var outputWG sync.WaitGroup

	rangesWG.Add(1)
	go func() {
		for ip := range ips {
			if isListening(ip, verbose, method) {
				output <- ip
				continue
			}
		}
		rangesWG.Done()
	}()

	go func() {
		for o := range output {
			if !verbose {
				if method != "all" {
					_ = o
				}
			}
		}
		outputWG.Done()
	}()

	// Close the output channel when the HTTP workers are done
	go func() {
		rangesWG.Wait()
		close(output)
	}()

	sc := bufio.NewScanner(os.Stdin)
	for sc.Scan() {
		ip := strings.ToLower(sc.Text())

		ips <- ip

	}
}

func isListening(ip string, verbose bool, method string) bool {

	client, err := NewWithCache()
	if err != nil {
		log.Fatal(err)
	}

	if found, provider, err := client.Check(net.ParseIP(ip)); found && provider != "" && err == nil {
		if verbose && (method == "all" || method == "cdn") {
			fmt.Println(ip + "-" + provider)
		} else if method == "all" || method == "cdn" {
			fmt.Println(ip)
		}

	} else {
		provider := "not CDN"
		if verbose && (method == "all" || method == "not") {
			fmt.Println(ip + "-" + provider)
		} else if method == "all" || method == "not" {
			fmt.Println(ip)
		}
	}

	if err != nil {
		return false
	}

	return true
}
