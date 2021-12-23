package main

import (
  "github.com/dmaes/nomad-logger/nomad"
  "github.com/dmaes/nomad-logger/promtail"

  "github.com/alexflint/go-arg"
)


var args struct {
  NomadAddress string `arg:"--nomad-address,env:NOMAD_ADDRESS" default:"http://localhost:4646"`
  NomadAllocsDir string `arg:"--nomad-allocs-dir,env:NOMAD_ALLOCS_DIR" default:"/var/lib/nomad/alloc"`
  NomadNodeID string `arg:"--nomad-node-id,env:NOMAD_NODE_ID" default:""`
  PromtailTargetsFile string `arg:"--promtail-targets-file,env:PROMTAIL_TARGETS_FILE" default:"/etc/promtail/nomad.yaml"`
}

func main() {
  arg.MustParse(&args)

  nomad := &nomad.Nomad {
    Address: args.NomadAddress,
    AllocsDir: args.NomadAllocsDir,
  }

  if args.NomadNodeID != "" {
    nomad.NodeID = args.NomadNodeID
  } else {
    nomad.SetNodeIDFromEnvs()
  }

  promtail := &promtail.Promtail {
    Nomad: nomad,
    TargetsFile: args.PromtailTargetsFile,
  }
  promtail.Run()

}
