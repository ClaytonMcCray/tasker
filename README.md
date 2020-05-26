# tasker
A change-driven CLI task manager.

## Usage
On launch, `tasker` will present you with a special task and a prompt:

```
*0. inactive
Active task was started on ...
...
>>
```

This task will have it's cumulative hours seperated from the rest
of the user's tasks. 

The asterisk indicates the currently activated task. At the prompt,
a user can enter

- The index of an existing task. This will switch that task to active
	and refresh all times.
- The name of a new task. This will create a new task and refresh times.
- Newline. This will refresh the times, leaving the currently activated task
	active. It is equivalent to entering the index of the active task.
