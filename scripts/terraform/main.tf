terraform {
  required_providers {
    azurerm = {
      source  = "hashicorp/azurerm"
      version = "=3.0.1"
    }
  }

  required_version = ">= 1.9.0"
}

provider "azurerm" {
  # The AzureRM Provider supports authenticating using via the Azure CLI, a Managed Identity
  # and a Service Principal. More information on the authentication methods supported by
  # the AzureRM Provider can be found here:
  # https://registry.terraform.io/providers/hashicorp/azurerm/latest/docs#authenticating-to-azure

  # The features block allows changing the behaviour of the Azure Provider, more
  # information can be found here:
  # https://registry.terraform.io/providers/hashicorp/azurerm/latest/docs/guides/features-block
  features {}
}

resource "azurerm_resource_group" "navix_rg" {
  name     = "navix_rg"
  location = "eastus2"
}

resource "azurerm_mssql_server" "navix_mssql_server" {
  name                         = "navix_db_server"
  resource_group_name          = azurerm_resource_group.navix_rg.name
  location                     = azurerm_resource_group.navix_rg.location
  version                      = "12.0"
  administrator_login = ""
  administrator_login_password = ""
}
/*
resource " azurerm_mssql_database" "navix_mssql_db" {
  name         = "Navix-db"
  server_id    = azurerm_mssql_server.navix_mssql_server.id
  license_type = "BasePrice"
  max_size_gb  = 2

  tags = {
    foo = "bar"
  }

  # prevent the possibility of accidental data loss
  lifecycle {
    prevent_destroy = true
  }
}
*/
/*
# Optional: SQL Server Firewall Rule to Allow Azure Services
resource "azurerm_sql_firewall_rule" "allow_azure_services" {
  name                = "AllowAzureServices"
  resource_group_name = azurerm_resource_group.rg.name
  server_name         = azurerm_sql_server.main.name
  start_ip_address    = "0.0.0.0"
  end_ip_address      = "0.0.0.0"
}
*/
/*
# Optional: SQL Server Firewall Rule to Allow Specific IP Range
resource "azurerm_sql_firewall_rule" "my_ip_range" {
  name                = "MyIPAddressRange"
  resource_group_name = azurerm_resource_group.rg.name
  server_name         = azurerm_sql_server.main.name
  start_ip_address    = "YOUR_START_IP"
  end_ip_address      = "YOUR_END_IP"
}
*/