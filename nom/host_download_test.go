package nom

import (
	"context"
	"strings"
	"testing"
)

func TestFormatActivityLabel_HostTag(t *testing.T) {
	t.Parallel()

	snap := ActivitySnapshot{
		Label:  "Build",
		Status: ActivityStatusRunning,
		Symbol: ActivityStatusRunning.GetSymbol(),
		Color:  ActivityStatusRunning.GetColor(),
		Host:   "builder-1",
	}

	display, _ := formatActivityLabel(snap)
	if !strings.Contains(display, "@builder-1") {
		t.Errorf("expected host tag @builder-1 in display, got: %q", display)
	}
}

func TestFormatActivityLabel_HostDormantWhenEmpty(t *testing.T) {
	t.Parallel()

	snap := ActivitySnapshot{
		Label:  "Build",
		Status: ActivityStatusRunning,
		Symbol: ActivityStatusRunning.GetSymbol(),
		Color:  ActivityStatusRunning.GetColor(),
	}

	display, _ := formatActivityLabel(snap)
	if strings.Contains(display, "@") {
		t.Errorf("host tag should be absent when Host is empty, got: %q", display)
	}
}

func TestFormatActivityLabel_DownloadBar(t *testing.T) {
	t.Parallel()

	snap := ActivitySnapshot{
		Label:    "Fetch",
		Status:   ActivityStatusRunning,
		Symbol:   ActivityStatusRunning.GetSymbol(),
		Color:    ActivityStatusRunning.GetColor(),
		Download: DownloadProgress{Downloaded: 50, Total: 100},
	}

	display, _ := formatActivityLabel(snap)
	if !strings.Contains(display, "50%") {
		t.Errorf("expected 50%% in download bar, got: %q", display)
	}

	if !strings.Contains(display, "█") || !strings.Contains(display, "░") {
		t.Errorf("expected filled/empty bar glyphs, got: %q", display)
	}
}

func TestFormatActivityLabel_DownloadBarDormantWhenCompleted(t *testing.T) {
	t.Parallel()

	// A completed activity should not show a download bar even if Download is set.
	snap := ActivitySnapshot{
		Label:    "Fetch",
		Status:   ActivityStatusCompleted,
		Symbol:   ActivityStatusCompleted.GetSymbol(),
		Color:    ActivityStatusCompleted.GetColor(),
		Download: DownloadProgress{Downloaded: 100, Total: 100},
	}

	display, _ := formatActivityLabel(snap)
	if strings.Contains(display, "▕") {
		t.Errorf("completed activity should not render a download bar, got: %q", display)
	}
}

func TestDownloadProgress_Fraction(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name string
		d    DownloadProgress
		want float64
	}{
		{"zero total", DownloadProgress{Downloaded: 50, Total: 0}, 0},
		{"half", DownloadProgress{Downloaded: 50, Total: 100}, 0.5},
		{"complete", DownloadProgress{Downloaded: 100, Total: 100}, 1},
		{"over", DownloadProgress{Downloaded: 150, Total: 100}, 1},
	}

	for _, tc := range tests {
		if got := tc.d.Fraction(); got != tc.want {
			t.Errorf("%s: Fraction = %v, want %v", tc.name, got, tc.want)
		}
	}
}

// hostDownloadEvent removed: ActivityStarted now carries Host and Download
// fields directly, so test events construct ActivityStarted{Host: ..., Download: ...}.

func TestSubscriber_PropagatesHostAndDownload(t *testing.T) {
	t.Parallel()

	sub := NewNOMSubscriber()
	ctx := context.Background()

	if err := sub.OnEvent(ctx, ActivityStarted{
		ID:       ActivityID("dl"),
		Name:     ActivityName("Download Deps"),
		Host:     "eu-west-1",
		Download: DownloadProgress{Downloaded: 700, Total: 1000},
	}); err != nil {
		t.Fatalf("OnEvent: %v", err)
	}

	snaps := sub.SnapshotActivities()

	snap, ok := snaps[ActivityID("dl")]
	if !ok {
		t.Fatal("activity not snapshotted")
	}

	if snap.Host != "eu-west-1" {
		t.Errorf("Host = %q, want eu-west-1", snap.Host)
	}

	if snap.Download.Total != 1000 || snap.Download.Downloaded != 700 {
		t.Errorf("Download = %+v, want {Downloaded:700 Total:1000}", snap.Download)
	}
}
