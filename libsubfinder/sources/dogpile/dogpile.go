//
// Written By : @Mzack9999
//
// Distributed Under MIT License
// Copyrights (C) 2018 Ice3man
//

// A golang client for Dogpile Subdomain Discovery
package dogpile

import (
	"fmt"
	"io/ioutil"
	"math/rand"
	"regexp"
	"strconv"
	"time"

	"github.com/Ice3man543/subfinder/libsubfinder/helper"
)

// all subdomains found
var subdomains []string

// Query function returns all subdomains found using the service.
func Query(domain string, state *helper.State, ch chan helper.Result) {

	var result helper.Result
	result.Subdomains = subdomains
	maxPages, _ := strconv.Atoi(state.CurrentSettings.DogpilePages)
	for currentPage := 0; currentPage <= maxPages; currentPage++ {
		url := "http://www.dogpile.com/search/web?q=" + domain + "&qsi=" + strconv.Itoa(currentPage*15+1)

		resp, err := helper.GetHTTPResponse(url, state.Timeout)
		if err != nil {
			result.Error = err
			ch <- result
			return
		}

		// Get the response body
		body, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			result.Error = err
			ch <- result
			return
		}

		reSub := regexp.MustCompile(`%.{4}`)
		src := reSub.ReplaceAllLiteralString(string(body), " ")

		re := helper.SubdomainRegex(domain)
		match := re.FindAllString(src, -1)

		for _, subdomain := range match {
			if state.Verbose == true {
				if state.Color == true {
					fmt.Printf("\n[%sDogpile%s] %s", helper.Red, helper.Reset, subdomain)
				} else {
					fmt.Printf("\n[Dogpile] %s", subdomain)
				}
			}

			subdomains = append(subdomains, subdomain)
		}
		time.Sleep(time.Duration((3 + rand.Intn(5))) * time.Second)
	}

	result.Subdomains = subdomains
	result.Error = nil
	ch <- result
}
