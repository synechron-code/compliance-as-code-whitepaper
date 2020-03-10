output "cluster_master_hostname" {
  value = azurerm_kubernetes_cluster.aks_cluster.fqdn
}

output "rg_name" {
  value = azurerm_kubernetes_cluster.aks_cluster.resource_group_name
}

