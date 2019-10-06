package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"github.com/jamiealquiza/envy"
	"github.com/jehiah/go-strftime"
	"io/ioutil"
	"log"
	"net"
	"net/http"
	"os"
	"sync"
	"time"
)

var (
	urlListJsonArray = flag.String("urls", "[]", "JSON array of URLs to DELETE. Each will be parsed with STRFTIME")
	timeOffset       = flag.Int("timeoffset", 8760, "STRFTIME will be <utc-now> minus this number of hours")
	loopDelay        = flag.Int("loopdelay", 0, "If specified: run continuously with this number of hours between actions. If unspecified: run once and exit")
	backCount        = flag.Int("backcount", 30, "Iterate back in time <backcount> times")
	backStep         = flag.Int("backstep", 24, "Each <backcount> iteration goes back in time <backstep> hours")
	concurrent       = flag.Int("concurrent", 2, "Number of HTTP requests to send in parallel")
	startupDelay     = flag.Int("startupdelay", 0, "Seconds to wait before doing anything (used to facilitate easier testing)")
	httpReqTimeout   = flag.Int("httptimeout", 30, "HTTP request timeout in seconds")
	dryRun           = flag.Bool("dryrun", false, "Show parsed target list without sending HTTP Delete request")
	wg               sync.WaitGroup
	urlTemplates     []string
)

func main() {

	gracefulStop()

	flag.Usage = func() {
		fmt.Fprintf(flag.CommandLine.Output(), "Usage of %s:\n", os.Args[0])
		flag.PrintDefaults()
		fmt.Fprintf(flag.CommandLine.Output(), "\nMIT License. Copyright (c) 2019 JamesJJ. https://github.com/JamesJJ/http-deleter\n")
	}

	envy.Parse("HTTP_DELETER")
	flag.Parse()

	time.Sleep(time.Duration(*startupDelay) * time.Second)

	deleteUrlChan := make(chan *string)

	for workers := 1; workers <= *concurrent; workers++ {
		log.Printf("Starting worker #%d", workers)
		wg.Add(1)
		go deleteWorker(&wg, deleteUrlChan, httpReqTimeout)
	}

	if errJson := json.Unmarshal([]byte(*urlListJsonArray), &urlTemplates); errJson != nil {
		flag.Usage()
		os.Exit(1)
	}

	for {
		for _, urlTemplate := range urlTemplates {
			for back := *backCount; back >= 0; back-- {

				targetTime := time.Now().UTC().Add(time.Hour * time.Duration(-*timeOffset)).Add(time.Hour * time.Duration(-*backStep*back))

				targetUrl := strftime.Format(urlTemplate, targetTime)

				if *dryRun {
					log.Printf("URL: %s", targetUrl)
				} else {
					deleteUrlChan <- &targetUrl
				}
			}
		}

		if *loopDelay == 0 {
			close(deleteUrlChan)
			break
		}
	}
	wg.Wait()
}

func deleteWorker(wg *sync.WaitGroup, inChan <-chan *string, timeout *int) {
	defer wg.Done()

	var url *string
	keepReading := true

	for keepReading {
		url, keepReading = <-inChan
		if keepReading {
			doDelete(url, timeout)
		}
	}
}

func doDelete(url *string, timeout *int) {

	var netTransport = &http.Transport{
		Dial: (&net.Dialer{
			Timeout: 4 * time.Second,
		}).Dial,
		TLSHandshakeTimeout: 4 * time.Second,
	}

	var netClient = &http.Client{
		Transport: netTransport,
		Timeout:   time.Second * time.Duration(*timeout),
	}

	req, errNewReq := http.NewRequest("DELETE", *url, nil)
	if errNewReq != nil {
		log.Printf("New Req Error: %v (%v)", *url, errNewReq)
		return
	}
	req.Header.Add("User-Agent", "HTTP-DELETER/1.0")

	res, errClientDo := netClient.Do(req)
	if errClientDo != nil {
		log.Printf("Client Do Error: %v (%v)", *url, errClientDo)
		return
	}
	defer res.Body.Close()

	body, errBody := ioutil.ReadAll(res.Body)
	if errBody != nil {
		log.Printf("Body Error: %v (%v)", *url, errBody)
	}
	log.Printf("DELETED: %s (%s)", *url, string(body))

}
