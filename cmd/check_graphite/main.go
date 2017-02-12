package main

import (
	"flag"
	"net/http"
	"net/url"

	"github.com/jasonhancock/go-nagios"
)

func main() {
	fs := flag.CommandLine
	p := nagios.NewPlugin("gaphite-data", fs)
	p.StringFlag("graphite", "", "Graphite render URL: http://localhost/graphite/render/")
	p.StringFlag("metric", "", "Target metric name")
	p.StringFlag("username", "", "Username for basic auth")
	p.StringFlag("password", "", "Password for basic auth")
	flag.Parse()

	graphite := p.OptRequiredString("graphite")
	url, err := url.Parse(graphite)
	if err != nil {
		p.Fatal("unable to parse graphite url")
	}

	metric := p.OptRequiredString("metric")
	username, _ := p.OptString("username")
	password, _ := p.OptString("password")

	v := url.Query()
	v.Add("format", "json")
	v.Add("target", metric)
	v.Add("from", "-15minutes")
	url.RawQuery = v.Encode()

	client := &http.Client{}

	p.Verbose(url.String())

	req, err := http.NewRequest("GET", url.String(), nil)
	if err != nil {
		p.Fatal(err)
	}

	if username != "" && password != "" {
		req.SetBasicAuth(username, password)
	}

	resp, err := client.Do(req)
	if err != nil {
		p.Fatal(err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		p.Fatalf("got non 200 http status from graphite: %d", resp.StatusCode)
	}

	value, err := parseGraphiteResponse(resp.Body, metric)
	if err != nil {
		p.Fatal(err)
	}

	code, err := p.CheckThresholds(value)
	if err != nil {
		p.Fatal(err)
	}

	p.Exit(code, "The value is %f", value)
}
