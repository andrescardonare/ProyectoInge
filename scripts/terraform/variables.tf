variable "sql_server_name" {
  default = "navix-my-sql-server"
}

variable "sql_database_name" {
  default = "navix-database"
}

variable "sql_admin_username" {
  default = "navix-admin"
}

variable "sql_admin_password" {
  description = "Password for SQL Admin user"
  sensitive = true
}