---
layout: "vmc"
page_title: "VMC: vmc_ns_service_group"
sidebar_current: "docs-nsxt-ns-service-resource-service_group"
description: |-
  Provides a resource to configure NS service group on NSX-T manager
---

# nsxt_ns_service_group

Provides a resource to configure NS service group on NSX-T manager

## Example Usage

```hcl

data "vmc_org" "my_org" {
	id = %q
}

data "vmc_connected_accounts" "accounts" {
	org_id = "${data.vmc_org.my_org.id}"
}

resource "vmc_sddc" "sddc_1" {
	org_id = "${data.vmc_org.my_org.id}"

	# storage_capacity    = 100
	sddc_name = "SDDC_Name"

	vpc_cidr      = "10.2.0.0/16"
	num_host      = 3
	provider_type = "Provider_type"

	region = "US_WEST_2"

	vxlan_subnet = "192.168.1.0/24"

	delay_account_link  = false
	skip_creating_vxlan = false
	sso_domain          = "vmc.local"

	deployment_type = "SingleAZ"
}
```

## Argument Reference

The following arguments are supported:

* `org_id` - (Required) ID of organization.
* `region` - (Required) 
* `sddc_name` - (Required)
* `storage_capacity` - (Optional) 
* `num_host` - (Required)
* `account_link_sddc_config` - (Optional)
* `vpc_cidr` - (Optional)
* `sddc_type` - (Optional)
* `vxlan_subnet` - (Optional) 
* `delay_account_link` - (Optional)
* `provider_type` - (Optional)
* `skip_creating_vxlan` - (Optional) 
* `sso_domain` - (Optional)
* `sddc_template_id` - (Optional) 
* `deployment_type` - (Optional) 