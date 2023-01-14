package core

import (
	"context"
	"sync"
	"sync/atomic"
)

type TaskStatus struct {
	Done   chan struct{}
	Status atomic.Int32
}

const (
	StatusNotStart   = 0
	StatusHasStarted = 1
	StatusDone       = 2
)

type Manager struct {
	Ctx        context.Context
	DomainData *DomainData
	taskInfo   map[string]*TaskStatus
}

func NewManager(ctx context.Context) *Manager {
	res := &Manager{
		Ctx:        ctx,
		DomainData: &DomainData{},
		taskInfo:   map[string]*TaskStatus{},
	}
	for _, task := range TaskMap {
		status := &TaskStatus{
			Done:   make(chan struct{}),
			Status: atomic.Int32{},
		}
		status.Status.Store(StatusNotStart)
		res.taskInfo[task.ID] = status
	}
	return res
}

func (mgr *Manager) executeTask(taskID string) {
	task := TaskMap[taskID]
	wg := &sync.WaitGroup{}
	waitDeps := []string{}
	defer func() {
		close(mgr.taskInfo[taskID].Done)
		mgr.taskInfo[taskID].Status.Store(StatusDone)
	}()
	for i, dep := range DependencyMap[taskID] {
		if mgr.taskInfo[dep].Status.Load() == StatusDone {
			continue
		}
		if mgr.taskInfo[dep].Status.CompareAndSwap(StatusNotStart, StatusHasStarted) {
			if i == len(DependencyMap[taskID])-1 {
				// optimization: reduce one goroutine
				mgr.executeTask(dep)
			} else {
				wg.Add(1)
				go func(id string) {
					defer wg.Done()
					mgr.executeTask(id)
				}(dep)
			}
		} else {
			waitDeps = append(waitDeps, dep)
		}
	}
	for _, dep := range waitDeps {
		select {
		case <-mgr.taskInfo[dep].Done:
		}
	}
	wg.Wait()
	task.Handle(mgr.DomainData)
}
