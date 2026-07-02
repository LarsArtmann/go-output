package nom

type ActivityID string

func (id ActivityID) String() string { return string(id) }
func (id ActivityID) IsZero() bool   { return id == "" }

type ActivityName string

func (n ActivityName) String() string { return string(n) }

type WorkflowID string

func (id WorkflowID) String() string { return string(id) }
func (id WorkflowID) IsZero() bool   { return id == "" }

type WorkflowName string

func (n WorkflowName) String() string { return string(n) }

// ActivityCategory is a user-defined semantic group (e.g. "build", "test",
// "deploy") used for tinting and collapsing related activities.
type ActivityCategory string

// RenderMode selects how the dependency tree is displayed.
type RenderMode int

const (
	// RenderModeTree draws the dependency tree with parent/child connectors,
	// matching nix-output-monitor's layout.
	RenderModeTree RenderMode = iota
	// RenderModeLayered groups activities by their DAG depth and renders each
	// layer horizontally, making parallel work explicit.
	RenderModeLayered
)

// String returns the human-readable name of the render mode: "tree" or "layered".
func (m RenderMode) String() string {
	switch m {
	case RenderModeTree:
		return "tree"
	case RenderModeLayered:
		return "layered"
	default:
		return "unknown"
	}
}

func NewActivityID(s string) ActivityID     { return ActivityID(s) }
func NewWorkflowID(s string) WorkflowID     { return WorkflowID(s) }
func NewActivityName(s string) ActivityName { return ActivityName(s) }
func NewWorkflowName(s string) WorkflowName { return WorkflowName(s) }
