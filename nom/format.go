package nom

import (
	"fmt"
	"strings"
	"time"
)

// ============================================================================
// OPERATION TYPE SYMBOLS
// ============================================================================
// OperationSymbol returns the symbol for an operation type.
func OperationSymbol(operationType string) string {
	switch operationType {
	case OperationTypeDownload:
		return SymbolDownload
	case OperationTypeUpload:
		return SymbolUpload
	default:
		return ""
	}
}

const (
	OperationTypeDownload = "download"
	OperationTypeUpload   = "upload"
)

// ============================================================================
// TIMING FORMATTING HELPERS
// ============================================================================
// TimingFormat is the format string for displaying timing (NOM-style).
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

	minutes := int(duration.Minutes())

	seconds := int(duration.Seconds()) % 60
	if seconds == 0 {
		return fmt.Sprintf("%dm", minutes)
	}

	return fmt.Sprintf("%dm%ds", minutes, seconds)
}

// FormatTimingInfo formats timing information for an activity.
func FormatTimingInfo(state *ActivityDisplayState) string {
	if state.IsRunning() {
		elapsed := time.Since(state.StartTime)
		return fmt.Sprintf("%s%s", SymbolTiming, FormatDuration(elapsed))
	}

	if state.IsCompleted() && !state.StartTime.IsZero() && !state.EndTime.IsZero() {
		duration := state.EndTime.Sub(state.StartTime)
		return fmt.Sprintf("%s%s", SymbolTiming, FormatDuration(duration))
	}

	if state.EstimatedTime > 0 {
		return fmt.Sprintf("%s%s", SymbolAverage, FormatDuration(state.EstimatedTime))
	}

	return ""
}

// GetActivitySummaryString generates a summary string for multiple activities
// Format: "⏵3↑2∑5" (3 running, 2 uploading, 5 total).
func GetActivitySummaryString(running, uploading, downloading, total int) string {
	var parts []string
	if running > 0 {
		parts = append(parts, fmt.Sprintf("%s%d", SymbolRunning, running))
	}

	if uploading > 0 {
		parts = append(parts, fmt.Sprintf("%s%d", SymbolUpload, uploading))
	}

	if downloading > 0 {
		parts = append(parts, fmt.Sprintf("%s%d", SymbolDownload, downloading))
	}

	if len(parts) > 0 {
		parts = append(parts, fmt.Sprintf("%s%d", SymbolTotal, total))
		return strings.Join(parts, "")
	}

	if total > 0 {
		return fmt.Sprintf("%s%d", SymbolTotal, total)
	}

	return ""
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
			return formatTimingWithSymbol(elapsed, SymbolTiming)
		}
	case ActivityStatusPending:
		if estimated > 0 && ShouldDisplayTiming(estimated) {
			return formatTimingWithSymbol(estimated, SymbolAverage)
		}
	case ActivityStatusPaused:
		if elapsed > 0 && ShouldDisplayTiming(elapsed) {
			return formatTimingWithSymbol(elapsed, SymbolPaused)
		}
	}

	return ""
}

func formatTimingWithSymbol(d time.Duration, symbol string) string {
	return fmt.Sprintf("%s%s", symbol, FormatDuration(d))
}
