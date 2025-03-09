# Increase log verbosity
log_level = "DEBUG"

# Setup data dir
data_dir = "/home/ludinnento/nomad-data"

# Enable the server
server {
  enabled          = true
  bootstrap_expect = 1
}

advertise {
  http = "0.0.0.0"
}