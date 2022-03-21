# Nomad Logger

This is a simple Go application that polls the [Nomad](https://www.nomadproject.io/) API for all allocations on a certain host,
and then updates the config of you log shipper of choice.

## Usage

```
Usage: nomad-logger [--nomad-addr NOMAD-ADDR] [--nomad-allocs-dir NOMAD-ALLOCS-DIR] [--nomad-node-id NOMAD-NODE-ID] [--reload-cmd RELOAD-CMD] [--log-shipper LOG-SHIPPER] [--fluentbit-conf-file FLUENTBIT-CONF-FILE] [--fluentbit-tag-prefix FLUENTBIT-TAG-PREFIX] [--promtail-targets-file PROMTAIL-TARGETS-FILE]

Options:
  --nomad-addr NOMAD-ADDR
                         The address of the Nomad API [default: http://localhost:4646, env: NOMAD_ADDR]
  --nomad-allocs-dir NOMAD-ALLOCS-DIR
                         The location of the Nomad allocations data. Used to set the path to the logfiles [default: /var/lib/nomad/alloc, env: NOMAD_ALLOCS_DIR]
  --nomad-node-id NOMAD-NODE-ID
                         The ID of the Nomad node to collect logs for. If empty, we'll suppose this also runs in as a nomad job, and the available env vars will be used to determine the Node ID [env: NOMAD_NODE_ID]
  --reload-cmd RELOAD-CMD
                         Optional command to execute after logshipper config has changed. Usefull to signal a service to reload it's config. Valid for fluentbit logshipper. [env: RELOAD_CMD]
  --log-shipper LOG-SHIPPER
                         The logshipper to use. Options: fluentbit, promtail [default: promtail, env: LOG_SHIPPER]
  --fluentbit-conf-file FLUENTBIT-CONF-FILE [default: /etc/fluent-bit/nomad.conf, env: FLUENTBIT_CONF_FILE]
  --fluentbit-tag-prefix FLUENTBIT-TAG-PREFIX [default: nomad, env: FLUENTBIT_TAG_PREFIX]
  --promtail-targets-file PROMTAIL-TARGETS-FILE
                         The promtail file_sd_config file where the generated config can be written. Will be completely overwritten, so don't put anything else there. [default: /etc/promtail/nomad.yaml, env: PROMTAIL_TARGETS_FILE]
  --help, -h             display this help and exit
```


### Example

There is an example nomad job file in `examples/nomad-job.hcl`.
If this is the only job you are running on that node,
than the resulting promtail `file_sd_config` file will look something like `examples/promtail-nomad.yaml`.


## Log shippers

### Promtail

https://grafana.com/docs/loki/latest/clients/promtail/

Will create/update a `file_sd_config` for promtail to use.
Promtail will watch this file for changes, so no need to signal promtail.


### Fluentbit

https://fluentbit.io/

`[INPUT]` and `[FILTER]` stanza's will be written into a dedicated file,
which can be `@INCLUDE`'ed from your main fluentbit config file.

Fluentbit does not watch it's config files.
So you either have to write something that watches to config file,
or use the `--reload-cmd` flag to execute a command every time the config file changes.
(examples: `--reload-cmd 'systemctl restart fluent-bit'`, `--reload-cmd 'touch a-canary-file-watched-by-a-wrapper'`)


## Installing/Building

You can just `go install github.com/dmaes/nomad-logger@latest` or `git clone` and `go build` this.

There is also a container (`ghcr.io/dmaes/nomad-logger`) that you can use.
This container is build for every commit, you can either use the commit sha or `latest` as tag.
