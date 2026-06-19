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

func NewActivityID(s string) ActivityID     { return ActivityID(s) }
func NewWorkflowID(s string) WorkflowID     { return WorkflowID(s) }
func NewActivityName(s string) ActivityName { return ActivityName(s) }
func NewWorkflowName(s string) WorkflowName { return WorkflowName(s) }
