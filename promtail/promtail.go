package promtail

import (
	"fmt"
	"log"
	"time"

	"github.com/dmaes/nomad-logger/nomad"
	"github.com/dmaes/nomad-logger/util"

	"github.com/hashicorp/nomad/api"
	"gopkg.in/yaml.v3"
)

type Promtail struct {
	Nomad       *nomad.Nomad
	TargetsFile string
}

func (p *Promtail) Run() {
	log.Println("Starting nomad-logger for Promtail")
	for {
		time.Sleep(1 * time.Second)
		allocs := p.Nomad.Allocs()
		p.UpdatePromtailTargets(allocs)
	}
}

func (p *Promtail) UpdatePromtailTargets(Allocs []*api.Allocation) {
	config := []*ScrapeConfig{}

	for _, alloc := range Allocs {
		config = append(config, p.AllocToScrapeConfigs(alloc)...)
	}

	yamlBytes, err := yaml.Marshal(&config)
	if err != nil {
		panic(err)
	}
	yamlString := string(yamlBytes)

	util.WriteConfig(yamlString, p.TargetsFile, "")
}

func (p *Promtail) AllocToScrapeConfigs(Alloc *api.Allocation) []*ScrapeConfig {
	configs := []*ScrapeConfig{}

	tasks, err := nomad.AllocTasks(Alloc)
	if err != nil {
		log.Fatalln(err)
	}

	for _, task := range tasks {
		configs = append(configs, p.AllocTaskStreamToScrapeConfig(Alloc, task, "stdout"))
		configs = append(configs, p.AllocTaskStreamToScrapeConfig(Alloc, task, "stderr"))
	}

	return configs
}

func (p *Promtail) AllocTaskStreamToScrapeConfig(Alloc *api.Allocation, Task *api.Task, Stream string) *ScrapeConfig {
	config := &ScrapeConfig{
		Targets: []string{"localhost"},
		Labels: map[string]string{
			"nomad_namespace":  Alloc.Namespace,
			"nomad_job":        Alloc.JobID,
			"nomad_task_group": Alloc.TaskGroup,
			"nomad_task":       Task.Name,
			"nomad_alloc_id":   Alloc.ID,
			"nomad_alloc_name": Alloc.Name,
			"nomad_node_id":    p.Nomad.NodeID,
			"nomad_log_stream": Stream,

			"__path__": fmt.Sprintf("%s/%s/alloc/logs/%s.%s.[0-9]*", p.Nomad.AllocsDir, Alloc.ID, Task.Name, Stream),
		},
	}
	return config
}

type ScrapeConfig struct {
	Targets []string          `yaml:"targets"`
	Labels  map[string]string `yaml:"labels"`
}
