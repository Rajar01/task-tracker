package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"math/rand"
	"os"
	"strconv"
	"strings"
	"time"
)

type Status int

const (
	TODO Status = iota
	INPROGRESS
	DONE
)

func (s Status) String() string {
	switch s {
	case TODO:
		return "todo"
	case INPROGRESS:
		return "in-progress"
	case DONE:
		return "done"
	default:
		return "unknown"
	}
}

type Task struct {
	Id          uint32    `json:"id"`
	Description string    `json:"description"`
	Status      Status    `json:"status"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

func AddTask(description string) {
	randomizer := rand.New(rand.NewSource(10))
	var task Task = Task{Id: randomizer.Uint32(), Description: description, Status: TODO, CreatedAt: time.Now(), UpdatedAt: time.Now()}
	tasks = append(tasks, task)
	WriteTasksToJson(tasks)
}

func UpdateTask(id uint32, description string) {
	for i := range tasks {
		if tasks[i].Id == id {
			tasks[i].Description = description
			tasks[i].UpdatedAt = time.Now()
			break
		}
	}

	WriteTasksToJson(tasks)
}

func DeleteTask(id uint32) {
	var filteredTasks []Task
	for _, task := range tasks {
		if task.Id != id {
			filteredTasks = append(filteredTasks, task)
		}
	}

	WriteTasksToJson(filteredTasks)
}

func MarkTask(id uint32, status Status) {
	for i := range tasks {
		if tasks[i].Id == id {
			tasks[i].Status = status
			tasks[i].UpdatedAt = time.Now()
			break
		}
	}

	WriteTasksToJson(tasks)
}

func GetTasks() {
	for _, task := range tasks {
		fmt.Printf("%+v\n", task)
	}
}

func GetTasksByStatus(status Status) {
	for _, task := range tasks {
		if task.Status == status {
			fmt.Printf("%+v\n", task)
		}
	}
}

func WriteTasksToJson(tasks []Task) {
	tasksJson, err := json.Marshal(tasks)
	if err != nil {
		log.Fatalf("Failed to convert task object to json string\n\n%v", err)
	}

	err = os.WriteFile(FILENAME, tasksJson, 0644)
	if err != nil {
		log.Fatalf("Failed to write task to %s file\n\n%v", FILENAME, err)
	}
}

func StringToStatus(status string) (Status, error) {
	status = strings.ToLower(status)

	switch status {
	case "todo":
		return TODO, nil
	case "in-progress":
		return INPROGRESS, nil
	case "done":
		return DONE, nil
	default:
		return -1, errors.New("Unknown task status")
	}
}

const FILENAME = "tasks.json"

var tasks []Task

func init() {
	if _, err := os.Stat(FILENAME); os.IsNotExist(err) {
		file, err := os.Create(FILENAME)
		if err != nil {
			log.Fatalf("Failed to create a %s file\n\n%v", FILENAME, err)
		}

		defer file.Close()
	} else {
		content, err := os.ReadFile(FILENAME)
		if err != nil {
			log.Fatalf("Failed to read a %s file\n\n%v", FILENAME, err)
		}

		err = json.Unmarshal(content, &tasks)
		if err != nil {
			log.Fatalf("Failed to turn all tasks.json content to array of task object\n\n%v", err)
		}
	}
}

func main() {
	switch os.Args[1] {
	case "add":
		if len(os.Args) < 3 {
			log.Fatal("Usage: task-cli add <description1> <description2> ...")
		}

		for _, description := range os.Args[2:] {
			AddTask(description)
		}
	case "update":
		if len(os.Args) != 4 {
			log.Fatal("Usage: task-cli update <id> <description>")
		}

		taskId, err := strconv.ParseUint(os.Args[2], 10, 32)
		if err != nil {
			log.Fatal("Usage: task-cli update <id> <description>")
		}

		UpdateTask(uint32(taskId), os.Args[3])
	case "delete":
		if len(os.Args) != 3 {
			log.Fatal("Usage: task-cli delete <id>")
		}

		taskId, err := strconv.ParseUint(os.Args[2], 10, 32)
		if err != nil {
			log.Fatal("Usage: task-cli delete <id>")
		}

		DeleteTask(uint32(taskId))
	case "list":
		if len(os.Args) > 3 {
			log.Fatal("Usage: task-cli list <status>")
		} else if len(os.Args) == 3 {
			status, err := StringToStatus(os.Args[2])
			if err != nil {
				log.Fatal(err)
			}

			GetTasksByStatus(status)
		} else {
			GetTasks()
		}
	default:
		log.Fatal("Unknown command")
	}
}
