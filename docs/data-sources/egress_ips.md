---
# generated by https://github.com/hashicorp/terraform-plugin-docs
page_title: "materialize_egress_ips Data Source - terraform-provider-materialize"
subcategory: ""
description: |-
  
---

# materialize_egress_ips (Data Source)



## Example Usage

```terraform
data "materialize_egress_ips" "all" {}

output "ips" {
  value = data.materialize_egress_ips.all.egress_ips
}
```

<!-- schema generated by tfplugindocs -->
## Schema

### Read-Only

- `egress_ips` (List of String) The egress IPs in the account
- `id` (String) The ID of this resource.
