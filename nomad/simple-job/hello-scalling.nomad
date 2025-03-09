job "hello" {
  datacenters = ["dc1"]
  type        = "service"

  group "hello" {
    scaling {
      enabled = true
      min     = 2
      max     = 3
      policy {}
    }


    network {
      port "http" {
        to = 8080
      }
    }

    service {
      name     = "hello-service"
      tags     = ["global", "hello"]
      port     = "http"
      provider = "consul"

      check {
        name     = "alive"
        type     = "tcp"
        interval = "10s"
        timeout  = "2s"
      }
    }


    task "hello" {
      driver = "docker"

      config {
        image = "ludinnento/hola"
        ports = ["http"]
      }

    }
  }
}