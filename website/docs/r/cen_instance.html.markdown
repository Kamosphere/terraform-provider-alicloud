---
subcategory: "Cloud Enterprise Network (CEN)"
layout: "alicloud"
page_title: "Alicloud: alicloud_cen_instance"
sidebar_current: "docs-alicloud-resource-cen-instance"
description: |-
  Provides a Alicloud CEN instance resource.
---

# alicloud_cen_instance

Provides a CEN instance resource. Cloud Enterprise Network (CEN) is a service that allows you to create a global network for rapidly building a distributed business system with a hybrid cloud computing solution. CEN enables you to build a secure, private, and enterprise-class interconnected network between VPCs in different regions and your local data centers. CEN provides enterprise-class scalability that automatically responds to your dynamic computing requirements.

For information about CEN and how to use it, see [What is Cloud Enterprise Network](https://www.alibabacloud.com/help/en/cen/developer-reference/api-cbn-2017-09-12-createcen).

-> **NOTE:** Available since v1.15.0.

## Example Usage

Basic Usage

<div style="display: block;margin-bottom: 40px;"><div class="oics-button" style="float: right;position: absolute;margin-bottom: 10px;">
  <a href="https://api.aliyun.com/api-tools/terraform?resource=alicloud_cen_instance&exampleId=582d1a64-5b30-168e-3f05-c804cba15fdddd9611e0&activeTab=example&spm=docs.r.cen_instance.0.582d1a645b&intl_lang=EN_US" target="_blank">
    <img alt="Open in AliCloud" src="https://img.alicdn.com/imgextra/i1/O1CN01hjjqXv1uYUlY56FyX_!!6000000006049-55-tps-254-36.svg" style="max-height: 44px; max-width: 100%;">
  </a>
</div></div>

```terraform
resource "alicloud_cen_instance" "example" {
  cen_instance_name = "tf_example"
  description       = "an example for cen"
}
```
## Argument Reference

The following arguments are supported:

* `name` - (Optional, Deprecated from v1.98.0+) Field `name` has been deprecated from version 1.98.0. Use `cen_instance_name` instead.
* `cen_instance_name` - (Optional, Available in v1.98.0+) The name of the CEN instance. Defaults to null. The name must be 2 to 128 characters in length and can contain letters, numbers, periods (.), underscores (_), and hyphens (-). The name must start with a letter, but cannot start with http:// or https://.
* `description` - (Optional) The description of the CEN instance. Defaults to null. The description must be 2 to 256 characters in length. It must start with a letter, and cannot start with http:// or https://.
* `tags` - (Optional, Available in v1.80.0+) A mapping of tags to assign to the resource.
* `protection_level` - (Optional, Available in 1.76.0+) Indicates the allowed level of CIDR block overlapping. Default value: `REDUCE`: Overlapping CIDR blocks are allowed. However, the overlapping CIDR blocks cannot be identical.
* `status` - (Optional) The Cen Instance current status.

## Timeouts

-> **NOTE:** Available in 1.48.0+.

The `timeouts` block allows you to specify [timeouts](https://www.terraform.io/docs/configuration-0-11/resources.html#timeouts) for certain actions:

* `create` - (Defaults to 6 mins) Used when creating the cen instance (until it reaches the initial `Active` status). 
* `delete` - (Defaults to 10 mins) Used when terminating the cen instance. 

## Attributes Reference

The following attributes are exported:

* `id` - The ID of the CEN instance.
* `status` - The Cen Instance current status.

## Import

CEN instance can be imported using the id, e.g.

```shell
$ terraform import alicloud_cen_instance.example cen-abc123456
```

