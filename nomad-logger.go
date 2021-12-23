package main

import (
  "github.com/dmaes/nomad-logger/nomad"
  "github.com/dmaes/nomad-logger/promtail"

  "github.com/alexflint/go-arg"
)


var args struct {
  NomadAddress string `arg:"--nomad-address,env:NOMAD_ADDRESS" default:"http://localhost:4646" help:"The address of the Nomad API"`
  NomadAllocsDir string `arg:"--nomad-allocs-dir,env:NOMAD_ALLOCS_DIR" default:"/var/lib/nomad/alloc" help:"The location of the Nomad allocations data. Used to set the path to the logfiles"`
  NomadNodeID string `arg:"--nomad-node-id,env:NOMAD_NODE_ID" default:"" help:"The ID of the Nomad node to collect logs for. If empty, we'll suppose this also runs in as a nomad job, and the available env vars will be used to determine the Node ID"`
  PromtailTargetsFile string `arg:"--promtail-targets-file,env:PROMTAIL_TARGETS_FILE" default:"/etc/promtail/nomad.yaml" help:"The promtail file_sd_config file where the generated config can be written. Will be completely overwritten, so don't put anything else there."`
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
