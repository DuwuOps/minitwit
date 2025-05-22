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
    db_host         = string
    db_port         = string
    db_name         = string
    docker_username = string
  })
}

resource "digitalocean_droplet" "minitwit_droplet" {
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
        ### Install Docker ###
        # Add Docker's official GPG key:
        "sudo apt update",
        "sudo apt install -y apt-transport-https ca-certificates curl software-properties-common",
        "curl -fsSL https://download.docker.com/linux/ubuntu/gpg | sudo apt-key add -",
        "sudo add-apt-repository \"deb [arch=amd64] https://download.docker.com/linux/ubuntu focal stable\" -y",
        "apt-cache policy docker-ce",
        # Add the repository to Apt sources:
        "sudo apt update",
        "sudo apt install -y docker-ce",
        "sudo systemctl start docker",
        "sudo systemctl enable docker",

        ### Start Minitwit Application ###
        "DB_USER=${var.docker_vars.db_user} \\",
        "DB_PASSWORD=${var.docker_vars.db_password} \\",
        "DB_HOST=${var.docker_vars.db_host} \\",
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
  value = digitalocean_droplet.minitwit_droplet.ipv4_address
}