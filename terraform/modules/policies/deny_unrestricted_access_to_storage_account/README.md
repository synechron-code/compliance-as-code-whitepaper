# Deny creating unrestricted network access storage account

Deny unrestricted network access in your storage account firewall settings. Instead, configure network rules so only applications from allowed networks can access the storage account. To allow connections from specific internet or on-premise clients, access can be granted to traffic from specific Azure virtual networks or to public internet IP address ranges. This policy also check against a list of whitelisted IP address. Only those are allowed to be added in the storage account firewall.

## Cloud Controls Objectives

This policy help to satisfy the following Common Control Objectives:

| Controls ID  | Objectives |
|---|---|
|SVD030|Protect cloud service network access by limiting access from the appropriate source network only, to prevent unauthorised access by external & public threats.|

## Intended Use

Prevent creating storage without proper proper network access control.

### Variables

management_group_id : the management group Id that's the policy definition is created against.
allowedAddressRanges : the IP address range that is allowed.

## Apply with Terraform

This should be applied to Azure as a policy and then assigned with appropriate parameters. This would be applied with the main azure-policy module.