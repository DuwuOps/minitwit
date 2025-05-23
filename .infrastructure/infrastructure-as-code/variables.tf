
variable "env_type" {
  description = "Deployment environment name (e.g. 'prod' or 'test')"
  type        = string
}

variable "digitalocean_token" {
  description = "DigitalOcean API token"
  type        = string
}

variable "ssh_vars" {
  description = "Variables for SSH Key Pair for DigitalOcean"
  type = object({
    secret_key_path = string
    username        = string
  })
}

locals {
  ssh_key_exists = fileexists(var.ssh_vars.secret_key_path)
}

variable "docker_vars" {
  description = "All variables used by Docker run & compose"
  type = object({
    db_user         = string
    db_password     = string
    db_port         = string
    db_name         = string
    docker_username = string
  })
}