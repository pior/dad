package tasks

type genericTaskActionCondition struct {
	pre  func(*Context) *actionResult
	post func(*Context) *actionResult
}

type genericTaskActionBuilder struct {
	desc string

	conditions     []*genericTaskActionCondition
	monitoredFiles []string
	monitoredTasks []string

	runFunc func(*Context) error
}

func actionBuilder(description string, runFunc func(*Context) error) *genericTaskActionBuilder {
	return &genericTaskActionBuilder{desc: description, runFunc: runFunc}
}

// On registers a new condition
func (a *genericTaskActionBuilder) On(condition *genericTaskActionCondition) *genericTaskActionBuilder {
	a.conditions = append(a.conditions, condition)
	return a
}

// OnFunc registers a condition defined as a single function
func (a *genericTaskActionBuilder) OnFunc(condFunc func(*Context) *actionResult) *genericTaskActionBuilder {
	a.On(&genericTaskActionCondition{pre: condFunc, post: condFunc})
	return a
}

// OnFileChange specifies that the action will run when a file changes or does not exist.
// The action will NOT fail if the file is not created.
func (a *genericTaskActionBuilder) OnFileChange(path string) *genericTaskActionBuilder {
	a.monitoredFiles = append(a.monitoredFiles, path)
	return a
}

// OnTaskChange specifies that the action will run when a task changes.
func (a *genericTaskActionBuilder) OnTaskChange(name string) *genericTaskActionBuilder {
	a.monitoredTasks = append(a.monitoredTasks, name)
	return a
}

// Build returns a new task action with the behaviour specified by the builder
func (a *genericTaskActionBuilder) Build() *genericTaskAction {
	return &genericTaskAction{
		builder: &genericTaskActionBuilder{ // Hand made copy
			desc:           a.desc,
			conditions:     append(a.conditions[:0:0], a.conditions...),
			monitoredFiles: append(a.monitoredFiles[:0:0], a.monitoredFiles...),
			monitoredTasks: append(a.monitoredTasks[:0:0], a.monitoredTasks...),
			runFunc:        a.runFunc,
		},
	}
}
