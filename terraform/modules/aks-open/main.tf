#https://github.com/terraform-providers/terraform-provider-azurerm/tree/master/examples/kubernetes
provider "azurerm" {
  version = "=2.0.0"
  features {}
}

resource "tls_private_key" "ssh_key" {
  algorithm = "RSA"
  rsa_bits  = 4096
}

#build AKS cluster
resource "azurerm_resource_group" "aks_rg" {
  count = var.rg_name == "create_one_for_me" ? 1 : 0

  location = var.location
  name     = "${var.name_prefix}-onprem-aks-rg"
  tags = merge(
    var.tags,
    {
      "env" = var.env
    },
  )
}

//agent_pool_profile block has been superseded by the default_node_pool, ignore IDE warnings
//noinspection MissingProperty
resource "azurerm_kubernetes_cluster" "aks_cluster" {
  name       = "${var.name_prefix}-aks"
  location   = var.location
  dns_prefix = "${var.name_prefix}-aks"
  resource_group_name = element(
    concat(azurerm_resource_group.aks_rg.*.name, [
    var.rg_name]),
    0,
  )

  linux_profile {
    admin_username = var.aks_admin_username

    ssh_key {
      key_data = var.node_public_ssh_key == "generate" ? tls_private_key.ssh_key.public_key_openssh : var.node_public_ssh_key
    }
  }

  //agent_pool_profile block has been superseded by the default_node_pool, ignore IDE warnings
  //noinspection HCLUnknownBlockType
  default_node_pool {
    name = "agentpool"

    enable_auto_scaling = false
    node_count          = var.node_count
    vm_size             = var.node_sku
    os_disk_size_gb     = 30

    # Required for advanced networking
    vnet_subnet_id = var.cluster_subnet_id
  }

  service_principal {
    client_id     = var.cluster_spn_id
    client_secret = var.cluster_spn_secret
  }

  role_based_access_control {
    enabled = true
  }

  network_profile {
    network_plugin = "azure"
  }

  tags = merge(
    var.tags,
    {
      "env" = var.env
    },
  )
}

/*
# Associate route table with subnet after building it, in case of no default route to Internet
resource "azurerm_subnet_route_table_association" "cluster_rt_association" {
  depends_on = ["azurerm_kubernetes_cluster.aks_cluster"]

  subnet_id      = "${azurerm_subnet.cluster_network.id}"
  route_table_id = "${azurerm_route_table.cluster_routes.id}"
}
*/
