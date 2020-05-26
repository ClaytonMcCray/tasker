package main

import (
	"testing"
)

func TestMakeTask(t *testing.T) {
	tsk := makeTask("name")

	if !tsk.active() {
		t.Errorf("task is inactive when it should be active")
	}

	if tsk.name != "name" {
		t.Errorf("task.name is %s but should be %s.", tsk.name, "name")
	}
}

func TestActiveDeactive(t *testing.T) {
	tsk := makeTask("atask")

	if !tsk.active() {
		t.Errorf("task is inactive falsely")
	}

	tsk.deactivate()
	if tsk.active() {
		t.Errorf("task is active after calling deactivate()")
	}

	tsk.stamp()
	if !tsk.active() {
		t.Errorf("task is inactive after stamping")
	}

	tsk.stamp()
	if !tsk.active() {
		t.Errorf("task is inactive after stamping")
	}
}

type mockIO struct {
	idx    int
	inputs []string
}

func (m *mockIO) Scan() bool {
	m.idx++
	return m.idx < len(m.inputs)
}

func (m *mockIO) Text() string {
	return m.inputs[m.idx]
}

func TestInputHandler(t *testing.T) {
	m := &mockIO{
		idx:    -1,
		inputs: []string{"JHU", "G&A", "0", "", "1"},
	}

	ts := make([]*task, 0)

	for range m.inputs {
		_, ts = inputHandler(m, ts)
	}

	if len(ts) != 2 {
		t.Errorf("len(task list) is %d but should be %d", len(ts), 2)
	}
}

func TestParseInputCharacter(t *testing.T) {
	inputs := []string{
		"a new task",
		"",
		"another new task",
		"",
		"1",
	}

	expectedTypeOfInputs := []int{
		newTask,
		refresh,
		newTask,
		refresh,
		taskIndex,
	}

	expectedNames := []string{
		"a new task",
		"",
		"another new task",
		"",
		"",
	}

	expectedActiveIdx := []int{
		-1,
		-1,
		-1,
		-1,
		1,
	}

	for i, taskName := range inputs {
		typeOf, name, active := parseInputCharacter(taskName)
		if typeOf != expectedTypeOfInputs[i] {
			t.Errorf("enum was %d but should be %d for %s", typeOf, expectedTypeOfInputs[i], name)
		}

		if name != expectedNames[i] {
			t.Errorf("name was %s but should be %s", name, expectedNames[i])
		}

		if active != expectedActiveIdx[i] {
			t.Errorf("activated index was %d but should be %d for %s", active,
				expectedActiveIdx[i], name)
		}
	}
}

func TestSplitActiveInactive(t *testing.T) {
	const (
		act   = 100.0
		inact = 150.0
	)

	aTsk := makeTask("active")
	aTsk.cumulative = act

	tmp1, tmp2 := splitActiveInactiveTime(aTsk)
	if tmp1 != act {
		t.Errorf("active value was %f but should be %f", tmp1, act)
	}

	if tmp2 != 0.0 {
		t.Errorf("inavtice value was %f but should be %f", tmp2, 0.0)
	}

	inTsk := makeTask(inactiveTask)
	inTsk.cumulative = inact

	tmp1, tmp2 = splitActiveInactiveTime(inTsk)
	if tmp1 != 0.0 {
		t.Errorf("active value was %f but should be %f", tmp1, 0.0)
	}

	if tmp2 != inact {
		t.Errorf("inavtice value was %f but should be %f", tmp2, inact)
	}
}
