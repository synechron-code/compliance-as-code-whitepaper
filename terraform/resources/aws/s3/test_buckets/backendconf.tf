terraform {
  backend "azurerm" {
    key = "testbuckets.def.terraform.tfstate"
  }
}