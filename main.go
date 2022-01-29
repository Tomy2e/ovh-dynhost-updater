package main

import (
	"context"
	"flag"
	"fmt"
	"io"
	"log"
	"net"
	"net/http"
	"net/url"
	"os"
	"strings"
	"time"
)

const (
	IPv4API               = "https://api.ipify.org/"
	MaxIPV4StringLen      = 15
	MaxDynHostResponseLen = 64
	DefaultDaemonDelay    = 10
)

func whatsMyIPv4(ctx context.Context) (net.IP, error) {
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, IPv4API, nil)
	if err != nil {
		return nil, err
	}

	res, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	data, err := io.ReadAll(http.MaxBytesReader(nil, res.Body, MaxIPV4StringLen))
	if err != nil {
		return nil, err
	}

	ip := net.ParseIP(string(data))

	if ip == nil {
		return nil, fmt.Errorf("failed to decode IP: %s", string(data))
	}

	if ip.To4() == nil {
		return nil, fmt.Errorf("not an IPv4: %s", string(data))
	}

	return ip, nil
}

func updateDynHostRecord(ctx context.Context, login, password, hostname string, ip net.IP) error {
	params := url.Values{
		"system":   {"dyndns"},
		"hostname": {hostname},
		"ip":       {ip.String()},
	}

	req, err := http.NewRequestWithContext(
		ctx,
		http.MethodGet,
		fmt.Sprintf("https://www.ovh.com/nic/update?%s", params.Encode()),
		nil,
	)
	if err != nil {
		return err
	}

	req.SetBasicAuth(login, password)

	res, err := http.DefaultClient.Do(req)
	if err != nil {
		return err
	}
	defer res.Body.Close()

	if res.StatusCode != http.StatusOK {
		return fmt.Errorf("unexpected status code %d", res.StatusCode)
	}

	data, err := io.ReadAll(http.MaxBytesReader(nil, res.Body, MaxDynHostResponseLen))
	if err != nil {
		return err
	}

	sd := string(data)

	if !strings.HasPrefix(sd, "nochg") && !strings.HasPrefix(sd, "good") {
		return fmt.Errorf("failed to update DynHost record: %s", sd)
	}

	return nil
}

func main() {
	ctx := context.Background()

	// Flag arguments
	hostname := flag.String("hostname", "", "Hostname of the DynHost record to update.")
	delay := flag.Uint("delay", DefaultDaemonDelay, "If the program is run as a daemon, number "+
		"of minutes to wait between updates.")
	daemon := flag.Bool("daemon", false, "Run the program as a daemon.")
	flag.Parse()

	if *hostname == "" {
		log.Fatal("Missing -hostname flag.")
	}

	// Env arguments
	login := os.Getenv("LOGIN")
	password := os.Getenv("PASSWORD")

	if login == "" {
		log.Fatal("Missing LOGIN environment variable.")
	}

	if password == "" {
		log.Fatal("Missing PASSWORD environment variable.")
	}

	// Lookup current DynHost IP
	var dynHostIP net.IP

	ips, err := net.LookupIP(*hostname)
	if err != nil {
		log.Fatalf("Failed to lookup IP: %v", err)
	}

	if len(ips) == 1 {
		dynHostIP = ips[0]
		log.Printf("Current IP of %s is %s\n", *hostname, dynHostIP)
	}

	// Run main logic.
	for {
		currentIP, err := whatsMyIPv4(ctx)
		if err != nil {
			log.Fatalf("Failed to get current IP: %v", err)
		}

		log.Println("Current IP of client is", currentIP)

		if !currentIP.Equal(dynHostIP) {
			log.Println("DynHost record needs to be updated")
			if err := updateDynHostRecord(ctx, login, password, *hostname, currentIP); err != nil {
				log.Fatalf("Failed to update DynHost record: %v", err)
			}
			dynHostIP = currentIP
		}

		// Escape early when not running as a daemon.
		if !*daemon {
			break
		}

		time.Sleep(time.Duration(*delay) * time.Minute)
	}
}
