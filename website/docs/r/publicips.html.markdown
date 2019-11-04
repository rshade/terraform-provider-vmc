---
layout: "vmc"
page_title: "VMC: vmc_publicips"
sidebar_current: "docs-vmc-resource-publicips"
description: |- 
  
---

# vmc_publicips

Provides a resource to allocate public IPs for an SDDC.

## Example Usage

```hcl
data "vmc_org" "my_org" {
	id = ""
}

resource "vmc_publicips" "publicip_1" {
	org_id = "${data.vmc_org.my_org.id}"
	sddc_id = "30aa9e93-766d-498b-92aa-75f3b5304a7e"
	name     = "DefaultIPs"
	private_ip = "10.105.167.133"
}
```

## Argument Reference

The following arguments are supported:

* `org_id` - (Required) Organization identifier.

* `sddc_id` - (Required) SDDC identifier.

* `private_ip` - (Required) Workload VM private IP to be assigned the public IP just allocated.

* `name` - (Required) Workload VM private IPs to be assigned the public IP just allocated.