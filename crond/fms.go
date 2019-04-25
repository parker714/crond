package crond

import (
	"cron-s/tasks"
	"encoding/json"
	"github.com/hashicorp/raft"
	"io"
)

type Fms struct {
	ctx *Context
}

func (f *Fms) Apply(l *raft.Log) interface{} {
	f.ctx.Crond.Log.Println("[DEBUG] fms: Apply")

	t := tasks.NewTask()
	if err := json.Unmarshal(l.Data, t); err != nil {
		f.ctx.Crond.Log.Println("[WARN] fms: Apply Unmarshal err", err)
		return nil
	}

	switch t.Status {
	case tasks.StatusAdd:
		tasks.Add(t)
	case tasks.StatusDel:
		tasks.Del(t)
	}

	return nil
}

func (f *Fms) Snapshot() (raft.FSMSnapshot, error) {
	f.ctx.Crond.Log.Println("[DEBUG] fms: Snapshot")

	return &FmsSnapshot{
		ctx: &Context{Crond: f.ctx.Crond},
	}, nil
}

func (f *Fms) Restore(serialized io.ReadCloser) error {
	f.ctx.Crond.Log.Println("[DEBUG] fpm: Restore")

	nh := tasks.NewHeap()
	if err := json.NewDecoder(serialized).Decode(nh); err != nil {
		return err
	}
	tasks.Init(nh)

	return nil
}
