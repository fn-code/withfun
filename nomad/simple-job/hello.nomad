job "hello" {
  datacenters = ["dc1"]
  type        = "service"

  group "hello" {
    count = 1

    network {
      port "http" {
        static = 8080
      }
    }


    task "hello" {
      driver = "docker"

      config {
        image = "thedojoseries/frontend"
        ports = ["http"]
      }

      service {
        name = "hello-service"
        port = "http"
        provider = "consul"
      }

    }
  }
}