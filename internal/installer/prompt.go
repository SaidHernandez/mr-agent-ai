package installer

import (
	"bufio"
	"fmt"
	"strconv"
	"strings"
)

// promptSelectOne shows a numbered list and returns a single selected index (-1 if skipped).
func promptSelectOne(reader *bufio.Reader, title string, items []string) int {
	fmt.Printf("%s\n\n", title)
	for i, item := range items {
		fmt.Printf("  %2d.  %s\n", i+1, item)
	}
	fmt.Print("\n  Enter a number (or press Enter to skip): ")

	input, _ := reader.ReadString('\n')
	input = strings.TrimSpace(input)

	if input == "" {
		return -1
	}

	n, err := strconv.Atoi(input)
	if err != nil || n < 1 || n > len(items) {
		return -1
	}
	return n - 1
}

// promptSelect shows a numbered list and returns the selected indices.
// Empty input or "all" selects everything.
func promptSelect(reader *bufio.Reader, title string, items []string) []int {
	fmt.Printf("%s\n\n", title)
	for i, item := range items {
		fmt.Printf("  %2d.  %s\n", i+1, item)
	}
	fmt.Print("\n  Enter numbers (e.g. 1,3 or 'all'): ")

	input, _ := reader.ReadString('\n')
	input = strings.TrimSpace(strings.ToLower(input))

	if input == "" || input == "all" {
		all := make([]int, len(items))
		for i := range items {
			all[i] = i
		}
		return all
	}

	seen := map[int]bool{}
	var indices []int
	for _, part := range strings.Split(input, ",") {
		n, err := strconv.Atoi(strings.TrimSpace(part))
		if err != nil || n < 1 || n > len(items) || seen[n-1] {
			continue
		}
		seen[n-1] = true
		indices = append(indices, n-1)
	}
	return indices
}
