package fluentbit

import (
  "fmt"
  "io/ioutil"
  "log"
  "os/exec"
  "time"

  "github.com/dmaes/nomad-logger/nomad"

  "github.com/hashicorp/nomad/api"
)


type Fluentbit struct {
  Nomad *nomad.Nomad
  ConfFile string
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

  f.WriteConfig(config)
}


func (f *Fluentbit) AllocToConfig(Alloc *api.Allocation) (string) {
  tag := fmt.Sprintf("%s.%s", f.TagPrefix, Alloc.ID)
  path := fmt.Sprintf("%s/%s/alloc/logs/*.[0-9]*", f.Nomad.AllocsDir, Alloc.ID)

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
  add nomad_group %s
  add nomad_alloc_id %s
  add nomad_alloc_name %s
  add nomad_node_id %s
`, tag, path, tag, Alloc.Namespace, Alloc.JobID, Alloc.TaskGroup, Alloc.ID, Alloc.Name, f.Nomad.NodeID)

  return config
}


func (f *Fluentbit) WriteConfig(config string) {
  oldConfig := ""
  oldConfBytes, err := ioutil.ReadFile(f.ConfFile)
  if err != nil && err.Error() != fmt.Sprintf("open %s: no such file or directory", f.ConfFile) {
    log.Fatal(err)
  } else if err == nil { oldConfig = string(oldConfBytes) }

  if oldConfig == config { return }

  log.Print("Updating config")
  ioutil.WriteFile(f.ConfFile, []byte(config), 0644)

  if f.ReloadCmd == "" { return }

  log.Print("Reloading fluentbit")
  out, cmdErr := exec.Command("/bin/sh", "-c", f.ReloadCmd).CombinedOutput()
  log.Print(out)
  if cmdErr != nil { log.Fatal(cmdErr) }
}

