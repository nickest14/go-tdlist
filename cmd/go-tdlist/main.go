package main

import (
	"bufio"
	"errors"
	"flag"
	"fmt"
	"io"
	"os"
	"strings"

	db "github.com/nickest14/go-tdlist/pkg/db"
	todo "github.com/nickest14/go-tdlist/pkg/todo"
	utils "github.com/nickest14/go-tdlist/pkg/utils"
)

func main() {
	if err := run(); err != nil {
		fmt.Fprintln(os.Stderr, err.Error())
		os.Exit(1)
	}
}

func run() error {
	add := flag.Bool("add", false, "Add a new todo")
	complete := flag.Int64("complete", 0, "Toggle todo completion")
	del := flag.Int64("del", 0, "Delete a todo")
	list := flag.Bool("list", false, "List all todos by date range, default is today")
	startDate := flag.String("start_date", "", "Start date")
	endDate := flag.String("end_date", "", "End date")

	flag.Parse()

	kv, _ := db.NewBuntDb("todos.db")
	defer kv.Close()

	kv.CreateJsonIndex("CreatedAt")

	switch {
	case *add:
		task, err := getInput(os.Stdin, flag.Args()...)
		if err != nil {
			return err
		}
		todo.AddTask(kv, task)

	case *complete > 0:
		todo.ToggleCompleteTodo(kv, *complete)
	case *del > 0:
		todo.Delete(kv, *del)
	case *list:
		startDate, err := utils.ParseDate(startDate)
		if err != nil {
			return err
		}
		endDate, err := utils.ParseDate(endDate)
		if err != nil {
			return err
		}
		endDate = utils.EndOfDay(endDate)

		if err := todo.GetTodos(kv, startDate, endDate); err != nil {
			return nil
		}
	default:
		return errors.New("Invalid command")
	}
	return nil
}

func getInput(r io.Reader, args ...string) (string, error) {

	if len(args) > 0 {
		return strings.Join(args, " "), nil
	}
	fmt.Println("Add a todo task: ")
	scanner := bufio.NewScanner(r)
	scanner.Scan()
	if err := scanner.Err(); err != nil {
		return "", err
	}

	text := scanner.Text()

	if len(text) == 0 {
		return "", errors.New("empty todo is not allowed")
	}

	return text, nil
}
