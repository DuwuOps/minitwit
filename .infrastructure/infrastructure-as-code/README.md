# Terraform Setup of Minitwit

This is a guide of how to use these infrastructure-as-code Terraform-files to setup DigitalOcean droplets with minitwit.

## Preparation-steps

1. Install Terraform on the [Terraform install page](https://developer.hashicorp.com/terraform/install)
2. Create a copy of `terraform.tfvars.example` named `terraform.tfvars` in this directory (`.infrastructure/infrastructure-as-code/`)
3. Login to your DigitalOcean account on [DigitalOcean's webapp](https://cloud.digitalocean.com/)
4. On DigitalOcean, select `API` in the left sidemenu, which should lead you to to [Applications & API](https://cloud.digitalocean.com/account/api/tokens)
5. Create a new personal access token with write access.
    - Easiest just to give it "Full Access".
    - Name it `your_itu_initials-computer_type-access_scope` (e.g. `babb-laptop-full`) and select an abritrary expiration date.
6. Copy and paste your new personal access token into `digitalocean_token` in the `terraform.tfvars`-file.
7. Fill out the rest of the variables in the `terraform.tfvars`-file
   - If in doubt, set `env_type = "test"`
   - The SSH key path can be found running the following:
     - Command prompt (bat): `dir %USERPROFILE%\.ssh\*`
     - Bash (Linux/macOS or WSL): `ls -al ~/.ssh`
   - Please insert the path (as bash command) to the **secret** key intp `secret_key_path` (e.g. `~/.ssh/id_rsa_do` or `~/.ssh/id_ed42069`)
   - In the `ssh_vars`'s `username` insert the name you called the personal access token you created in step 5, so `your_itu_initials-computer_type-access_scope` (e.g. `babb-laptop-full`)
   - Please insert the rest of the variable values

## Running Terraform-files

First, make this directory your working directory

```pwsh
cd .infrastructure/infrastructure-as-code/
```

Initialize Terraform files

```pwsh
terraform init
```

Finally, execute all Terraform actions in the working directory. You can remove `--auto-approve` if you want to confirm the steps it's gonna take before it starts.

```pwsh
terraform apply --auto-approve
```

## Destroy Droplets and SSH-keys via CLI

If you don't have DigititalOcean CLI installed yet (this can be checked by typing `doctl version` into your terminal), then you can install the [DigitalOcean CLI tool here](https://docs.digitalocean.com/reference/doctl/how-to/install/).

- If you have not, authenticate yourself via the DigitalOcean CLI tool:

```pwsh
doctl auth init
```

- Delete droplets (replace ENV's value with the used env_type-value):

```pwsh
ENV=""
doctl compute droplet delete --force $ENV-worker-1
doctl compute droplet delete --force $ENV-worker-2
doctl compute droplet delete --force $ENV-manager
doctl compute droplet delete --force $ENV-database
```

- Get your SSH-key's ID:

```pwsh
doctl compute ssh-key list
```

- Delete your SSH-key (replace KEYID's value with your SSH-key's ID):

```pwsh
KEYID=""
doctl compute ssh-key delete --force $KEYID
```
