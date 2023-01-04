package fluentbit

import (
	"bytes"
	"fmt"
	"log"
	"strings"
	"text/template"
	"time"

	"github.com/dmaes/nomad-logger/nomad"
	"github.com/dmaes/nomad-logger/util"

	"github.com/hashicorp/nomad/api"

	_ "embed"
)

type Fluentbit struct {
	Nomad     *nomad.Nomad
	ConfFile  string
	TagPrefix string
	Parser    string
	ReloadCmd string
}

//go:embed fluentbit-conf.gotmpl
var FluentbitConfTmpl string

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
	tagPrefix := f.Nomad.TaskMetaGet(*Task, "fluentbit.tag-prefix", f.TagPrefix)
	tag := fmt.Sprintf("%s.%s.%s.%s", tagPrefix, Alloc.ID, Task.Name, Stream)

	path := fmt.Sprintf("%s/%s/alloc/logs/%s.%s.[0-9]*", f.Nomad.AllocsDir, Alloc.ID, Task.Name, Stream)

	parser := f.Nomad.TaskMetaGet(*Task, "fluentbit.parser", f.Parser)

	filterParsersStr := f.Nomad.TaskMetaGet(*Task, "fluentbit.filter-parsers", "")
	filterParsers := make([]*FluentbitFilterParser, 0)
	if filterParsersStr != "" {
		for _, fp := range strings.Split(filterParsersStr, ",") {
			if fp != "" {
				fpSplit := strings.Split(fp, ":")
				filterParsers = append(filterParsers, &FluentbitFilterParser{
					Key:    fpSplit[0],
					Parser: fpSplit[1],
				})
			}
		}
	}

	fluentbitConfig := &FluentbitConfig{
		Tag:            tag,
		Path:           path,
		Parser:         parser,
		FilterParsers:  filterParsers,
		NomadNamespace: Alloc.Namespace,
		NomadJob:       Alloc.JobID,
		NomadTaskGroup: Alloc.TaskGroup,
		NomadTask:      Task.Name,
		NomadAllocID:   Alloc.ID,
		NomadAllocName: Alloc.Name,
		NomadNodeID:    f.Nomad.NodeID,
		NomadLogStream: Stream,
	}

	tpl := template.Must(template.New("fluentbit-conf").Parse(FluentbitConfTmpl))
	var tplBuffer bytes.Buffer
	err := tpl.Execute(&tplBuffer, fluentbitConfig)
	if err != nil {
		panic(err)
	}

	return tplBuffer.String()
}

type FluentbitConfig struct {
	Tag            string
	Path           string
	Parser         string
	FilterParsers  []*FluentbitFilterParser
	NomadNamespace string
	NomadJob       string
	NomadTaskGroup string
	NomadTask      string
	NomadAllocID   string
	NomadAllocName string
	NomadNodeID    string
	NomadLogStream string
}

type FluentbitFilterParser struct {
	Key    string
	Parser string
}
