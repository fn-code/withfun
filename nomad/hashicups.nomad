job "example-job" {
  group {
    task "simple-task" {
      service {
        name     = "simple-service"
        port     = "8080"
        provider = "consul"
      }
    }
  }
  group "product-api" {
    task "product-api" {
      ## ...
      template {
        data        = <<EOH
{{ range service "database" }}
DB_CONNECTION="host={{ .Address }} port={{ .Port }} user=user password=password dbname=db_name"
{{ end }}
EOH
        destination = "local/env.txt"
        env         = true
      }
    }
  }
}