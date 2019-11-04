layout: "vmc"
page_title: "VMC: connected_accounts"
sidebar_current: "docs-vmc-datasource-connected-accounts"
description: A connected accounts data source.
---

# vmc_connected_accounts

The connected accounts data source get a list of connected accounts.

## Example Usage

```hcl
data "vmc_connected_accounts" "my_accounts" {
  org_id = "${data.vmc_org.my_org.id}"
}
```

## Argument Reference

* `org_id` - (Required) Organization identifier.