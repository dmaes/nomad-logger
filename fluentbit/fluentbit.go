package fluentbit

import (
	"fmt"
	"log"
	"time"

	"github.com/dmaes/nomad-logger/nomad"
	"github.com/dmaes/nomad-logger/util"

	"github.com/hashicorp/nomad/api"
)

type Fluentbit struct {
	Nomad     *nomad.Nomad
	ConfFile  string
	TagPrefix string
	ReloadCmd string
}

func (f *Fluentbit) Run() {
	log.Println("Starting nomad-logger for Fluentbit")
	for {
		time.Sleep(1 * time.Second)
		allocs := f.Nomad.Allocs()
		f.UpdateConf(allocs)
	}
}

func (f *Fluentbit) UpdateConf(Allocs []*api.Allocation) {
	config := ""

	for _, alloc := range Allocs {
		config += f.AllocToConfig(alloc)
	}

	util.WriteConfig(config, f.ConfFile, f.ReloadCmd)
}

func (f *Fluentbit) AllocToConfig(Alloc *api.Allocation) string {
	tasks, err := nomad.AllocTasks(Alloc)
	if err != nil {
		log.Fatalln(err)
	}

	config := ""

	for _, task := range tasks {
		config += f.AllocTaskStreamToConfig(Alloc, task, "stdout")
		config += f.AllocTaskStreamToConfig(Alloc, task, "stderr")
	}

	return config
}

func (f *Fluentbit) AllocTaskStreamToConfig(Alloc *api.Allocation, Task *api.Task, Stream string) string {
	tag := fmt.Sprintf("%s.%s.%s.%s", f.TagPrefix, Alloc.ID, Task.Name, Stream)
	path := fmt.Sprintf("%s/%s/alloc/logs/%s.%s.[0-9]*", f.Nomad.AllocsDir, Alloc.ID, Task.Name, Stream)

	config := fmt.Sprintf(`
[INPUT]
  name tail
  tag %s
  path %s
[FILTER]
  name modify
  match %s
  add nomad_namespace %s
  add nomad_job %s
  add nomad_task_group %s
  add nomad_task %s
  add nomad_alloc_id %s
  add nomad_alloc_name %s
  add nomad_node_id %s
  add nomad_log_stream %s
  `, tag, path, tag, Alloc.Namespace, Alloc.JobID, Alloc.TaskGroup, Task.Name, Alloc.ID, Alloc.Name, f.Nomad.NodeID, Stream)

	return config
}
