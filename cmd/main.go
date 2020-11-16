package main

import (
	"flag"
	"fmt"
	"go-sca/rules"
	"log"
	"os"
	"regexp"
)

const usageDoc = `Calculate cyclomatic complexities of Go functions.
Usage:
    gosca [flags] <Go file or directory> ...
Flags:
    -over N               show functions with complexity > N only and
                          return exit code 1 if the set is non-empty
    -top N                show the top N most complex functions only
    -avg, -avg-short      show the average complexity over all functions;
                          the short option prints the value without a label
    -total, -total-short  show the total complexity for all functions;
                          the short option prints the value without a label
    -ignore REGEX         exclude files matching the given regular expression
The output fields for each line are:
<complexity> <package> <function> <file:line:column>
`

func main() {
	ignore := flag.String("ignore", "", "exclude files matching the given regular expression")

	log.SetFlags(0)
	log.SetPrefix("gosca: ")
	flag.Usage = usage
	flag.Parse()
	paths := flag.Args()
	if len(paths) == 0 {
		usage()
	}
	allStats := rules.Analyze(paths, regex(*ignore))
	printStats(allStats)
}

func printStats(s rules.Stats) {
	for _, stat := range s {
		fmt.Println(stat)
	}
}

func regex(expr string) *regexp.Regexp {
	if expr == "" {
		return nil
	}
	re, err := regexp.Compile(expr)
	if err != nil {
		log.Fatal(err)
	}
	return re
}

func usage() {
	_, _ = fmt.Fprintf(os.Stderr, usageDoc)
	os.Exit(2)
}
