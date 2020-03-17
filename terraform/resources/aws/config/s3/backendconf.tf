terraform {
  backend "azurerm" {
    key = "awsconfig.def.terraform.tfstate"
  }
}