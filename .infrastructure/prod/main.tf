terraform {
  required_providers {
    digitalocean = {
      source  = "digitalocean/digitalocean"
      version = "~> 2.0"
    }
  }
}

#From https://medium.com/@lilnya79/getting-started-with-digitalocean-terraform-and-docker-a-step-by-step-guide-ef43b0513f51 
variable "do_token" {
  description = "DigitalOcean API token"
  type        = string
}

provider "digitalocean" {
  token = var.do_token
}

# Create a new SSH key
resource "digitalocean_ssh_key" "default" {
  name       = "Terraform_Minitwit_Key"
  public_key = file("~/.ssh/id_do_rsa.pub")
}

#From https://medium.com/@lilnya79/getting-started-with-digitalocean-terraform-and-docker-a-step-by-step-guide-ef43b0513f51
resource "digitalocean_droplet" "minitwit_droplet" {
  name      = "prod-web"
  region    = "ams3"
  size      = "s-1vcpu-1gb"
  image     = "ubuntu-24-10-x64"
  ssh_keys  = [
    digitalocean_ssh_key.default.fingerprint
  ] 

  connection {
    type        = "ssh"
    user        = "root"
    private_key = file("~/.ssh/id_do_rsa")
    host        = self.ipv4_address
  }

  provisioner "file" {
    source      = "../../docker-compose.yml"
    destination = "~/.deploy/docker-compose.yml"
  }


  provisioner "file" {
    source      = "../../docker-compose.deploy.yml"
    destination = "~/.deploy/docker-compose.deploy.yml"
  }

  provisioner "remote-exec" {
    inline = [
        ### Install Docker on Droplet ###
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
        # Install Docker Compose
        "sudo apt install -y docker-compose-plugin",

        ### Start Minitwit Application ###
        "cd ~/.deploy/",
        "export DB_USER=${DB_USER}",
        "export DB_PASSWORD=${DB_PASSWORD}",
        "export DB_HOST=${DB_HOST}",
        "export DB_PORT=${DB_PORT}",
        "export DB_NAME=${DB_NAME}",
        "export DB_NAME=${DB_NAME}",
        "export DOCKER_USERNAME=${DOCKER_USERNAME}",
        "docker compose -f docker-compose.yml -f docker-compose.deploy.yml up -d --pull always",
        "unset DB_USER",
        "unset DB_PASSWORD",
        "unset DB_HOST",
        "unset DB_PORT",
        "unset DB_NAME",
        "unset DB_NAME",
        "unset DOCKER_USERNAME"
    ]
  }
}

output "droplet_public_ip" {
  value = digitalocean_droplet.minitwit_droplet.ipv4_address
}