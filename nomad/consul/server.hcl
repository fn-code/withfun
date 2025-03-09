server                 = true
bootstrap_expect       = 1
bind_addr              = "127.0.0.1"
client_addr            = "0.0.0.0"
ui                     = true
data_dir               = "/home/ludinnento/consul-data",
datacenter             = "dc1"
data_dir               = "/opt/consul"
encrypt                = "qDOPBEr+/oUVeOFQOnVypxwDaHzLrD+lvjo5vCEBbZ0="
verify_incoming        = true
verify_outgoing        = true
verify_server_hostname = true
connect {
  enabled = true
}

ports {
  grpc = 8502
}

ui_config {
  enabled = true
}