package promtail

import (
  "fmt"
  "io/ioutil"
  "log"
  "time"

  "github.com/dmaes/nomad-logger/nomad"

  "github.com/hashicorp/nomad/api"
  "gopkg.in/yaml.v3"
)


type Promtail struct {
  Nomad *nomad.Nomad
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
    config = append(config, p.AllocToScrapeConfig(alloc))
  }

  data, err := yaml.Marshal(&config)
  if err != nil { panic(err) }
  ioutil.WriteFile(p.TargetsFile, data, 0)
}


func (p *Promtail) AllocToScrapeConfig(Alloc *api.Allocation) (*ScrapeConfig) {
  config := &ScrapeConfig{
    Targets: []string{ "localhost" },
    Labels: map[string]string {
      "nomad_namespace": Alloc.Namespace,
      "nomad_job": Alloc.JobID,
      "nomad_group": Alloc.TaskGroup,
      "nomad_alloc_id": Alloc.ID,
      "nomad_alloc_name": Alloc.Name,
      "nomad_node_id": p.Nomad.NodeID,

      // Log files have an integer suffix (e.g. example.stdout.0, example.stderr.123)
      "__path__": fmt.Sprintf("%s/%s/alloc/logs/*.[0-9]*", p.Nomad.AllocsDir, Alloc.ID),
    },
  }
  return config

}


type ScrapeConfig struct {
  Targets []string `yaml:"targets"`
  Labels map[string]string `yaml:"labels"`
}
