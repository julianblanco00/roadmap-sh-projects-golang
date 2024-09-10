package main

import (
	"bufio"
	"encoding/json"
	"errors"
	"fmt"
	"math/rand"
	"os"
	"strconv"
	"strings"
	"time"
)

const (
	FILE_NAME        = "tasks.json"
	DONE             = "done"
	IN_PROGRESS      = "in-progress"
	TODO             = "todo"
	ADD              = "add"
	UPDATE           = "update"
	DELETE           = "delete"
	EXIT             = "exit"
	MARK_DONE        = "mark-done"
	MARK_IN_PROGRESS = "mark-in-progress"
	LIST             = "list"
)

func createTasksJSONFile() {
	fi, err := os.Stat(FILE_NAME)
	if fi != nil {
		return
	}

	if errors.Is(err, os.ErrNotExist) {
		err = os.WriteFile(FILE_NAME, nil, 0644)
		if err != nil {
			panic(err)
		}
		return
	}

	if err != nil {
		panic(err)
	}
}

type Task struct {
	Id          int       `json:"id"`
	Description string    `json:"description"`
	Status      string    `json:"status"`
	CreatedAt   time.Time `json:"createdAt"`
	UpdatedAt   time.Time `json:"updatedAt"`
}

func newTask(description string) Task {
	return Task{
		Id:          rand.Intn(1000),
		Description: description,
		Status:      TODO,
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
	}
}

func readTasks() (map[int]Task, error) {
	bytes, err := os.ReadFile(FILE_NAME)
	if err != nil {
		return nil, err
	}

	tasks := make(map[int]Task)
	json.Unmarshal(bytes, &tasks)

	return tasks, nil
}

func writeTasks(tasks map[int]Task) error {
	bytes, err := json.Marshal(tasks)
	if err != nil {
		return err
	}

	err = os.WriteFile(FILE_NAME, bytes, 0644)
	if err != nil {
		return err
	}

	return nil
}

func addTask(task Task) {
	tasks, err := readTasks()
	if err != nil {
		fmt.Println("Could not read the tasks.json file")
		return
	}

	tasks[task.Id] = task

	err = writeTasks(tasks)

	if err != nil {
		fmt.Println("Could not write to the tasks.json file")
		return
	}

	println("\U0001F44D", "Task added successfully", task.Id)
}

func updateTask(taskId string, description string) {
	taskIdInt, err := strconv.Atoi(taskId)
	if err != nil {
		fmt.Println("Invalid task id")
		return
	}

	tasks, err := readTasks()
	if err != nil {
		fmt.Println("Could not read the tasks.json file")
		return
	}

	tasks[taskIdInt] = Task{
		Id:          taskIdInt,
		Description: description,
		Status:      tasks[taskIdInt].Status,
		CreatedAt:   tasks[taskIdInt].CreatedAt,
		UpdatedAt:   time.Now(),
	}

	err = writeTasks(tasks)
	if err != nil {
		fmt.Println("Could not write to the tasks.json file")
		return
	}

	println("\U0001F44D", "Task updated successfully")
}

func deleteTask(taskId string) {
	taskIdInt, err := strconv.Atoi(taskId)
	if err != nil {
		fmt.Println("Invalid task id")
		return
	}

	tasks, err := readTasks()
	if err != nil {
		fmt.Println("Could not read the tasks.json file")
		return
	}

	delete(tasks, taskIdInt)

	err = writeTasks(tasks)
	if err != nil {
		fmt.Println("Could not write to the tasks.json file")
		return
	}

	println("\U0001F44D", "Task deleted successfully")
}

func listAllTasks() {
	tasks, err := readTasks()
	if err != nil {
		fmt.Println("Could not read the tasks.json file")
		return
	}

	if len(tasks) == 0 {
		fmt.Println("No tasks found")
		return
	}

	for _, task := range tasks {
		fmt.Println(task.Id, task.Description, task.Status)
	}
}

func listTasksWithFilter(status string) {
	tasks, err := readTasks()
	if err != nil {
		fmt.Println("Could not read the tasks.json file")
		return
	}

	found := false
	for _, task := range tasks {
		if task.Status == status {
			found = true
			fmt.Println(task.Id, task.Description, task.Status)
		}
	}

	if !found {
		fmt.Println("No tasks found with status", status)
	}
}

func updateTaskStatus(taskId string, status string) {
	tasks, err := readTasks()
	if err != nil {
		fmt.Println("Could not read the tasks.json file")
		return
	}

	taskIdInt, err := strconv.Atoi(taskId)
	if err != nil {
		fmt.Println("Invalid task id")
		return
	}

	tasks[taskIdInt] = Task{
		Id:          taskIdInt,
		Description: tasks[taskIdInt].Description,
		Status:      status,
		CreatedAt:   tasks[taskIdInt].CreatedAt,
		UpdatedAt:   time.Now(),
	}

	err = writeTasks(tasks)
	if err != nil {
		fmt.Println("Could not write to the tasks.json file")
		return
	}

	println("\U0001F44D", "Task status updated successfully")
}

func readFromCommandLine() {
	reader := bufio.NewReader(os.Stdin)
	fmt.Println("\nWelcome to the task manager")
	fmt.Println("Commands:")
	fmt.Println("add \"<task description>\"")
	fmt.Println("update <task id> \"<new task description>\"")
	fmt.Println("delete <task id>")
	fmt.Println("list")
	fmt.Println("list <status>")
	fmt.Println("mark-in-progress <task id>")
	fmt.Println("mark-done <task id>")

	for {
		fmt.Printf("\nEnter command: ")

		command, _ := reader.ReadString('\n')
		trimedCommand := strings.TrimSpace(command)
		splitCommand := strings.Split(trimedCommand, " ")

		action := splitCommand[0]

		var id string
		if len(splitCommand) > 1 {
			id = splitCommand[1]
		}

		switch action {
		case ADD:
			description := strings.Split(trimedCommand, "\"")[1]
			task := newTask(strings.ReplaceAll(description, "\"", ""))
			addTask(task)
		case UPDATE:
			description := strings.Split(trimedCommand, "\"")[1]
			updateTask(id, strings.ReplaceAll(description, "\"", ""))
		case DELETE:
			deleteTask(id)
		case LIST:
			if id == "" {
				listAllTasks()
			} else {
				listTasksWithFilter(id)
			}
		case MARK_IN_PROGRESS:
			updateTaskStatus(id, IN_PROGRESS)
		case MARK_DONE:
			updateTaskStatus(id, DONE)
		case EXIT:
			os.Exit(0)
		default:
			fmt.Println("Invalid command")
		}
	}
}

func main() {
	createTasksJSONFile()
	readFromCommandLine()
}
