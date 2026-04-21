package theme

import (
	"fmt"
	"os"
	"strings"
	"time"

	"golang.org/x/term"
)

// ANSI escape constants.
const (
	ansiReset     = "\033[0m"
	ansiBold      = "\033[1m"
	ansiDim       = "\033[2m"
	ansiGreen     = "\033[32m"
	ansiGray      = "\033[90m"
	ansiRed       = "\033[31m"
	ansiYellow    = "\033[33m"
	ansiBoldCyan  = "\033[1;36m"
)

// RenderPanel returns a full box-drawing panel with all widget results.
func RenderPanel(results []Result) string {
	const boxWidth = 47
	border := strings.Repeat("─", boxWidth-2)

	var b strings.Builder
	b.WriteString(ansiBold + "┌" + border + "┐" + ansiReset + "\n")
	b.WriteString(ansiBold + "│" + ansiReset)
	b.WriteString("  mr-agent-ai  context panel               ")
	b.WriteString(ansiBold + "│" + ansiReset + "\n")
	b.WriteString(ansiBold + "├" + border + "┤" + ansiReset + "\n")

	for _, r := range results {
		name := padRight(r.WidgetName, 10)
		value := formatValue(r)
		line := fmt.Sprintf("  %s%s%s  %s%-31s%s",
			ansiBoldCyan, name, ansiReset,
			valueColor(r), truncate(value, 31), ansiReset,
		)
		b.WriteString(ansiBold + "│" + ansiReset)
		b.WriteString(line)
		b.WriteString(ansiBold + "│" + ansiReset + "\n")
	}

	b.WriteString(ansiBold + "├" + border + "┤" + ansiReset + "\n")
	ts := time.Now().Format("15:04:05")
	footer := fmt.Sprintf("  %supdated %s%s", ansiDim, ts, ansiReset)
	b.WriteString(ansiBold + "│" + ansiReset)
	b.WriteString(footer)
	b.WriteString(strings.Repeat(" ", 47-2-len("  updated "+ts)))
	b.WriteString(ansiBold + "│" + ansiReset + "\n")
	b.WriteString(ansiBold + "└" + border + "┘" + ansiReset + "\n")

	return b.String()
}

// RenderStatusLine returns a compact single-line string for Claude Code's status bar.
// cols is the terminal width; the line is truncated to fit.
func RenderStatusLine(results []Result, cols int) string {
	var parts []string
	for _, r := range results {
		if r.Value == "" && r.Err == nil {
			continue
		}
		val := formatValue(r)
		if val == "(no data)" || val == "(no GITHUB_TOKEN)" || val == "(no LINEAR_API_KEY)" {
			continue
		}
		color := valueColor(r)
		parts = append(parts, fmt.Sprintf("%s%s: %s%s", color, r.WidgetName, val, ansiReset))
	}

	if len(parts) == 0 {
		return ""
	}

	sep := ansiGray + " | " + ansiReset
	line := " " + strings.Join(parts, sep) + " "

	// Truncate to terminal width (approximate, ignoring ANSI escape lengths)
	visible := stripANSI(line)
	if cols > 0 && len(visible) > cols {
		line = " " + strings.Join(parts[:len(parts)-1], sep) + " "
	}

	return line + "\n"
}

// termWidth returns the current terminal column count, defaulting to 80.
func termWidth() int {
	w, _, err := term.GetSize(int(os.Stdout.Fd())) // #nosec G115
	if err != nil || w <= 0 {
		return 80
	}
	return w
}

func formatValue(r Result) string {
	if r.Err != nil {
		msg := r.Err.Error()
		if len(msg) > 28 {
			msg = msg[:27] + "…"
		}
		return "[err] " + msg
	}
	if r.Value == "" {
		return "(no data)"
	}
	return r.Value
}

func valueColor(r Result) string {
	if r.Err != nil {
		return ansiRed
	}
	v := r.Value
	if v == "" || v == "(no data)" || strings.HasPrefix(v, "(no ") {
		return ansiGray
	}
	if strings.HasPrefix(v, "[err]") {
		return ansiRed
	}
	if strings.HasPrefix(v, "[timeout]") {
		return ansiYellow
	}
	return ansiGreen
}

func padRight(s string, n int) string {
	if len(s) >= n {
		return s[:n]
	}
	return s + strings.Repeat(" ", n-len(s))
}

func truncate(s string, max int) string {
	if len(s) <= max {
		return s
	}
	return s[:max-1] + "…"
}

// stripANSI removes ANSI escape sequences from s for length calculation.
func stripANSI(s string) string {
	var b strings.Builder
	inEsc := false
	for _, r := range s {
		if r == '\033' {
			inEsc = true
			continue
		}
		if inEsc {
			if r == 'm' {
				inEsc = false
			}
			continue
		}
		b.WriteRune(r)
	}
	return b.String()
}
