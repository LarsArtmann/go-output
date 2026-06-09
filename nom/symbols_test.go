package nom

import (
	"testing"
)

func TestSymbolConstants(t *testing.T) {
	t.Parallel()

	symbols := map[string]string{
		"SymbolRunning":   SymbolRunning,
		"SymbolCompleted": SymbolCompleted,
		"SymbolPaused":    SymbolPaused,
		"SymbolFailed":    SymbolFailed,
		"SymbolDownload":  SymbolDownload,
		"SymbolUpload":    SymbolUpload,
		"SymbolTiming":    SymbolTiming,
		"SymbolAverage":   SymbolAverage,
		"SymbolTotal":     SymbolTotal,
	}

	for name, value := range symbols {
		t.Run(name, func(t *testing.T) {
			t.Parallel()

			if value == "" {
				t.Errorf("%s should not be empty", name)
			}
		})
	}
}

func TestOperationTypeConstants(t *testing.T) {
	t.Parallel()

	if OperationTypeDownload != "download" {
		t.Errorf("OperationTypeDownload = %q, want %q", OperationTypeDownload, "download")
	}

	if OperationTypeUpload != "upload" {
		t.Errorf("OperationTypeUpload = %q, want %q", OperationTypeUpload, "upload")
	}
}
