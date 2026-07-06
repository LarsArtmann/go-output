// Package nom provides NOM-style real-time progress visualization for
// long-running workflows, inspired by nix-output-monitor.
//
// The package is built around three main components:
//
//   - [NOMSubscriber] receives workflow/activity events via OnEvent and
//     maintains the live activity state (status, timing, dependency tree).
//   - [DependencyTree] holds the parent/child relationships between activities
//     and renders them as a UTF-8 box-drawing tree, sorted by priority
//     (failed > running > paused > pending > completed).
//   - [InlineRenderer] draws the tree to an io.Writer in-place using ANSI
//     escape codes — no alt-screen takeover. Each redraw moves the cursor up,
//     clears the previous frame, and repaints, so output from prior steps
//     remains scrollable above.
//
// # Quick start
//
//	sub := nom.NewNOMSubscriber()
//	renderer := nom.NewInlineRenderer(sub, os.Stdout, 20)
//	renderer.Start(ctx, 100*time.Millisecond)
//	defer renderer.Finish(nil)
//
//	sub.OnEvent(ctx, &workflowStarted{...})
//	sub.OnEvent(ctx, &activityStarted{...})
//	// ... activities run ...
//	sub.OnEvent(ctx, &workflowCompleted{...})
//
// # Concurrency model
//
// The subscriber and tree share the same *Activity pointers — mutations via
// SetRunning/SetCompleted/SetFailed are instantly visible to both without any
// sync call. All event handlers and rendering take the subscriber's RWMutex to
// prevent garbled frames. See ADR 007 for the composition rationale.
//
// # Timing cache
//
// Completed activity durations are persisted to ~/.cache/nom-timing.csv so
// future runs can display estimated times for pending activities (median of
// historical samples).
package nom
