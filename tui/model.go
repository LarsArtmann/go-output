package tui
import (
"time"
tea "charm.land/bubbletea/v2"
)
// ============================================================================
// BUBBLE TEA MODEL IMPLEMENTATION
// ============================================================================
// TickCmd creates a recurring tick command for the progress model.
func TickCmd() tea.Cmd {
	return tea.Tick(100*time.Millisecond, func(t time.Time) tea.Msg {
		return TickMsg(t)
	})
}
// Init implements tea.Model interface.
func (m *ProgressModel) Init() tea.Cmd {
	return TickCmd()
}
// Update implements tea.Model interface with unified state management.
func (m *ProgressModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		return m.handleWindowSize(msg)
	case tea.KeyPressMsg:
		return m.handleKeyPress(msg)
	case ProgressUpdateMsg:
		return m.handleProgressUpdate(msg)
	case TickMsg:
		return m.handleTick(msg)
	}
	return m, nil
}
func (m *ProgressModel) handleWindowSize(msg tea.WindowSizeMsg) (tea.Model, tea.Cmd) {
	m.width = msg.Width
	m.height = msg.Height
	return m, nil
}
func (m *ProgressModel) handleKeyPress(msg tea.KeyPressMsg) (tea.Model, tea.Cmd) {
	switch msg.String() {
	case "ctrl+c", "q":
		return m, tea.Quit
	}
	return m, nil
}
func (m *ProgressModel) handleProgressUpdate(msg ProgressUpdateMsg) (tea.Model, tea.Cmd) {
	if !m.workflowState.CanAcceptUpdates() {
		return m, nil
	}
	m.lastUpdate = time.Now()
	switch msg.Type {
	case ProgressUpdate:
		m.currentProgress = msg.Progress
		if msg.Progress >= 100.0 && m.workflowState.CanTransitionTo(WorkflowStateCompleted) {
			m.workflowState = WorkflowStateCompleted
		}
	case MessageUpdate:
		m.currentMessage = msg.Message
	case StepUpdate:
		// Step updates are handled in ReportStep method
	}
	return m, nil
}
func (m *ProgressModel) handleTick(msg TickMsg) (tea.Model, tea.Cmd) {
	if !m.workflowState.CanAcceptTicks() {
		return m, nil
	}
	m.lastUpdate = time.Time(msg)
	if m.nomSubscriber != nil && m.nomSubscriber.IsEnabled() {
		m.syncNOMSubscriber()
	}
	return m, TickCmd()
}
func (m *ProgressModel) syncNOMSubscriber() {
	m.nomSubscriber.UpdateRunningActivityElapsed()
	m.nomSubscriber.SyncActivityTimingToTree()
	m.activities = m.nomSubscriber.GetActivities()
	m.dependencyTree = m.nomSubscriber.GetDependencyTree()
	if m.nomSubscriber.IsWorkflowRunning() {
		if m.workflowState == WorkflowStateIdle {
			m.workflowState = WorkflowStateRunning
			m.startTime = m.nomSubscriber.GetStartTime()
		}
	} else if m.workflowState == WorkflowStateRunning {
		m.updateWorkflowCompletionState()
	}
}
func (m *ProgressModel) updateWorkflowCompletionState() {
	running, completed, failed, _ := m.nomSubscriber.GetActivityCounts()
	if failed > 0 {
		m.workflowState = WorkflowStateErrored
	} else if running == 0 && completed > 0 {
		m.workflowState = WorkflowStateCompleted
	}
}
