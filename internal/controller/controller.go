package controller

import (
	"fmt"
	"time"
	"todolist/internal/tasks"
)

type Storage interface {
	Create(task tasks.Task) error
	Read() (tasks.TaskList, error)
	Change(name, parameter, value string) error
	Delete(name string) error
}

type Controller struct {
	storage Storage
}

func NewController(storage Storage) *Controller {
	return &Controller{storage: storage}
}

func (c *Controller) Create(name, description, deadline string) error {
	t, err := time.Parse("2006-01-02", deadline)
	if err != nil {
		return err
	}

	return c.storage.Create(tasks.Task{
		Name:        name,
		Description: description,
		Deadline:    t,
	})
}

func (c *Controller) Read() error {
	list, err := c.storage.Read()
	if err != nil {
		return err
	}

	for i := range list {
		fmt.Printf("status: %s, name: %s, description: %s, deadline: %s\n",
			list[i].Status, list[i].Name, list[i].Description, list[i].Deadline)
	}

	return nil
}

func (c *Controller) Change(name, parameter, value string) error {
	if err := c.storage.Change(name, parameter, value); err != nil {
		return err
	}

	return nil
}

func (c *Controller) Delete(name string) error {
	if err := c.storage.Delete(name); err != nil {
		return err
	}

	return nil
}
