---
subcategory: "Network Load Balancer (NLB)"
layout: "alicloud"
page_title: "Alicloud: alicloud_nlb_server_group_server_attachment"
description: |-
  Provides a Alicloud Network Load Balancer (NLB) Server Group Server Attachment resource.
---

# alicloud_nlb_server_group_server_attachment

Provides a Network Load Balancer (NLB) Server Group Server Attachment resource.

Network Server Load Balancer.

For information about Network Load Balancer (NLB) Server Group Server Attachment and how to use it, see [What is Server Group Server Attachment](https://www.alibabacloud.com/help/en/server-load-balancer/latest/addserverstoservergroup-nlb).

-> **NOTE:** Available since v1.192.0.

## Example Usage

Basic Usage

```terraform
variable "name" {
  default = "tf-example"
}
data "alicloud_resource_manager_resource_groups" "default" {}
resource "alicloud_vpc" "default" {
  vpc_name   = var.name
  cidr_block = "10.4.0.0/16"
}
resource "alicloud_nlb_server_group" "default" {
  resource_group_id        = data.alicloud_resource_manager_resource_groups.default.ids.0
  server_group_name        = var.name
  server_group_type        = "Ip"
  connection_drain_timeout = 10
  connection_drain_enabled = true
  vpc_id                   = alicloud_vpc.default.id
  scheduler                = "Wrr"
  protocol                 = "TCP"
  health_check {
    health_check_enabled = false
  }
  address_ip_version = "Ipv4"
}

resource "alicloud_nlb_server_group_server_attachment" "default" {
  server_type     = "Ip"
  server_id       = "10.0.0.0"
  description     = var.name
  port            = 80
  server_group_id = alicloud_nlb_server_group.default.id
  weight          = 100
  server_ip       = "10.0.0.0"
}
```

## Argument Reference

The following arguments are supported:
* `description` - (Optional) The description of the servers.
The description must be 2 to 256 characters in length, and can contain letters, digits, commas (,), periods (.), semicolons (;), forward slashes (/), at signs (@), underscores (\_), and hyphens (-).
* `port` - (Optional, ForceNew, Computed, Int) The port that is used by the backend server. Valid values: `1` to `65535`.
* `server_group_id` - (Required, ForceNew) The ID of the server group.
* `server_id` - (Required, ForceNew) The ID of the server.

  - If the server group type is `Instance`, set the ServerId parameter to the ID of an Elastic Compute Service (ECS) instance, an elastic network interface (ENI), or an elastic container instance. These backend servers are specified by `Ecs`, `Eni`, or `Eci`.
  - If the server group type is `Ip`, set the ServerId parameter to an IP address.
* `server_ip` - (Optional, ForceNew, Computed) The IP address of the server. If the server group type is `Ip`, set the ServerId parameter to an IP address.
* `server_type` - (Required, ForceNew) The type of the backend server. Valid values:

  - `Ecs`: ECS instance
  - `Eni`: ENI
  - `Eci`: an elastic container instance
  - `Ip`: an IP address
* `weight` - (Optional, Computed, Int) The weight of the backend server. Valid values: `0` to `100`. Default value: `100`. If the weight of a backend server is set to `0`, no requests are forwarded to the backend server.


## Attributes Reference

The following attributes are exported:
* `id` - The ID of the resource supplied above.The value is formulated as `<server_group_id>_<server_id>_<server_ip>_<server_type>_<port>`.
* `zone_id` - The zone ID of the server.
* `status` - The status of the resource

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://www.terraform.io/docs/configuration-0-11/resources.html#timeouts) for certain actions:
* `create` - (Defaults to 5 mins) Used when create the Server Group Server Attachment.
* `delete` - (Defaults to 5 mins) Used when delete the Server Group Server Attachment.
* `update` - (Defaults to 5 mins) Used when update the Server Group Server Attachment.

## Import

Network Load Balancer (NLB) Server Group Server Attachment can be imported using the id, e.g.

```shell
$ terraform import alicloud_nlb_server_group_server_attachment.example <server_group_id>_<server_id>_<server_ip>_<server_type>_<port>
```