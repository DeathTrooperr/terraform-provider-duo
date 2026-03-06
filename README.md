# mss.terraform.provider.duo
Custom terraform provider for Cisco Duo Admin API.

## Features
- `duo_group`: Manage Duo Groups.
- `duo_application`: Manage Duo Applications.
- `duo_user`: Manage Duo Users and their attributes (email, realname, notes, status).
- `duo_policy`: Manage Duo Policies.
- `duo_settings`: Manage global Duo settings (lockout threshold, duration, etc).

## Usage

```hcl
terraform {
  required_providers {
    duo = {
      source = "srmullaney/duo"
    }
  }
}

provider "duo" {
  integration_key = "DIXXXXXXXXXXXXXXXXXX"
  secret_key      = "XXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXX"
  api_host        = "api-XXXXXXXX.duosecurity.com"
}

resource "duo_group" "baseline" {
  name        = "Baseline-Group"
  description = "Managed by Terraform"
}
```

## Development

Build the provider:
```bash
go build -o terraform-provider-duo ./cmd/terraform-provider-duo
```
