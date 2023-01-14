package core

import "fmt"

type TaskType int

const (
	TaskTypeLoad     = 1
	TaskTypeAssemble = 2
)

type Task struct {
	ID          string
	Description string
	Type        TaskType
	Handle      func(mgr *DomainData) error
}

var (
	TaskMap       = make(map[string]*Task)
	DependencyMap = make(map[string][]string)
)

func RegisterTask(ID, description string, taskType TaskType, handlerFunc func(data *DomainData) error) {
	task := &Task{
		ID:          ID,
		Description: description,
		Type:        taskType,
		Handle:      handlerFunc,
	}
	TaskMap[task.ID] = task
}

func RegisterDependencies(ID string, dependencies []string) {
	DependencyMap[ID] = append(DependencyMap[ID], dependencies...)
}

func Init() {
	RegisterTask("id1", "", TaskTypeLoad, func(data *DomainData) error {
		fmt.Println("id1")
		return nil
	})
	RegisterTask("id2", "", TaskTypeLoad, func(data *DomainData) error {
		fmt.Println("id2")
		return nil
	})
	id := "id3"
	RegisterTask(id, id, TaskTypeLoad, func(data *DomainData) error {
		fmt.Println(id)
		return nil
	})
	id1 := "id4"
	RegisterTask(id1, id1, TaskTypeLoad, func(data *DomainData) error {
		fmt.Println(id1)
		return nil
	})
	RegisterDependencies("id2", []string{"id1"})
	RegisterDependencies("id3", []string{"id1"})
	RegisterDependencies("id4", []string{"id2", "id3"})
}
