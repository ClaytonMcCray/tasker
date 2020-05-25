package main

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"time"
)

const (
	newTask = iota
	taskIndex
	refresh
)

type task struct {
	name       string
	cumulative float64
	lastStart  time.Time
}

func (t *task) stamp() float64 {
	if t.lastStart == (time.Time{}) {
		t.lastStart = time.Now()
	}

	t.cumulative += time.Since(t.lastStart).Hours()
	t.lastStart = time.Now()
	return t.cumulative
}

func (t *task) deactivate() {
	t.lastStart = time.Time{}
}

func (t *task) active() bool {
	return t.lastStart != (time.Time{})
}

func makeTask(name string) *task {
	return &task{
		name:       name,
		cumulative: 0.0,
		lastStart:  time.Now(),
	}
}

func printTasks(tasks []*task) {
	var activeT *task
	var toPrint string
	cumulativeTotalTime := 0.0
	for idx, currTask := range tasks {
		if currTask.active() {
			activeT = currTask
			toPrint += "*"
		} else {
			toPrint += " "
		}

		toPrint += fmt.Sprintf("%d. %.2f %s\n", idx, currTask.cumulative, currTask.name)
		cumulativeTotalTime += currTask.cumulative
	}

	if activeT == nil {
		return
	}

	fmt.Print(toPrint)
	fmt.Printf("Active task was started on %s.\n", activeT.lastStart.Format(time.UnixDate))
	fmt.Printf("Hours worked: %.2f.\n", cumulativeTotalTime)
}

func deactivateAll(tasks []*task) int {
	activeIdx := 0
	for idx, t := range tasks {
		if t.active() {
			t.stamp()
			activeIdx = idx
		}

		t.deactivate()
	}

	return activeIdx
}

func stampActiveTask(tasks []*task) {
	for _, t := range tasks {
		if t.active() {
			t.stamp()
		}
	}
}

type inOut interface {
	Scan() bool
	Text() string
}

// Determine what type of input the user handed to the app.
func parseInputCharacter(input string) (typeOfInput int, name string, idxActivate int) {
	if len(input) < 1 {
		return refresh, "", -1
	}

	taskIdxToStamp, err := strconv.Atoi(string(input[0]))

	if err != nil {
		return newTask, input, -1
	} else {
		return taskIndex, "", taskIdxToStamp
	}
}

// Determine what task to activate and/or add a new task.
func inputHandler(scanner inOut, tasks []*task) (int, []*task) {
	scanner.Scan()
	input := scanner.Text()
	previouslyActive := deactivateAll(tasks)

	typeOfInput, name, idx := parseInputCharacter(input)

	switch typeOfInput {
	case newTask:
		tasks = append(tasks, makeTask(name))
		return len(tasks) - 1, tasks
	case refresh:
		return previouslyActive, tasks
	default:
		return idx, tasks
	}
}

// TODO: add continuous logging and an inactive state
func main() {
	tasks := make([]*task, 0)
	var taskIdxToStamp int
	scanner := bufio.NewScanner(os.Stdin)

	for {
		printTasks(tasks)
		fmt.Print(">>  ")
		taskIdxToStamp, tasks = inputHandler(scanner, tasks)
		stampActiveTask(tasks)
		if len(tasks) > 0 {
			tasks[taskIdxToStamp].stamp()
		}
	}
}
