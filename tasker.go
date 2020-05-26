package main

import (
	"bufio"
	"fmt"
	"os"
	"os/exec"
	"runtime"
	"strconv"
	"time"
)

const (
	newTask = iota
	taskIndex
	refresh
)

const (
	inactiveTask = "inactive"
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

func splitActiveInactiveTime(t *task) (active, inactive float64) {
	if t.name == inactiveTask {
		return 0.0, t.cumulative
	} else {
		return t.cumulative, 0.0
	}
}

func printTasks(tasks []*task) {
	var activeT *task
	var toPrint string
	var cumulativeActive float64
	var cumulativeInactive float64

	for idx, currTask := range tasks {
		if currTask.active() {
			activeT = currTask
			toPrint += "*"
		} else {
			toPrint += " "
		}

		toPrint += fmt.Sprintf("%d. %.2f %s\n", idx, currTask.cumulative, currTask.name)
		tmpAct, tmpInact := splitActiveInactiveTime(currTask)
		cumulativeActive += tmpAct
		cumulativeInactive += tmpInact
	}

	if activeT == nil {
		return
	}

	fmt.Print(toPrint)
	fmt.Printf("Active task was started on %s.\n", activeT.lastStart.Format(time.UnixDate))
	fmt.Printf("Hours active: %.2f \tHours inactive: %.2f.\n", cumulativeActive, cumulativeInactive)
}

func clearTerminal() {
	cmd := exec.Command("")
	switch runtime.GOOS {
	case "windows":
		cmd = exec.Command("cmd", "/c", "cls")
	default:
		cmd = exec.Command("clear")
	}

	cmd.Stdout = os.Stdout
	cmd.Run()
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

	tasks = append(tasks, makeTask(inactiveTask))

	for {
		clearTerminal()
		printTasks(tasks)
		fmt.Print(">>  ")
		taskIdxToStamp, tasks = inputHandler(scanner, tasks)
		stampActiveTask(tasks)
		if len(tasks) > 0 {
			tasks[taskIdxToStamp].stamp()
		}
	}
}
