terraform {
  backend "azurerm" {
    key = "awsconfig_bootstrap.def.terraform.tfstate"
  }
}