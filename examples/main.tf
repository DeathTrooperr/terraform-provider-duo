terraform {
  required_providers {
    duo = {
      source = "DeathTrooperr/duo"
    }
  }
}

provider "duo" {
  # integration_key = "DIXXXXXXXXXXXXXXXXXX"
  # secret_key      = "XXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXX"
  # api_host        = "api-XXXXXXXX.duosecurity.com"
}

# Example Group
resource "duo_group" "baseline" {
  name        = "Baseline-Group"
  description = "Managed by Terraform"
}

# Example Application
resource "duo_application" "web_sdk" {
  name = "Web SDK App"
  type = "web_sdk"
}

# Example User with Attributes
resource "duo_user" "admin" {
  username = "admin-user"
  realname = "Admin User"
  email    = "admin@example.com"
  status   = "Active"
  notes    = "Baseline user"
}

# Example Policy
resource "duo_policy" "strict" {
  name = "Strict Policy"
}

# Example Global Settings
resource "duo_settings" "global" {
  lockout_threshold = 5
  lockout_duration  = 30
}
