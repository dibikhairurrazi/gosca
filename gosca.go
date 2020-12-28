package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"regexp"

	"github.com/dibikhairurrazi/gosca/rules"
)

const usageDoc = `Calculate cyclomatic complexities of Go functions.
Usage:
    gosca [flags] <Go file or directory> ...
Flags:
	-ignore REGEX         exclude files matching the given regular expression
	-max-cyclo N		  show function with cyclomatic complexity of N or higher
	-max-cogni N 		  show function with cognitive complexity of N or higher
The output fields for each line are:
Function <function>() on package <package> have cyclomatic complexity of <complexity> (exceeding <N>) consider refactoring.
`

var (
	cycloThreshold     *int
	cognitiveThreshold *int
)

func main() {
	cycloThreshold = flag.Int("max-cyclo", 0, "show functions with complexity > N only")
	cognitiveThreshold = flag.Int("max-cogni", 0, "show functions with complexity > N only")
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
	if *cycloThreshold == 0 && *cognitiveThreshold == 0 {
		printAll(allStats)
	} else {
		printStats(allStats)
	}

	os.Exit(0)

	// go vet integration
	/* multichecker.Main(
		atomic.Analyzer,
		loopclosure.Analyzer,
	) */
}

func printStats(s rules.Stats) {
	for _, stat := range s {
		if stat.Cyclomatic > *cycloThreshold {
			fmt.Printf("Function %v() on package %v have cyclomatic complexity of %d (exceeding %d) consider refactoring.\n", stat.FuncName, stat.PkgName, stat.Cyclomatic, *cycloThreshold)
		}

		if stat.Cognitive > *cognitiveThreshold {
			fmt.Printf("Function %v() on package %v have cognitive complexity of %d (exceeding %d) consider refactoring.\n", stat.FuncName, stat.PkgName, stat.Cognitive, *cognitiveThreshold)
		}
	}
}

func printAll(s rules.Stats) {
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
	fmt.Fprintf(os.Stderr, usageDoc)
	os.Exit(2)
}
