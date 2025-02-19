package storage

import (
	"encoding/json"
	"errors"
	"io"
	"todolist/internal/tasks"
)

type Task interface {
	Change(parameter, value string) error
}

type File interface {
	io.ReadWriteCloser
	Truncate(size int64) error
	Seek(offset int64, whence int) (int64, error)
}

type Storage struct {
	file File
}

func NewStorage(file File) *Storage {
	return &Storage{file: file}
}

func (s *Storage) Create(task tasks.Task) error {
	list, err := s.Read()
	if err != nil {
		return err
	}

	for i := range list {
		if list[i].Name == task.Name {
			return errors.New("Такая задача уже присутствует")
		}
	}

	if err := s.writeTasks(append(list, &task)); err != nil {
		return err
	}

	return nil
}

func (s *Storage) Read() (tasks.TaskList, error) {
	readAll, err := io.ReadAll(s.file)
	if err != nil {
		return nil, err
	}

	var list tasks.TaskList
	if len(readAll) == 0 {
		return list, nil
	}
	if err := json.Unmarshal(readAll, &list); err != nil {
		return nil, err
	}

	return list, nil
}

func (s *Storage) Change(name, parameter, value string) error {
	var flagUpdate bool

	allTasks, err := s.Read()
	if err != nil {
		return err
	}
	for i := range allTasks {
		if allTasks[i].Name == name {
			flagUpdate = true
			if err := allTasks[i].Change(parameter, value); err != nil {
				return err
			}
			break
		}
	}
	if !flagUpdate {
		return errors.New("Такого параметра не существует.")
	}

	if err := s.writeTasks(allTasks); err != nil {
		return err
	}

	return nil
}

func (s *Storage) writeTasks(allTasks tasks.TaskList) error {
	bytes, err := json.Marshal(allTasks)
	if err != nil {
		return err
	}
	if err := s.file.Truncate(0); err != nil {
		return err
	}
	if _, err := s.file.Seek(0, 0); err != nil {
		return err
	}
	if _, err := s.file.Write(bytes); err != nil {
		return err
	}

	return nil
}

func (s *Storage) Delete(name string) error {
	allTasks, err := s.Read()
	if err != nil {
		return err
	}
	for i := range allTasks {
		if allTasks[i].Name == name {
			allTasks = append(allTasks[:i], allTasks[i+1:]...)
			if err := s.writeTasks(allTasks); err != nil {
				return err
			}
			return nil
		}
	}
	return errors.New("Такой задачи не существует.")
}
