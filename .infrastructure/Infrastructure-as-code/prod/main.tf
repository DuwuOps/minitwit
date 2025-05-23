terraform {
  required_providers {
    digitalocean = {
      source  = "digitalocean/digitalocean"
      version = "~> 2.0"
    }
  }
}

# Access token
variable "digitalocean_token" {
  description = "DigitalOcean API token"
  type        = string
}

provider "digitalocean" {
  token = var.digitalocean_token
}


# SSH key
variable "ssh_key_location" {
  description = "DigitalOcean SSH-publickey"
  type        = string
}

resource "digitalocean_ssh_key" "default" {
  name       = "Terraform_Prod_Env_Key"
  public_key = file("${var.ssh_key_location}.pub")
}

# Setup Doplet via SSH
variable "docker_vars" {
 description = "This is a variable of type object"
  type = object({
    db_user         = string
    db_password     = string
    db_port         = string
    db_name         = string
    docker_username = string
  })
}

resource "digitalocean_droplet" "database_droplet" {
  name      = "prod-database"
  region    = "ams3"
  size      = "s-1vcpu-1gb"
  image     = "ubuntu-24-10-x64"
  ssh_keys  = [
    digitalocean_ssh_key.default.fingerprint
  ]
  tags = [
    "minitwit",
    "database",
    "prod"
  ]

  connection {
    type        = "ssh"
    user        = "root"
    private_key = file(var.ssh_vars.secret_key_path)
    host        = self.ipv4_address
  }

  provisioner "remote-exec" {
    inline = [
      # Note: -o DPkg::Lock::Timeout=20 is a flag for apt and apt-get that makes apt/apt-get wait till dpkg is unlocked by earlier/other commands. (https://unix.stackexchange.com/a/277255)
      # Add Docker's official GPG key (https://docs.docker.com/engine/install/ubuntu/):
      "sudo apt-get update -y -o DPkg::Lock::Timeout=20",
      "sudo apt-get install -y -o DPkg::Lock::Timeout=20 ca-certificates curl",
      "sudo install -m 0755 -d /etc/apt/keyrings",
      "sudo curl -fsSL https://download.docker.com/linux/ubuntu/gpg -o /etc/apt/keyrings/docker.asc",
      "sudo chmod a+r +/etc/apt/keyrings/docker.asc",

      # Add the repository to Apt sources (https://docs.docker.com/engine/install/ubuntu/):
      "echo \\",
      "  \"deb [arch=$(dpkg --print-architecture) signed-by=/etc/apt/keyrings/docker.asc] https://download.docker.com/linux/ubuntu \\",
      "  $(lsb_release -cs) stable\" | \\",
      "  sudo tee /etc/apt/sources.list.d/docker.list > /dev/null",
      "sudo apt-get update -y -o DPkg::Lock::Timeout=20",

      "sudo apt-get install -y -o DPkg::Lock::Timeout=20 docker-ce docker-ce-cli containerd.io docker-buildx-plugin docker-compose-plugin",

      # Start docker
      "sudo systemctl start docker",
      "sudo systemctl enable docker",

      # Install PostgreSQL (using Docker instead of system package)
      "sudo docker run -d --name database \\",
      " -e POSTGRES_PASSWORD=${var.docker_vars.db_password} \\",
      " -e POSTGRES_USER=${var.docker_vars.db_user} \\",
      " -e POSTGRES_DB=${var.docker_vars.db_name} \\",
      " -p ${var.docker_vars.db_port}:5432 \\",
      " postgres"
    ]
  }
}

resource "digitalocean_droplet" "web_droplet" {
  name      = "prod-web"
  region    = "ams3"
  size      = "s-1vcpu-1gb"
  image     = "ubuntu-24-10-x64"
  ssh_keys  = [
    digitalocean_ssh_key.default.fingerprint
  ]
  tags = [
    "minitwit",
    "app",
    "prod"
  ]

  connection {
    type        = "ssh"
    user        = "root"
    private_key = file(var.ssh_key_location)
    host        = self.ipv4_address
  }

  provisioner "file" {
    source      = "../../../docker-compose.yml"
    destination = "docker-compose.yml"
  }

  provisioner "file" {
    source      = "../../../docker-compose.deploy.yml"
    destination = "docker-compose.deploy.yml"
  }

  provisioner "remote-exec" {
    inline = [
      # Note: -o DPkg::Lock::Timeout=20 is a flag for apt and apt-get that makes apt/apt-get wait till dpkg is unlocked by earlier/other commands. (https://unix.stackexchange.com/a/277255)
      # Add Docker's official GPG key (https://docs.docker.com/engine/install/ubuntu/):
      "sudo apt-get update -y -o DPkg::Lock::Timeout=20",
      "sudo apt-get install -y -o DPkg::Lock::Timeout=20 ca-certificates curl",
      "sudo install -m 0755 -d /etc/apt/keyrings",
      "sudo curl -fsSL https://download.docker.com/linux/ubuntu/gpg -o /etc/apt/keyrings/docker.asc",
      "sudo chmod a+r +/etc/apt/keyrings/docker.asc",

      # Add the repository to Apt sources (https://docs.docker.com/engine/install/ubuntu/):
      "echo \\",
      "  \"deb [arch=$(dpkg --print-architecture) signed-by=/etc/apt/keyrings/docker.asc] https://download.docker.com/linux/ubuntu \\",
      "  $(lsb_release -cs) stable\" | \\",
      "  sudo tee /etc/apt/sources.list.d/docker.list > /dev/null",
      "sudo apt-get update -y -o DPkg::Lock::Timeout=20",

      "sudo apt-get install -y -o DPkg::Lock::Timeout=20 docker-ce docker-ce-cli containerd.io docker-buildx-plugin docker-compose-plugin",

      # Start docker
      "sudo systemctl start docker",
      "sudo systemctl enable docker",

      ### Start Minitwit Application ###
      "DB_USER=${var.docker_vars.db_user} \\",
      "DB_PASSWORD=${var.docker_vars.db_password} \\",
      "DB_HOST=${digitalocean_droplet.database_droplet.ipv4_address} \\",
      "DB_PORT=${var.docker_vars.db_port} \\",
      "DB_NAME=${var.docker_vars.db_name} \\",
      "DOCKER_USERNAME=${var.docker_vars.docker_username} \\",
      "docker compose \\",
      "  -f docker-compose.yml \\",
      "  -f docker-compose.deploy.yml \\",
      "  up -d --pull always",

      "mkdir ~/.deploy/",
      "mv -f docker-compose.yml ~/.deploy/docker-compose.yml",
      "mv -f docker-compose.deploy.yml ~/.deploy/docker-compose.deploy.yml"
    ]
  }
}

output "droplet_public_ip" {
  value = digitalocean_droplet.web_droplet.ipv4_address
}