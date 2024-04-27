package todo

import (
	"encoding/json"
	"fmt"
	"strconv"
	"time"

	"github.com/alexeyco/simpletable"
	db "github.com/nickest14/go-tdlist/pkg/db"
)

const (
	keyPrefix  = "task_"
	timeFormat = "2006-01-02 15:04 Mon"
)

type item struct {
	Id          int64
	Task        string
	Done        bool
	CreatedAt   time.Time
	CompletedAt time.Time
}

func parseItem(data string) (item, error) {
	var i item
	err := json.Unmarshal([]byte(data), &i)
	if err != nil {
		return i, err
	}
	return i, nil
}

func AddTask(kv *db.KV, task string) {
	now := time.Now().Unix()
	todo := item{
		Id:          now,
		Task:        task,
		Done:        false,
		CreatedAt:   time.Now(),
		CompletedAt: time.Time{},
	}
	value, err := json.Marshal(todo)
	if err != nil {
		panic(err)
	}
	key := keyPrefix + strconv.FormatInt(now, 10)
	kv.Add(key, value)
}

func GetTodos(kv *db.KV, startTime, endTime time.Time) error {
	todoStrings, err := kv.GetDescendRange(
		"CreatedAt",
		fmt.Sprintf(`{"CreatedAt": "%s"}`, endTime.Format(time.RFC3339Nano)),
		fmt.Sprintf(`{"CreatedAt": "%s"}`, startTime.Format(time.RFC3339Nano)),
	)
	if err != nil {
		return err
	}
	var todos Todos
	for _, taskString := range todoStrings {
		todo, err := parseItem(taskString)
		if err != nil {
			return err
		}
		todos = append(todos, todo)
	}
	todos.Print()
	return nil
}

func ToggleCompleteTodo(kv *db.KV, id int64) error {
	key := keyPrefix + strconv.FormatInt(id, 10)
	value, err := kv.Get(key)
	if err != nil {
		return err
	}
	var todo item
	err = json.Unmarshal([]byte(value), &todo)
	if err != nil {
		return err
	}
	if todo.Done {
		todo.Done = false
		todo.CompletedAt = time.Time{}
	} else {
		todo.Done = true
		todo.CompletedAt = time.Now()
	}
	byteValue, err := json.Marshal(todo)
	if err != nil {
		return err
	}
	kv.Add(key, byteValue)
	return nil
}

func Delete(kv *db.KV, id int64) error {
	key := keyPrefix + strconv.FormatInt(id, 10)
	err := kv.Delete(key)
	return err
}

type Todos []item

func (t *Todos) Print() {
	table := simpletable.New()

	table.Header = &simpletable.Header{
		Cells: []*simpletable.Cell{
			{Align: simpletable.AlignCenter, Text: "Id"},
			{Align: simpletable.AlignCenter, Text: "Task"},
			{Align: simpletable.AlignCenter, Text: "Done?"},
			{Align: simpletable.AlignRight, Text: "CreatedAt"},
			{Align: simpletable.AlignRight, Text: "CompletedAt"},
		},
	}

	var cells [][]*simpletable.Cell
	for _, item := range *t {
		task := blue(item.Task)
		done := blue("no")
		if item.Done {
			task = green(item.Task)
			done = green("yes")
		}
		var completedAt string
		if item.CompletedAt.IsZero() {
			completedAt = "-"
		} else {
			completedAt = item.CompletedAt.Format(timeFormat)
		}
		item.CompletedAt.Format(timeFormat)
		cells = append(cells, *&[]*simpletable.Cell{
			{Text: fmt.Sprintf("%d", item.Id)},
			{Text: task},
			{Text: done},
			{Text: item.CreatedAt.Format(timeFormat), Align: simpletable.AlignCenter},
			{Text: completedAt, Align: simpletable.AlignCenter},
		})
	}

	table.Body = &simpletable.Body{Cells: cells}

	table.Footer = &simpletable.Footer{Cells: []*simpletable.Cell{
		{Align: simpletable.AlignCenter, Span: 5, Text: red(fmt.Sprintf("You have %d pending todos", t.CountPending()))},
	}}

	table.SetStyle(simpletable.StyleUnicode)

	table.Println()
}

func (t *Todos) CountPending() int {
	total := 0
	for _, item := range *t {
		if !item.Done {
			total++
		}
	}

	return total
}
