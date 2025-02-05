package main

import (
	"encoding/json"
	"fmt"
	"log"
	"math/rand"
	"os"
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
	tasksJson, err := json.Marshal(tasks)
	if err != nil {
		log.Fatalf("Failed to convert task object to json string\n\n%v", err)
	}

	err = os.WriteFile(FILENAME, tasksJson, 0644)
	if err != nil {
		log.Fatalf("Failed to write task to %s file\n\n%v", FILENAME, err)
	}
}

func UpdateTask(id uint, description string) {}

func DeleteTask(id uint) {}

func MarkTask(id uint, status Status) {}

func GetTasks() {}

func GetTasksByStatus(status Status) {}

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
}
