layout: "vmc"
page_title: "VMC: connected_accounts"
sidebar_current: "docs-vmc-datasource-connected-accounts"
description: A connected accounts data source.
---

# vmc_connected_accounts

This data source provides information about a IP pool configured in NSX.

## Example Usage

```hcl
data "vmc_connected_accounts" "my_accounts" {
  org_id = "${data.vmc_org.my_org.id}"
}
```

## Argument Reference

* `org_id` - (Required) ID of the organization