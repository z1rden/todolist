package main

import (
	"errors"
	"github.com/spf13/cobra"
	"log"
	"os"
	"path/filepath"
	"todolist/internal/controller"
	"todolist/internal/storage"
)

type Controller interface {
	Create(name, description, deadline string) error
	Read() error
	Change(name, parameter, value string) error
	Delete(name string) error
}

func main() {
	var rootCmd cobra.Command
	var taskName string
	var taskDescription string
	var taskDeadline string
	var taskParameter string
	var parameterValue string
	var cntrllr Controller

	rootCmd = cobra.Command{
		Use:     "root",
		Version: "0.0.1",
		// Именно PersistentPreRunE, потому что функция Run явно не будет вызываться
		PersistentPreRunE: func(cmd *cobra.Command, args []string) error {
			const name = "db/storage.json"

			file, err := os.OpenFile(name, os.O_RDWR, 0666)
			// Ошибки, возвращаемые файловой системой, могут быть сравнены с ошибками типа os.ErrNotExist
			// с помощью errors.Is.
			if errors.Is(err, os.ErrNotExist) {
				// Если поставить os.ModeDir, то ее редактировать сможет только администратор.
				err := os.MkdirAll(filepath.Dir(name), os.ModePerm)
				if err != nil {
					return err
				}

				file, err = os.Create(name)
				if err != nil {
					return err
				}
			} else if err != nil {
				return err
			}

			cntrllr = controller.NewController(storage.NewStorage(file))
			return nil
		},
		Run: func(cmd *cobra.Command, args []string) {
			log.Println("start")
		},
	}

	rootCmd.PersistentFlags().StringVarP(&taskName, "name", "n", "",
		"Название задачи")
	rootCmd.PersistentFlags().StringVarP(&taskDescription, "description", "d", "",
		"Описание задачи")
	rootCmd.PersistentFlags().StringVar(&taskDeadline, "deadline", "",
		"Время окончания задачи. Формат ввода: 2006-01-02")
	rootCmd.PersistentFlags().StringVar(&taskParameter, "parameter", "",
		"Параметр задачи, который Вы хотите изменить")
	rootCmd.PersistentFlags().StringVar(&parameterValue, "value", "",
		"Значение параметра задачи, который Вы хотите установить")

	rootCmd.AddCommand(&cobra.Command{
		Use:   "create",
		Short: "",
		Long:  "",
		RunE: func(cmd *cobra.Command, args []string) error {
			log.Println("create")

			return cntrllr.Create(taskName, taskDescription, taskDeadline)
		},
	})

	rootCmd.AddCommand(&cobra.Command{
		Use:   "read",
		Short: "",
		Long:  "",
		RunE: func(cmd *cobra.Command, args []string) error {
			log.Println("read")

			return cntrllr.Read()
		},
	})

	rootCmd.AddCommand(&cobra.Command{
		Use:   "change",
		Short: "",
		Long:  "",
		RunE: func(cmd *cobra.Command, args []string) error {
			log.Println("change")

			return cntrllr.Change(taskName, taskParameter, parameterValue)
		},
	})

	rootCmd.AddCommand(&cobra.Command{
		Use:   "delete",
		Short: "",
		Long:  "",
		RunE: func(cmd *cobra.Command, args []string) error {
			log.Println("delete")

			return cntrllr.Delete(taskName)
		},
	})

	err := rootCmd.Execute()
	if err != nil {
		log.Fatalln(err)
	}
	// В терминал выводит, когда эти значения записались в переменные.
	// К примеру, если запустить программу go run cmd/todolist/main.go -n "test" -d "description"
	// --deadline "2025-01-05", то в log запишутся
	// Для проверки дебаггером нужно его настроить:
	// 1. Нажать вертикальное многоточие в меню запуска (справа вверху)
	// 2. Edit...
	// 3. Изменить program arguments (пример: -n "test" -d "test" --deadline "2025-02-18")
	log.Println(taskName, taskDescription, taskDeadline)
}
