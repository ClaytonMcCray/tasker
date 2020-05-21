package main

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"time"
)

const (
	BLANK = "\n"
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
		}

		toPrint += fmt.Sprintf("%d. %.2f %s", idx, currTask.cumulative, currTask.name)
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

// TODO: add continuous logging and an inactive state
func main() {
	tasks := make([]*task, 0)
	reader := bufio.NewReader(os.Stdin)

	for {
		printTasks(tasks)
		fmt.Print(">>  ")
		input, _ := reader.ReadString('\n')
		taskIdxToStamp, err := strconv.Atoi(string(input[0]))
		previouslyActive := deactivateAll(tasks)
		if err != nil && input != BLANK {
			tasks = append(tasks, makeTask(input))
			continue
		} else if input == BLANK {
			taskIdxToStamp = previouslyActive
		}

		stampActiveTask(tasks)
		tasks[taskIdxToStamp].stamp()
	}
}
