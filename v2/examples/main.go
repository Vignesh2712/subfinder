package main

import (
	"bytes"
	"context"
	"io"
	"log"
	"strconv"

	"github.com/projectdiscovery/subfinder/v2/pkg/runner"
	contextutil "github.com/projectdiscovery/utils/context"
)

func main() {
	subfinderOpts := &runner.Options{
		Threads:            10, // Thread controls the number of threads to use for active enumerations
		Timeout:            30, // Timeout is the seconds to wait for sources to respond
		MaxEnumerationTime: 10, // MaxEnumerationTime is the maximum amount of time in mins to wait for enumeration
		// ResultCallback: func(s *resolve.HostEntry) {
		// callback function executed after each unique subdomain is found
		// },
		// ProviderConfig: "your_provider_config.yaml",
		// and other config related options
	}

	// disable timestamps in logs / configure logger
	log.SetFlags(0)

	subfinder, err := runner.NewRunner(subfinderOpts)
	if err != nil {
		log.Fatalf("failed to create subfinder runner: %v", err)
	}

	output := &bytes.Buffer{}
	// To run subdomain enumeration on a single domain
	ctx, _ := contextutil.WithValues(context.Background(), contextutil.ContextArg("All"), contextutil.ContextArg(strconv.FormatBool(subfinderOpts.All)))
	if err = subfinder.EnumerateSingleDomainWithCtx(ctx, "hackerone.com", []io.Writer{output}); err != nil {
		log.Fatalf("failed to enumerate single domain: %v", err)
	}

	// To run subdomain enumeration on a list of domains from file/reader
	// file, err := os.Open("domains.txt")
	// if err != nil {
	// 	log.Fatalf("failed to open domains file: %v", err)
	// }
	// defer file.Close()
	// if err = subfinder.EnumerateMultipleDomainsWithCtx(context.Background(), file, []io.Writer{output}); err != nil {
	// 	log.Fatalf("failed to enumerate subdomains from file: %v", err)
	// }

	// print the output
	log.Println(output.String())
}
