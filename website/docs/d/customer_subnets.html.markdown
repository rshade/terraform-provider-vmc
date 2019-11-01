---
layout: "vmc"
page_title: "VMC: customer_subnets"
sidebar_current: "docs-vmc-datasource-customer-subnets"
description: A customer subnets data source.
---

# vmc_customer_subnets

The switching profile data source provides information about switching profiles configured in NSX. A switching profile is a template that defines the settings of one or more logical switches. There can be both factory default and user defined switching profiles. One example of a switching profile is a quality of service (QoS) profile which defines the QoS settings of all switches that use the defined switch profile.

## Example Usage

```hcl
data "vmc_customer_subnets" "my_subnets" {
	org_id = "${data.vmc_org.my_org.id}"
	region = "us-west-2"
}
```

## Argument Reference

* `org_id` - (Required) 

* `region` - (Required) 

* `num_hosts` - (Optional) 

* `sddc_id` - (Optional) 

* `force_refresh` - (Optional) 

* `instance_type` - (Optional) 