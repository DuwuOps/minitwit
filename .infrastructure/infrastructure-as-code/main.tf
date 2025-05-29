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

  backups = var.enable_backups
  backup_policy {
    plan    = "weekly"
    weekday = "TUE"
    hour    = 4
  }

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

########## Setup Manager Droplet via SSH ##########
resource "digitalocean_droplet" "manager_droplet" {
  name      = "${var.env_type}-manager"
  region    = "ams3"
  size      = "s-1vcpu-1gb"
  image     = "ubuntu-24-10-x64"
  ssh_keys  = [
    digitalocean_ssh_key.default.fingerprint
  ]
  tags = [
    "minitwit",
    "app",
    "manager",
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
      # Initialize Docker Swarm
      "docker swarm init --advertise-addr ${self.ipv4_address}",
      
      # Get join tokens for workers and save it in a file
      "docker swarm join-token worker -q > /tmp/worker.token",

      "mkdir ~/.deploy/",
      "mv -f docker-compose.yml ~/.deploy/docker-compose.yml",
      "mv -f docker-compose.deploy.yml ~/.deploy/docker-compose.deploy.yml"
    ]
  }
}

########## Setup Worker Droplets via SSH ##########
resource "digitalocean_droplet" "worker_droplet" {
  count     = 2
  name      = "${var.env_type}-worker-${count.index + 1}"
  region    = "ams3"
  size      = "s-1vcpu-1gb"
  image     = "ubuntu-24-10-x64"
  ssh_keys  = [
    digitalocean_ssh_key.default.fingerprint
  ]
  tags = [
    "minitwit",
    "app",
    "worker",
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
      # Join the swarm as worker
      "docker swarm join --token $(ssh -o StrictHostKeyChecking=no -i ${var.ssh_vars.secret_key_path} root@${digitalocean_droplet.manager_droplet.ipv4_address} 'cat /tmp/worker.token') ${digitalocean_droplet.manager_droplet.ipv4_address}:2377"
    ]
  }

  depends_on = [digitalocean_droplet.manager_droplet]  # Don't do this before manager is set up
}

########## Deploy Stack on Manager ##########
resource "null_resource" "deploy_stack" {
  connection {
    type        = "ssh"
    user        = "root"
    private_key = file(var.ssh_vars.secret_key_path)
    host        = digitalocean_droplet.manager_droplet.ipv4_address
  }

  provisioner "remote-exec" {
    inline = [
      ### Start Minitwit Application ###
      "cd ~/.deploy/",
      "DB_USER=${var.docker_vars.db_user} \\",
      "DB_PASSWORD=${var.docker_vars.db_password} \\",
      "DB_HOST=${digitalocean_droplet.database_droplet.ipv4_address} \\",
      "DB_PORT=${var.docker_vars.db_port} \\",
      "DB_NAME=${var.docker_vars.db_name} \\",
      "DOCKER_USERNAME=${var.docker_vars.dockerhub_username} \\",
      "docker stack deploy -c docker-compose.yml -c docker-compose.deploy.yml --prune minitwit"
    ]
  }

  depends_on = [digitalocean_droplet.worker_droplet] # Don't do this before workers are set up
}

output "manager_web_address" {
  value = "http://${digitalocean_droplet.manager_droplet.ipv4_address}/"
}

output "worker_ips" {
  value = digitalocean_droplet.worker_droplet[*].ipv4_address
}