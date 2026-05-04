# Terraform CLI Basics

- **Repo:** https://github.com/rsdenck/skillnux/infrastructure/terraform_pull.go
- **Description:** Infrastructure as Code tool for provisioning resources
- **Category:** Infrastructure

## Commands

### terraform init
Initialize a Terraform working directory.

### terraform plan
Generate and show an execution plan.

### terraform apply
Apply the changes required to reach the desired state.

### terraform destroy
Destroy Terraform-managed infrastructure.

### terraform fmt
Reformat your configuration in the standard style.

### terraform validate
Check whether the configuration is valid.

## Install
```bash
wget -O- https://apt.releases.hashicorp.com/gpg | gpg --dearmor | sudo tee /usr/share/keyrings/hashicorp-archive-keyring.gpg
echo "deb [signed-by=/usr/share/keyrings/hashicorp-archive-keyring.gpg] https://apt.releases.hashicorp.com $(lsb_release -cs) main" | sudo tee /etc/apt/sources.list.d/hashicorp.list
sudo apt update && sudo apt install terraform
```
