# Increase log verbosity
log_level = "DEBUG"

# Setup data dir
data_dir = "/home/ludinnento/nomad-data"

name = "fncode02-cl"

# Enable the client
client {
    enabled = true
    servers = ["192.168.64.3:4647"]
}

ports {
    http = 5656
}