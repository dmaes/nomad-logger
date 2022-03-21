package main

import (
  "log"

  "github.com/dmaes/nomad-logger/fluentbit"
  "github.com/dmaes/nomad-logger/nomad"
  "github.com/dmaes/nomad-logger/promtail"

  "github.com/alexflint/go-arg"
)


var args struct {
  NomadAddress string `arg:"--nomad-addr,env:NOMAD_ADDR" default:"http://localhost:4646" help:"The address of the Nomad API"`
  NomadAllocsDir string `arg:"--nomad-allocs-dir,env:NOMAD_ALLOCS_DIR" default:"/var/lib/nomad/alloc" help:"The location of the Nomad allocations data. Used to set the path to the logfiles"`
  NomadNodeID string `arg:"--nomad-node-id,env:NOMAD_NODE_ID" default:"" help:"The ID of the Nomad node to collect logs for. If empty, we'll suppose this also runs in as a nomad job, and the available env vars will be used to determine the Node ID"`
  ReloadCmd string `arg:"--reload-cmd,env:RELOAD_CMD" default:"" help:"Optional command to execute after logshipper config has changed. Usefull to signal a service to reload it's config. Valid for fluentbit logshipper."`
  LogShipper string `arg:"--log-shipper,env:LOG_SHIPPER" default:"promtail" help:"The logshipper to use. Options: fluentbit, promtail"`
  FluentbitConfFile string `arg:"--fluentbit-conf-file,env:FLUENTBIT_CONF_FILE" default:"/etc/fluent-bit/nomad.conf" help "The file in which we can write our input's and stuff. Will be completely overwritten, should be '@INCLUDE'ed from main config file."`
  FluentbitTagPrefix string `arg:"--fluentbit-tag-prefix,env:FLUENTBIT_TAG_PREFIX" default:"nomad" help "Prefix to use for fluentbit tags. Full tag will be '$prefix.$allocId"`
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

  switch args.LogShipper {
  case "fluentbit":
    fluentbit := &fluentbit.Fluentbit {
      Nomad: nomad,
      ConfFile: args.FluentbitConfFile,
      TagPrefix: args.FluentbitTagPrefix,
      ReloadCmd: args.ReloadCmd,
    }
    fluentbit.Run()
  case "promtail":
    promtail := &promtail.Promtail {
      Nomad: nomad,
      TargetsFile: args.PromtailTargetsFile,
    }
    promtail.Run()
  default:
    log.Fatalf("Invalid log shipper type '%s'", args.LogShipper)
  }

}
