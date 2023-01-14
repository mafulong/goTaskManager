package core

import (
	"context"
	"testing"
)

func TestManager_executeTask(t *testing.T) {
	Init()
	ctx := context.Background()
	mgr := NewManager(ctx)
	mgr.executeTask("id4")
}
