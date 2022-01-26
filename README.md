# Nomad Logger

This is a simple Go application that polls the [Nomad](https://www.nomadproject.io/) API for all allocations on a certain host,
and then updates a [promtail](https://grafana.com/docs/loki/latest/clients/promtail/) `file_sd_config` file,
so promtail can scrape and ship logs.


## Usage

```
Usage: nomad-logger [--nomad-address NOMAD-ADDRESS] [--nomad-allocs-dir NOMAD-ALLOCS-DIR] [--nomad-node-id NOMAD-NODE-ID] [--promtail-targets-file PROMTAIL-TARGETS-FILE]

Options:
  --nomad-address NOMAD-ADDRESS
                         The address of the Nomad API [default: http://localhost:4646, env: NOMAD_ADDR]
  --nomad-allocs-dir NOMAD-ALLOCS-DIR
                         The location of the Nomad allocations data. Used to set the path to the logfiles [default: /var/lib/nomad/alloc, env: NOMAD_ALLOCS_DIR]
  --nomad-node-id NOMAD-NODE-ID
                         The ID of the Nomad node to collect logs for. If empty, we'll suppose this also runs in as a nomad job, and the available env vars will be used to determine the Node ID [env: NOMAD_NODE_ID]
  --promtail-targets-file PROMTAIL-TARGETS-FILE
                         The promtail file_sd_config file where the generated config can be written. Will be completely overwritten, so don't put anything else there. [default: /etc/promtail/nomad.yaml, env: PROMTAIL_TARGETS_FILE]
  --help, -h             display this help and exit
```


### Example

There is an example nomad job file in `examples/nomad-job.hcl`.
If this is the only job you are running on that node,
than the resulting promtail `file_sd_config` file will look something like `examples/promtail-nomad.yaml`.


## Installing/Building

You can just `go install github.com/dmaes/nomad-logger@latest` or `git clone` and `go build` this.
