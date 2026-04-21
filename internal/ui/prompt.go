package ui

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"
)

// Select shows a numbered list and returns the selected indices.
// Empty input or "all" selects everything.
func Select(title string, items []string) []int {
	fmt.Printf("%s\n\n", title)
	for i, item := range items {
		fmt.Printf("  %2d.  %s\n", i+1, item)
	}
	fmt.Print("\n  Enter numbers (e.g. 1,3 or 'all'): ")

	reader := bufio.NewReader(os.Stdin)
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

// SelectOne shows a numbered list and returns a single selected index (-1 if skipped).
func SelectOne(title string, items []string) int {
	fmt.Printf("%s\n\n", title)
	for i, item := range items {
		fmt.Printf("  %2d.  %s\n", i+1, item)
	}
	fmt.Print("\n  Enter a number (or press Enter to skip): ")

	reader := bufio.NewReader(os.Stdin)
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
