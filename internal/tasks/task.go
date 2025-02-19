package tasks

import (
	"errors"
	"time"
)

type TaskList []*Task
type Task struct {
	Name        string
	Description string
	Deadline    time.Time
	Status      Status
}

func (t *Task) Change(parameter, value string) error {
	switch parameter {
	case "name":
		t.Name = value
	case "description":
		t.Description = value
	case "deadline":
		valueTime, err := time.Parse("2006-01-02", value)
		if err != nil {
			return err
		}
		t.Deadline = valueTime
	case "status":
		switch value {
		case "TO DO":
			t.Status = ToDo
		case "IN PROGRESS":
			t.Status = InProgress
		case "DONE":
			t.Status = Done
		}
	default:
		return errors.New("Такого параметра не существует")
	}

	return nil
}
