job "logging" {
  region = "someregion"
  datacenters = ["somedc"]

  type = "system"

  priority = 100



  group "nomad-logger" {
    count = 1

    task "nomad-logger" {
      driver = "podman"

      config {
        image = "ghcr.io/dmaes/nomad-logger:latest"
        volumes = [
          "/var/lib/nomad/volumes/logging/promtail:/var/lib/promtail",
        ]
      }

      env {
        NOMAD_ADDR = "http://nomad.service.somedc.consul:4646"
        PROMTAIL_TARGETS_FILE = "/var/lib/promtail/nomad.yaml"
      }

      resources {
        cpu = 100
        memory = 25
      }

      leader = true
    }
  }

  group "promtail" {
    count = 1

    task "promtail" {
      driver = "podman"

      config {
        image = "registry.example.com/nomad/logging/promtail"
        volumes = [
          "local/promtail.yaml:/etc/promtail/config.yaml",
          "/var/lib/nomad/alloc:/var/lib/nomad/alloc:ro",
          "/var/lib/nomad/volumes/logging/promtail:/var/lib/promtail",
        ]
      }

      template {
        destination = "local/promtail.yaml"
        data = <<EOF
---
server:
  http_listen_port: 9080
  grpc_listen_port: 0
clients:
  - url: http://loki.example.com/loki/api/v1/push
positions:
  filename: /var/lib/promtail/positions.yaml
scrape_configs:
  - job_name: nomad
    file_sd_configs:
      - files: ['/var/lib/promtail/nomad.yaml']
    pipeline_stages:
      - match:
          selector: '{job=~"nomad-.+"}'
          stages:
            - regex:
                expression: '^(?P<timestamp>[\d-T:\.+]+) (?P<stream>stdout|stderr) . (?P<output>.+)$'
            - labels:
                stream:
            - timestamp:
                format: RFC3339Nano
                source: timestamp
            - output:
                source: output


EOF
      }

      resources {
        cpu = 100
        memory = 100
      }

      leader = true
    }
  }

}

