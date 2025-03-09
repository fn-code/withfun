job "gitea" {
  region = "global"
  datacenters = [
    "DC0",
  ]
  type = "service"
  group "svc" {
    count = 1
    volume "gitea-data" {
      type      = "host"
      source    = "gitea-data"
      read_only = false
    }
    volume "gitea-db" {
      type      = "host"
      source    = "gitea-db"
      read_only = false
    }
    restart {
      attempts = 5
      delay    = "30s"
    }
    task "app" {
      driver = "docker"
      volume_mount {
        volume      = "gitea-data"
        destination = "/data"
        read_only   = false
      }
      config {
        image = "gitea/gitea:linux-arm64"
        port_map {
          http     = 3000
          ssh_pass = 22
        }
      }
      env = {
        "APP_NAME"   = "Gitea: Git with a cup of tea"
        "RUN_MODE"   = "prod"
        "SSH_DOMAIN" = "git.example.com"
        "SSH_PORT"   = "22"
        "ROOT_URL"   = "http://git.example.com/"
        "USER_UID"   = "1002"
        "USER_GID"   = "1002"
        "DB_TYPE"    = "postgres"
        "DB_HOST"    = "${NOMAD_ADDR_db_db}"
        "DB_NAME"    = "gitea"
        "DB_USER"    = "gitea"
        "DB_PASSWD"  = "gitea"
      }
      resources {
        cpu    = 200
        memory = 256
        network {
          port "http" {}
          port "ssh_pass" {
            static = "2222"
          }
        }
      }
      service {
        name = "gitea-gui"
        port = "http"
      }
    }
    task "db" {
      driver = "docker"
      volume_mount {
        volume      = "gitea-db"
        destination = "/var/lib/postgresql/data"
        read_only   = false
      }
      config {
        image = "postgres:10-alpine"
        port_map {
          db = 5432
        }
      }
      env {
        "POSTGRES_USER"     = "gitea"
        "POSTGRES_PASSWORD" = "gitea"
        "POSTGRES_DB"       = "gitea"
      }
      resources {
        cpu    = 200
        memory = 128
        network {
          port "db" {}
        }
      }
    }
  }
}