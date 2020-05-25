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
