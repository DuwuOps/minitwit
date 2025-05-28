
variable "env_type" {
  description = "Deployment environment name (e.g. 'prod' or 'test')"
  type        = string
}

variable "digitalocean_token" {
  description = "DigitalOcean API token"
  type        = string
}

variable "enable_backups" {
  description = "Turn on DigitalOcean automated backups for the DB droplet"
  type        = bool
  default     = true
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
    dockerhub_username = string
  })
}

locals {
  docker_install_script = [
      # Note: -o DPkg::Lock::Timeout=20 is a flag for apt and apt-get that makes apt/apt-get wait till dpkg is unlocked by earlier/other commands. (https://unix.stackexchange.com/a/277255)
      # Add Docker's official GPG key (https://docs.docker.com/engine/install/ubuntu/):
      "sudo apt-get update -y -o DPkg::Lock::Timeout=20",
      "sudo apt-get install -y -o DPkg::Lock::Timeout=20 ca-certificates curl",

      # Wait till /var/lib/apt/lists/lock is unlocked (taken from https://askubuntu.com/a/1451841)
      "while sudo fuser /var/lib/apt/lists/lock >/dev/null 2>&1; do",
      "  echo \"/var/lib/apt/lists/lock is locked..\"",
      "  sleep 1",
      "done",

      "sudo install -m 0755 -d /etc/apt/keyrings",
      "sudo curl -fsSL https://download.docker.com/linux/ubuntu/gpg -o /etc/apt/keyrings/docker.asc",
      "sudo chmod a+r /etc/apt/keyrings/docker.asc",

      # Add the repository to Apt sources (https://docs.docker.com/engine/install/ubuntu/):
      "echo \\",
      "  \"deb [arch=$(dpkg --print-architecture) signed-by=/etc/apt/keyrings/docker.asc] https://download.docker.com/linux/ubuntu \\",
      "  $(lsb_release -cs) stable\" | \\",
      "  sudo tee /etc/apt/sources.list.d/docker.list > /dev/null",
      "sudo apt-get update -y -o DPkg::Lock::Timeout=20",

      # Install the latest versions of Docker packages
      "sudo apt-get install -y -o DPkg::Lock::Timeout=20 docker-ce docker-ce-cli containerd.io docker-buildx-plugin docker-compose-plugin",

      # Wait till /var/lib/apt/lists/lock is unlocked (taken from https://askubuntu.com/a/1451841)
      "while sudo fuser /var/lib/apt/lists/lock >/dev/null 2>&1; do",
      "  echo \"/var/lib/apt/lists/lock is locked..\"",
      "  sleep 1",
      "done",

      # Start docker
      "sudo systemctl start docker",
      "sudo systemctl enable docker",
    ]
}