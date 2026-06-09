package nom

import "errors"

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
func MustWorkflowID(s string) WorkflowID {
	if s == "" {
		panic("WorkflowID must not be empty")
	}

	return WorkflowID(s)
}

func MustActivityID(s string) ActivityID {
	if s == "" {
		panic("ActivityID must not be empty")
	}

	return ActivityID(s)
}
func NewActivityID(s string) ActivityID     { return ActivityID(s) }
func NewWorkflowID(s string) WorkflowID     { return WorkflowID(s) }
func NewActivityName(s string) ActivityName { return ActivityName(s) }
func NewWorkflowName(s string) WorkflowName { return WorkflowName(s) }
func ParseActivityID(s string) (ActivityID, error) {
	if s == "" {
		return "", errors.New("activity ID must not be empty")
	}

	return ActivityID(s), nil
}

func ParseWorkflowID(s string) (WorkflowID, error) {
	if s == "" {
		return "", errors.New("workflow ID must not be empty")
	}

	return WorkflowID(s), nil
}
