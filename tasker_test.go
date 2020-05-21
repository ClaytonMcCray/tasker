package main

import "testing"

func TestMakeTask(t *testing.T) {
	tsk := makeTask("name")

	if !tsk.active() {
		t.Errorf("task is inactive when it should be active")
	}

	if tsk.name != "name" {
		t.Errorf("task.name is %s but should be %s.", tsk.name, "name")
	}
}
