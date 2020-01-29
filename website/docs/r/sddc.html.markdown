---
layout: "vmc"

page_title: "VMC: vmc_sddc"
sidebar_current: "docs-vmc-resource-sddc"

description: |-
  Provides a resource to provision SDDC.
---

# vmc_sddc

Provides a resource to provision SDDC.

## Example Usage

```hcl

provider "vmc" {
  refresh_token = var.api_token
}

data "vmc_org" "my_org" {
  id =  var.org_id
}

data "vmc_connected_accounts" "my_accounts" {
  org_id = data.vmc_org.my_org.id
  account_number = var.aws_account_number
}

data "vmc_customer_subnets" "my_subnets" {
  org_id               = data.vmc_org.my_org.id
  connected_account_id = data.vmc_connected_accounts.my_accounts.ids[0]
  region               = var.sddc_region
}

resource "vmc_sddc" "sddc_1" {
  org_id = data.vmc_org.my_org.id

  sddc_name           = var.sddc_name
  vpc_cidr            = var.vpc_cidr
  num_host            = 3
  provider_type       = "AWS"
  region              = data.vmc_customer_subnets.my_subnets.region
  vxlan_subnet        = var.vxlan_subnet
  delay_account_link  = false
  skip_creating_vxlan = false
  sso_domain          = "vmc.local"

  deployment_type = "SingleAZ"

  account_link_sddc_config {
    customer_subnet_ids  = [data.vmc_customer_subnets.my_subnets.ids[0]]
    connected_account_id = data.vmc_connected_accounts.my_accounts.ids[0]
  }
  timeouts {
    create = "300m"
    update = "300m"
    delete = "180m"
  }
}
```

## Argument Reference

The following arguments are supported:

* `org_id` - (Required) Organization identifier.

* `region` - (Required)  The region of the cloud resources to work in.

* `sddc_name` - (Required) Name of the SDDC.

* `storage_capacity` - (Optional) The storage capacity value to be requested for the sddc primary cluster, 
   in GiBs. If provided, instead of using the direct-attached storage, a capacity value amount of 
   separable storage will be used.

* `num_host` - (Required) The number of hosts.

* `account_link_sddc_config` - (Optional) The account linking configuration object.

* `host_instance_type` -  (Optional) The instance type for the esx hosts in the primary cluster of the SDDC.

* `vpc_cidr` - (Optional) AWS VPC IP range. Only prefix of 16 or 20 is currently supported.

* `sddc_type` - (Optional) Denotes the sddc type , if the value is null or empty, the type is considered 
   as default.

* `vxlan_subnet` - (Optional) VXLAN IP subnet in CIDR for compute gateway.

* `delay_account_link` - (Optional)  Boolean flag identifying whether account linking should be delayed 
   or not for the SDDC.

* `provider_type` - (Optional)  Determines what additional properties are available based on cloud 
   provider. Acceptable values are "ZEROCLOUD" and "AWS" with AWS as the default value.

* `skip_creating_vxlan` - (Optional) Boolean value to skip creating vxlan for compute gateway for SDDC provisioning.

* `sso_domain` - (Optional) The SSO domain name to use for vSphere users. If not specified, vmc.local will be used.

* `sddc_template_id` - (Optional) If provided, configuration from the template will applied to the provisioned SDDC.

* `deployment_type` - (Optional) Denotes if request is for a SingleAZ or a MultiAZ SDDC. Default is SingleAZ.

* `cluster_id` - (Optional) Cluster identifier.

## Attributes Reference

In addition to arguments listed above, the following attributes are exported:

* `id` - SDDC identifier.
