# GitHub CLI

- **Repo:** cli/cli
- **Description:** GitHub on the command line
- **Category:** Git

## Commands

### gh auth login
Authenticate with GitHub.

### gh repo create
Create a new GitHub repository.

### gh workflow list
List workflows in a repository.

### gh pr create
Create a pull request.

### gh issue create
Create a new issue.

## Install
```bash
type -p curl >/dev/null || (sudo apt update && sudo apt install curl -y)
curl -fsSL https://cli.github.com/packages/githubcli-archive-keyring.gpg -o /usr/share/keyrings/githubcli-archive-keyring.gpg
echo "deb [signed-by=/usr/share/keyrings/githubcli-archive-keyring.gpg] https://cli.github.com/packages stable main" | sudo tee /etc/apt/sources.list.d/github-cli.list > /dev/null
sudo apt update && sudo apt install gh -y
```
