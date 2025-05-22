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
  name       = "Terraform_Minitwit_Test_Env_Key"
  public_key = file("~/.ssh/id_do_test_env_rsa.pub")
}

#From https://medium.com/@lilnya79/getting-started-with-digitalocean-terraform-and-docker-a-step-by-step-guide-ef43b0513f51
resource "digitalocean_droplet" "minitwit_droplet" {
  name      = "test-web"
  region    = "ams3"
  size      = "s-1vcpu-1gb"
  image     = "ubuntu-24-10-x64"
  ssh_keys  = [
    digitalocean_ssh_key.default.fingerprint
  ] 

  connection {
    type        = "ssh"
    user        = "root"
    private_key = file("~/.ssh/id_do_test_env_rsa")
    host        = self.ipv4_address
  }

  provisioner "remote-exec" {
    inline = [
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
        "sudo docker run -d -p 0.0.0.0:80:8000 --restart=always -v /var/minitwit:/app/tmp tingariussorensen/minitwit:latest"
    ]
  }
}

output "droplet_public_ip" {
  value = digitalocean_droplet.minitwit_droplet.ipv4_address
}