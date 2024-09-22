datacenter           = "dc-jakarta"
node_name            = "node-jakarta"
data_dir             = "/var/lib/consul/data"
bootstrap            = true
server               = true
enable_debug         = true
enable_script_checks = true
rejoin_after_leave   = true
leave_on_terminate   = true
bootstrap_expect     = 1
raft_protocol        = 3

ports {
  http     = -1
  https    = 8501
  grpc     = -1
  grpc_tls = 8503
}

ui_config {
  enabled = true
}

performance {
  raft_multiplier = 1
}

http_config {
  response_headers = {
    "Access-Control-Allow-Origin" = "*"
  }
}

acl {
  enabled                  = true
  enable_token_persistence = true
  default_policy           = "deny"

  tokens {
    master = "a6c0ab05-1089-4172-b848-5519d78f0675"
    agent  = "a6c0ab05-1089-4172-b848-5519d78f0675"
  }
}

tls {
  defaults = {
    cert_file = "/var/lib/consul/cert.crt"
    key_file  = "/var/lib/consul/key.key"
  }
}