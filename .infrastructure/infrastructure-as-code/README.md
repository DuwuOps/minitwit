This is a guide of how to use these infrastructure-as-code Terraform-files to setup DigitalOcean droplets with minitwit.

## Preparation-steps

1. Install Terraform (https://developer.hashicorp.com/terraform/install#linux)
2. Create a copy of `terraform.tfvars.example` named `terraform.tfvars` in this directory (`.infrastructure/infrastructure-as-code/`)
3. Login to your DigitalOcean account (https://cloud.digitalocean.com/)
4. On DigitalOcean, go to Applications & API (https://cloud.digitalocean.com/account/api/tokens)
5. Create a new personal access token with write access. 
    - Easiest just to give it "Full Access".
    - Name does not matter and you decide expiration date.
6. Copy and paste your new personal access token into `env_type` in the `terraform.tfvars`-file.
7. Fill out the rest of the variables in the `terraform.tfvars`-file.


## Running Terraform-files

First, make this directory your working directory

```
cd .infrastructure/infrastructure-as-code/
```

Initialize Terraform files

```
terraform init
```

Finally, execute all Terraform actions in the working directory. You can remove `--auto-approve` if you want to confirm the steps it's gonna take before it starts.

```
terraform apply --auto-approve
```
