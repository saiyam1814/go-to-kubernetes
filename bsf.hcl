
packages {
  development = ["go@1.21.6", "gotools@0.16.1", "delve@1.22.0"]
  runtime     = ["cacert@3.95"]
}

gomodule {
  name    = "go-server-example"
  src     = "./."
  ldFlags = null
  tags    = null
  doCheck = false
}

oci "dev" {
  name         = "ttl.sh/godev/go:dev"
  cmd          = ["/bin/go-server-example"]
  envVars = ["foo=bar"]
  exposedPorts = ["8080/tcp"]
}
