package rules

import (
	"fmt"
	"go/token"
	"sort"
)

// Stat is statistic of the complexity.
type Stat struct {
	PkgName    string
	FuncName   string
	Cyclomatic int
	Cognitive  int
	Complexity int
	Pos        token.Position
}

func (s Stat) String() string {
	return fmt.Sprintf("cyclomatic: %d, cognitive: %d, %s %s %s", s.Cyclomatic, s.Cognitive, s.PkgName, s.FuncName, s.Pos)
}

// Stats hold the complexities of many functions.
type Stats []Stat

// AverageComplexity calculates the average complexity of the complexities in s.
func (s Stats) AverageComplexity() float64 {
	return float64(s.TotalComplexity()) / float64(len(s))
}

// TotalComplexity calculates the total sum of all complexities in s.
func (s Stats) TotalComplexity() uint64 {
	total := uint64(0)
	for _, stat := range s {
		total += uint64(stat.Complexity)
	}
	return total
}

// SortAndFilter sorts the complexities in s in descending order
// and returns a slice of s limited to the 'top' N entries with a complexity
// greater than 'over'. If 'top' is negative, i.e. -1, it does
// not limit the result. If 'over' is <= 0 it does not limit the result either,
// because a function has a base complexity of at least 1.
func (s Stats) SortAndFilter(top, over int) Stats {
	result := make(Stats, len(s))
	copy(result, s)
	sort.Sort(byComplexityDesc(result))
	for i, stat := range result {
		if i == top {
			return result[:i]
		}
		if stat.Complexity <= over {
			return result[:i]
		}
	}
	return result
}

type byComplexityDesc Stats

func (s byComplexityDesc) Len() int {
	return len(s)
}

func (s byComplexityDesc) Swap(i, j int) {
	s[i], s[j] = s[j], s[i]
}

func (s byComplexityDesc) Less(i, j int) bool {
	return s[i].Complexity >= s[j].Complexity
}
