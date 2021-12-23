package nomad

import (
  "log"
  "os"

  "github.com/hashicorp/nomad/api"
)


type Nomad struct {
  Address string
  AllocsDir string
  NodeID string
}


func (n *Nomad) Client() (*api.Client) {
  config := *api.DefaultConfig()
  config.Address = n.Address
  client, err := api.NewClient(&config)
  if err != nil { log.Fatal(err) }
  return client
}


func (n *Nomad) SetNodeIDFromEnvs() () {
  q := &api.QueryOptions{ Namespace: os.Getenv("NOMAD_NAMESPACE") }
  alloc, _, err := n.Client().Allocations().Info(os.Getenv("NOMAD_ALLOC_ID"), q)
  if err != nil { log.Fatalln(err) };
  log.Printf("Found node id %s using env vars\n", alloc.NodeID)
  n.NodeID = alloc.NodeID
}


func (n *Nomad) Allocs() ([]*api.Allocation) {
  allocs, _, err := n.Client().Nodes().Allocations(n.NodeID, nil)
  if err != nil { panic(err) }
  return allocs
}
