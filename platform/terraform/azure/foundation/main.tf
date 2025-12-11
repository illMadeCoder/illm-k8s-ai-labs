terraform {
  required_version = ">= 1.0"

  # GitLab-managed state backend - configured via CLI in .gitlab-ci.yml
  backend "http" {}

  required_providers {
    azurerm = {
      source  = "hashicorp/azurerm"
      version = "~> 3.0"
    }
    azuread = {
      source  = "hashicorp/azuread"
      version = "~> 2.0"
    }
  }
}

provider "azurerm" {
  features {
    key_vault {
      purge_soft_delete_on_destroy = true
    }
  }
}

provider "azuread" {}

data "azurerm_client_config" "current" {}

# Resource Group for all lab resources
resource "azurerm_resource_group" "lab" {
  name     = var.resource_group_name
  location = var.location

  tags = {
    project     = "illm-k8s-lab"
    environment = "lab"
    managed_by  = "gitlab-ci"
  }
}

# Key Vault for secrets
resource "azurerm_key_vault" "lab" {
  name                = var.key_vault_name
  location            = azurerm_resource_group.lab.location
  resource_group_name = azurerm_resource_group.lab.name
  tenant_id           = data.azurerm_client_config.current.tenant_id
  sku_name            = "standard"

  # Allow the current user (GitLab CI SP) to manage secrets
  access_policy {
    tenant_id = data.azurerm_client_config.current.tenant_id
    object_id = data.azurerm_client_config.current.object_id

    secret_permissions = [
      "Get", "List", "Set", "Delete", "Purge", "Recover"
    ]
  }

  tags = {
    project     = "illm-k8s-lab"
    environment = "lab"
    managed_by  = "gitlab-ci"
  }
}

# Service Principal for External Secrets Operator
resource "azuread_application" "eso" {
  display_name = "illm-k8s-lab-eso"
}

resource "azuread_service_principal" "eso" {
  client_id = azuread_application.eso.client_id
}

resource "azuread_service_principal_password" "eso" {
  service_principal_id = azuread_service_principal.eso.id
  end_date_relative    = "8760h" # 1 year
}

# Grant ESO SP access to Key Vault secrets (read-only)
resource "azurerm_key_vault_access_policy" "eso" {
  key_vault_id = azurerm_key_vault.lab.id
  tenant_id    = data.azurerm_client_config.current.tenant_id
  object_id    = azuread_service_principal.eso.object_id

  secret_permissions = [
    "Get", "List"
  ]
}

# Store ESO credentials in Key Vault (for bootstrapping)
resource "azurerm_key_vault_secret" "eso_client_id" {
  name         = "eso-client-id"
  value        = azuread_application.eso.client_id
  key_vault_id = azurerm_key_vault.lab.id

  depends_on = [azurerm_key_vault.lab]
}

resource "azurerm_key_vault_secret" "eso_client_secret" {
  name         = "eso-client-secret"
  value        = azuread_service_principal_password.eso.value
  key_vault_id = azurerm_key_vault.lab.id

  depends_on = [azurerm_key_vault.lab]
}
