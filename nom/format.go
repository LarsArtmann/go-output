package nom

import (
	"fmt"
	"time"
)

// ============================================================================
// TIMING FORMATTING HELPERS
// ============================================================================
// TimingFormat is a legacy constant. FormatDuration now handles all formatting
// internally and does not use this constant.
//
// Deprecated: Use FormatDuration() for duration display. This constant will be
// removed in v2.
//
//nolint:staticcheck // kept for backward compatibility, remove in v2
const TimingFormat = "%.1fs"

// FormatDuration formats a duration in NOM style (e.g., "1.5s", "2m30s").
func FormatDuration(duration time.Duration) string {
	if duration < time.Second {
		return fmt.Sprintf("%.0fms", float64(duration.Milliseconds()))
	}

	if duration < time.Minute {
		// Use integer math (tenths of a second) to avoid float rounding
		// across the minute boundary (59.95s would display as "60.0s" with %.1f).
		tenths := int(duration / (100 * time.Millisecond))
		return fmt.Sprintf("%d.%ds", tenths/10, tenths%10)
	}

	if duration < time.Hour {
		minutes := int(duration.Minutes())
		seconds := int(duration.Seconds()) % 60

		if seconds == 0 {
			return fmt.Sprintf("%dm", minutes)
		}

		return fmt.Sprintf("%dm%ds", minutes, seconds)
	}

	// Hours: avoid unwieldy "90m" or "1440m" for long-running workflows.
	hours := int(duration.Hours())
	remaining := duration - time.Duration(hours)*time.Hour
	minutes := int(remaining.Minutes())

	if minutes == 0 {
		return fmt.Sprintf("%dh", hours)
	}

	return fmt.Sprintf("%dh%dm", hours, minutes)
}

// ShouldDisplayTiming determines if timing should be displayed.
// Only shows timing if duration >= 1 second to reduce noise.
func ShouldDisplayTiming(duration time.Duration) bool {
	return duration >= time.Second
}

// FormatActivityNodeTiming formats timing information for a ActivityNode in NOM style.
// Includes conditional display (hide if < 1s) and appropriate symbols.
func FormatActivityNodeTiming(status ActivityStatus, elapsed, estimated time.Duration) string {
	switch status {
	case ActivityStatusRunning, ActivityStatusCompleted, ActivityStatusFailed:
		if elapsed > 0 && ShouldDisplayTiming(elapsed) {
			return FormatDuration(elapsed)
		}
	case ActivityStatusPending:
		if estimated > 0 && ShouldDisplayTiming(estimated) {
			return formatTimingWithSymbol(estimated, SymbolAverage)
		}
	}

	return ""
}

func formatTimingWithSymbol(d time.Duration, symbol Symbol) string {
	return fmt.Sprintf("%s%s", symbol, FormatDuration(d))
}
