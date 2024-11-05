terraform {
  required_providers {
    azurerm = {
      source  = "hashicorp/azurerm"
      version = "~> 3.0.2"
    }
  }

  required_version = ">= 1.1.0"
}

provider "azurerm" {
  features {}
}

resource "azurerm_resource_group" "rg" {
  name     = "NavixRG"
  location = "East US"
}

resource "azurerm_sql_server" "main" {
  name                         = var.sql_server_name
  resource_group_name          = azurerm_resource_group.rg.name
  location                     = azurerm_resource_group.rg.location
  version                      = "12.0" # SQL Server version
  administrator_login          = var.sql_admin_username
  administrator_login_password = var.sql_admin_password

  tags = {
    environment = "demo"
  }
}

# SQL Database
resource "azurerm_sql_database" "main" {
  name                = var.sql_database_name
  resource_group_name = azurerm_resource_group.rg.name
  location            = azurerm_resource_group.rg.location
  server_name         = azurerm_sql_server.main.name
  sku_name            = "Basic"
  max_size_gb         = 2

  tags = {
    environment = "demo"
  }
}

# Optional: SQL Server Firewall Rule to Allow Azure Services
resource "azurerm_sql_firewall_rule" "allow_azure_services" {
  name                = "AllowAzureServices"
  resource_group_name = azurerm_resource_group.rg.name
  server_name         = azurerm_sql_server.main.name
  start_ip_address    = "0.0.0.0"
  end_ip_address      = "0.0.0.0"
}

# Optional: SQL Server Firewall Rule to Allow Specific IP Range
resource "azurerm_sql_firewall_rule" "my_ip_range" {
  name                = "MyIPAddressRange"
  resource_group_name = azurerm_resource_group.rg.name
  server_name         = azurerm_sql_server.main.name
  start_ip_address    = "YOUR_START_IP"
  end_ip_address      = "YOUR_END_IP"
}
