# Nomad-Deploy-Notifier
[![docker-publish](https://github.com/Allan-Nava/Nomad-Deploy-Notifier/actions/workflows/docker-publish.yml/badge.svg)](https://github.com/Allan-Nava/Nomad-Deploy-Notifier/actions/workflows/docker-publish.yml) [![Go Build](https://github.com/Allan-Nava/Nomad-Deploy-Notifier/actions/workflows/go-build.yml/badge.svg)](https://github.com/Allan-Nava/Nomad-Deploy-Notifier/actions/workflows/go-build.yml)

Send Nomad deployment messages to slack in GO Lang

# Deployment

Application can be launched using docker  
For example inside a Nomad. Don't forget to change Slack stuff.

```hcl
job "nomad-deploy-notifier" {
    datacenters = ["dc1"]
    type = "service"
    group "nomad-deploy-notifier" {
        task "nomad-deploy-notifier" {
            driver = "docker"
            env {
              "SLACK_TOKEN": "SLACK_TOKEN",
              "SLACK_CHANNEL":"SLACK_CHANNEL"
            }
            config {
                image = "ghcr.io/allan-nava/Nomad-Deploy-Notifier:latest"
                network_mode = "host"
            }
        }
    }
}
```



Inspired by https://github.com/drewbailey/nomad-deploy-notifier

