---
subcategory: "Classic Load Balancer (SLB)"
layout: "alicloud"
page_title: "Alicloud: alicloud_slb_attachment"
sidebar_current: "docs-alicloud-resource-slb-attachment"
description: |-
  Provides an Application Load Banlancer Attachment resource.
---

# alicloud_slb_attachment

-> **DEPRECATED:** This resource has been deprecated from v1.153.0 and using [alicloud_slb_backend_server](https://www.terraform.io/docs/providers/alicloud/r/slb_backend_server) instead.

Add a group of backend servers (ECS instance) to the Server Load Balancer or remove them from it.

## Example Usage

```terraform
variable "name" {
  default = "slb-attachment-example"
}

data "alicloud_zones" "default" {
  available_disk_category     = "cloud_efficiency"
  available_resource_creation = "VSwitch"
}

data "alicloud_instance_types" "default" {
  availability_zone = data.alicloud_zones.default.zones[0].id
  cpu_core_count    = 1
  memory_size       = 2
}

data "alicloud_images" "default" {
  name_regex  = "^ubuntu_18.*64"
  most_recent = true
  owners      = "system"
}

resource "alicloud_vpc" "default" {
  name       = var.name
  cidr_block = "172.16.0.0/16"
}

resource "alicloud_vswitch" "default" {
  vpc_id       = alicloud_vpc.default.id
  cidr_block   = "172.16.0.0/16"
  zone_id      = data.alicloud_zones.default.zones[0].id
  vswitch_name = var.name
}

resource "alicloud_security_group" "default" {
  name   = var.name
  vpc_id = alicloud_vpc.default.id
}

resource "alicloud_instance" "default" {
  image_id                   = data.alicloud_images.default.images[0].id
  instance_type              = data.alicloud_instance_types.default.instance_types[0].id
  internet_charge_type       = "PayByTraffic"
  internet_max_bandwidth_out = "5"
  system_disk_category       = "cloud_efficiency"
  security_groups            = [alicloud_security_group.default.id]
  instance_name              = var.name
  vswitch_id                 = alicloud_vswitch.default.id
}

resource "alicloud_slb_load_balancer" "default" {
  load_balancer_name = var.name
  vswitch_id         = alicloud_vswitch.default.id
  load_balancer_spec = "slb.s1.small"
}

resource "alicloud_slb_attachment" "default" {
  load_balancer_id = alicloud_slb_load_balancer.default.id
  instance_ids     = [alicloud_instance.default.id]
  weight           = 90
}
```

## Argument Reference

The following arguments are supported:

* `load_balancer_id` - (Required) ID of the load balancer.
* `instance_ids` - (Required) A list of instance ids to added backend server in the SLB.
* `weight` - (Optional) Weight of the instances. Valid value range: [0-100]. Default to 100.
* `server_type` - (Optional, Available in 1.60.0+) Type of the instances. Valid value ecs, eni. Default to ecs.
* `delete_protection_validation` - (Optional, Available in 1.63.0+) Checking DeleteProtection of SLB instance before deleting. If true, this resource will not be deleted when its SLB instance enabled DeleteProtection. Default to false.
* `backend_servers` - (Optional)The backend servers of the load balancer.

## Attributes Reference

The following attributes are exported:

* `id` - ID of the resource.

## Import

Load balancer attachment can be imported using the id or load balancer id, e.g.

```shell
$ terraform import alicloud_slb_attachment.example lb-abc123456
```
