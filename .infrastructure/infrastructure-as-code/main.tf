terraform {
  required_providers {
    digitalocean = {
      source  = "digitalocean/digitalocean"
      version = "~> 2.0"
    }
  }
}

########## Access token ##########
provider "digitalocean" {
  token = var.digitalocean_token
}

########## SSH key ##########
# Only create the key if it doesn't exist on disk
resource "tls_private_key" "ssh_key" {
  count     = local.ssh_key_exists ? 0 : 1
  algorithm = "RSA"
  rsa_bits  = 2048
}

# Save the generated key to disk (if created)
resource "local_file" "private_key" {
  count           = local.ssh_key_exists ? 0 : 1
  content         = tls_private_key.ssh_key[0].private_key_pem
  filename        = var.ssh_vars.secret_key_path
  file_permission = "0600"
}

resource "local_file" "public_key" {
  count           = local.ssh_key_exists ? 0 : 1
  content         = tls_private_key.ssh_key[0].public_key_openssh
  filename        = "${var.ssh_vars.secret_key_path}.pub"
  file_permission = "0644"
}

resource "digitalocean_ssh_key" "default" {
  name       = "${var.ssh_vars.username}-${var.env_type}-key"
  public_key = local.ssh_key_exists ? file("${var.ssh_vars.secret_key_path}.pub") : tls_private_key.ssh_key[0].public_key_openssh
}

########## Setup Database-Droplet via SSH ##########
resource "digitalocean_droplet" "database_droplet" {
  name      = "${var.env_type}-database"
  region    = "ams3"
  size      = "s-1vcpu-1gb"
  image     = "ubuntu-24-10-x64"
  ssh_keys  = [
    digitalocean_ssh_key.default.fingerprint
  ]
  tags = [
    "minitwit",
    "database",
    var.env_type
  ]

  connection {
    type        = "ssh"
    user        = "root"
    private_key = file(var.ssh_vars.secret_key_path)
    host        = self.ipv4_address
  }

  provisioner "remote-exec" {
    inline = local.docker_install_script
  }

  provisioner "remote-exec" {
    inline = [
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

########## Setup Web-Droplet via SSH ##########
resource "digitalocean_droplet" "web_droplet" {
  name      = "${var.env_type}-web"
  region    = "ams3"
  size      = "s-1vcpu-1gb"
  image     = "ubuntu-24-10-x64"
  ssh_keys  = [
    digitalocean_ssh_key.default.fingerprint
  ]
  tags = [
    "minitwit",
    "app",
    var.env_type
  ]

  connection {
    type        = "ssh"
    user        = "root"
    private_key = file(var.ssh_vars.secret_key_path)
    host        = self.ipv4_address
  }

  provisioner "file" {
    source      = "../../docker-compose.yml"
    destination = "docker-compose.yml"
  }

  provisioner "file" {
    source      = "../../docker-compose.deploy.yml"
    destination = "docker-compose.deploy.yml"
  }

  provisioner "remote-exec" {
    inline = local.docker_install_script
  }

  provisioner "remote-exec" {
    inline = [
      ### Start Minitwit Application ###
      "DB_USER=${var.docker_vars.db_user} \\",
      "DB_PASSWORD=${var.docker_vars.db_password} \\",
      "DB_HOST=${digitalocean_droplet.database_droplet.ipv4_address} \\",
      "DB_PORT=${var.docker_vars.db_port} \\",
      "DB_NAME=${var.docker_vars.db_name} \\",
      "DOCKER_USERNAME=${var.docker_vars.dockerhub_username} \\",
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
  value = "http://${digitalocean_droplet.web_droplet.ipv4_address}/"
}