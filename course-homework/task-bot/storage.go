package main

import (
	"errors"
	"sort"
	"sync"
)

var NotFound = errors.New("not found")
var NotYourTask = errors.New("not your task")

type User struct {
	ID       int64
	Username string
}

type Task struct {
	ID       int
	Title    string
	Owner    User
	Assignee *User
}

type TaskManager struct {
	tasks  map[int]*Task
	mu     sync.RWMutex
	nextID int
}

func NewTaskManager() *TaskManager {
	return &TaskManager{
		mu:    sync.RWMutex{},
		tasks: make(map[int]*Task),
	}
}

func (tm *TaskManager) Create(owner User, title string) *Task {
	tm.mu.Lock()
	defer tm.mu.Unlock()
	tm.nextID++
	id := tm.nextID
	task := &Task{ID: id, Title: title, Owner: owner}
	tm.tasks[id] = task
	return task
}

func (tm *TaskManager) GetByID(id int) (*Task, error) {
	tm.mu.RLock()
	defer tm.mu.RUnlock()
	task, ok := tm.tasks[id]
	if !ok {
		return nil, NotFound
	}
	return task, nil
}

func (tm *TaskManager) ListAll() []*Task {
	tm.mu.RLock()
	defer tm.mu.RUnlock()
	tasksList := make([]*Task, 0, len(tm.tasks))
	for _, task := range tm.tasks {
		tasksList = append(tasksList, task)
	}
	sort.Slice(tasksList, func(i, j int) bool {
		return tasksList[i].ID < tasksList[j].ID
	})
	return tasksList
}

func (tm *TaskManager) Assign(id int, user User) (task *Task, prev *User, err error) {
	tm.mu.Lock()
	defer tm.mu.Unlock()

	task, ok := tm.tasks[id]
	if !ok {
		return nil, nil, NotFound
	}

	prev = task.Assignee
	task.Assignee = &user

	return task, prev, nil
}

func (tm *TaskManager) Unassign(id int, user User) (task *Task, prev *User, err error) {
	tm.mu.Lock()
	defer tm.mu.Unlock()
	task, ok := tm.tasks[id]
	if !ok {
		return nil, nil, NotFound
	}
	if task.Assignee == nil {
		return nil, nil, NotYourTask
	}
	if task.Assignee.ID != user.ID {
		return nil, nil, NotYourTask
	}
	prev = task.Assignee
	task.Assignee = nil
	return task, prev, nil
}

func (tm *TaskManager) Resolve(id int, user User) (task *Task, prev *User, err error) {
	tm.mu.Lock()
	defer tm.mu.Unlock()
	task, ok := tm.tasks[id]
	if !ok {
		return nil, nil, NotFound
	}
	if task.Assignee == nil {
		return nil, nil, NotYourTask
	}
	if task.Assignee.ID != user.ID {
		return nil, nil, NotYourTask
	}
	prev = task.Assignee
	delete(tm.tasks, id)
	return task, prev, nil
}

func (tm *TaskManager) ListByAssignee(userID int64) []*Task {
	tm.mu.RLock()
	defer tm.mu.RUnlock()
	tasksList := make([]*Task, 0, len(tm.tasks))
	for _, task := range tm.tasks {
		if task.Assignee != nil && task.Assignee.ID == userID {
			tasksList = append(tasksList, task)
		}
	}
	sort.Slice(tasksList, func(i, j int) bool {
		return tasksList[i].ID < tasksList[j].ID
	})
	return tasksList
}

func (tm *TaskManager) ListByOwner(userID int64) []*Task {
	tm.mu.RLock()
	defer tm.mu.RUnlock()
	tasksList := make([]*Task, 0, len(tm.tasks))
	for _, task := range tm.tasks {
		if task.Owner.ID == userID {
			tasksList = append(tasksList, task)
		}
	}
	sort.Slice(tasksList, func(i, j int) bool {
		return tasksList[i].ID < tasksList[j].ID
	})
	return tasksList
}
