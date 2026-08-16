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
//	sub.OnEvent(ctx, nom.WorkflowStarted{ID: nom.NewWorkflowID("wf"), Name: nom.NewWorkflowName("Build")})
//	sub.OnEvent(ctx, nom.ActivityStarted{ID: nom.NewActivityID("build"), Name: nom.NewActivityName("Build")})
//	// ... activities run ...
//	sub.OnEvent(ctx, nom.WorkflowCompleted{})
//
// # Concurrency model
//
// The subscriber owns the mutable activity state; the dependency tree holds
// only IDs and topology. Renderers never share pointers with the subscriber:
// they consume [NOMSubscriber.SnapshotActivities], an immutable value copy
// taken under the subscriber's read lock, so every frame is a consistent
// point-in-time view. All event handlers take the subscriber's RWMutex.
// See ADR 007 for the composition rationale.
//
// # Timing cache
//
// Completed activity durations are persisted to ~/.cache/nom-timing.csv so
// future runs can display estimated times for pending activities (median of
// historical samples).
package nom
