package main

import (
	"./clients"
	"flag"
	"log"
	"time"
)

const timeLayout = "2006-01-02T15:04:05Z"

func main() {
	startPtr := flag.String("start", time.Now().UTC().Add(-1*time.Hour).Format(timeLayout), "Start time in ISO 8601")
	endPtr := flag.String("end", time.Now().UTC().Format(timeLayout), "End time in ISO 8601")
	granularityPtr := flag.Int("granularity", 60, "Desired timeslice in seconds")
	environmentPtr := flag.String("environment", "sandbox", "Environment to execute the request in")
	flag.Parse()

	log.Printf("start = %v", *startPtr)
	log.Printf("end = %v", *endPtr)
	log.Printf("granularity = %v", *granularityPtr)
	log.Printf("environment = %v", *environmentPtr)

	var client *clients.Client = nil
	switch *environmentPtr {
	case "sandbox":
		client = clients.NewSandboxClient()
	case "production":
		client = clients.NewProductionClient()
	}

	var startTime, endTime *time.Time
	if nil != startPtr {
		t, err := time.Parse(timeLayout, *startPtr)
		startTime = &t
		if nil != err {
			log.Fatal(err)
		}
	}
	if nil != endPtr {
		t, err := time.Parse(timeLayout, *endPtr)
		endTime = &t
		if nil != err {
			log.Fatal(err)
		}
	}
	var granularity clients.HistoricRateGranularity = 60
	if nil != granularityPtr {
		granularity = clients.HistoricRateGranularity(*granularityPtr)
	}

	output, err := clients.GetProductHistoricRates(client, "ETH-USD", startTime, endTime, granularity)
	if nil != err {
		log.Fatal(err)
	}

	log.Printf("output = %v", output)
}
