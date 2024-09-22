datacenter           = "dc-jakarta"
node_name            = "node-jakarta-01"
data_dir             = "/Users/restuwahyusaputra/.consul"
server               = true
enable_debug         = true
enable_script_checks = true
rejoin_after_leave   = true
leave_on_terminate   = true
bootstrap            = true
bootstrap_expect     = 1
raft_protocol        = 3
client_addr          = "0.0.0.0"
bind_addr            = "0.0.0.0"

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
    master = "a6c0ab05-1089-4172-b848-5519d78f0675" // Change with your token
    agent  = "a6c0ab05-1089-4172-b848-5519d78f0675" // Change with your token
  }
}

tls {
  defaults = {
    cert_file       = "/Users/restuwahyusaputra/.self_signed/cert.crt"
    key_file        = "/Users/restuwahyusaputra/.self_signed/key.key"
    verify_incoming = false
  }
}

// You can pass register service manual via config like this or you add service through API
service {
  id      = "dev-discovery-api"
  name    = "dev-discovery-api"
  tags    = ["api", "discovery", "dev"]
  address = "127.0.0.1"
  port    = 5000
  checks = [
    {
      id              = "dev-discovery-api-health"
      name            = "dev-discovery-api-health"
      method          = "GET"
      http            = "http://127.0.0.1:5000/api/v1/ping"
      interval        = "1s"
      timeout         = "10s"
      tls_skip_verify = true
    }
  ]
  token = "a6c0ab05-1089-4172-b848-5519d78f0675" // Change with your token
}