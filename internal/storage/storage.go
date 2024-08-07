package storage

import (
	"fmt"
	"sync"
	"tasks/internal/entities"

	gonanoid "github.com/matoous/go-nanoid/v2"
)

type Storage struct {
	m     sync.Mutex
	Tasks map[string]entities.Task 
	Users map[string]entities.User
}

func NewStorage() *Storage {
	return &Storage{
		Tasks: make(map[string]entities.Task),
		Users: make(map[string]entities.User),
	}
}

func (s *Storage) GetTasks() []entities.Task {
	tasks := make([]entities.Task, 0, len(s.Tasks))

	for _, task := range s.Tasks {
		tasks = append(tasks, task)
	}

	return tasks
}

func (s *Storage) CreateTask(task entities.Task) (string, bool) {
	s.m.Lock()

	defer s.m.Unlock()

	id, err := gonanoid.New()

	if err != nil {
		return "", false
	}

	task.ID = id
	s.Tasks[task.ID] = task

	fmt.Print(task)
	return id, true
}

func (s *Storage) UpdadeTask(id string, task entities.Task) bool {
	s.m.Lock()

	defer s.m.Unlock()

	t, ok := s.Tasks[id]

	if !ok {
		return false
	}

	task.ID = t.ID
	s.Tasks[task.ID] = task

	return true
}

func (s *Storage) DeleteTask(id string) bool {
	s.m.Lock()

	defer s.m.Unlock()

	_, ok := s.Tasks[id]

	if !ok {
		return false
	}

	delete(s.Tasks, id)
	return true
}

func (s *Storage) GetUser(username string) (entities.User, bool) {
	fmt.Print(username)
	s.m.Lock()

	defer s.m.Unlock()

	user, ok := s.Users[username]

	fmt.Print(user)
	return user, ok
}

func (s *Storage) CreateUser(u entities.User) bool {
	s.m.Lock()

	defer s.m.Unlock()

	_, ok := s.Users[u.Username]

	if ok {
		return false
	}

	s.Users[u.Username] = u
	return true
}
