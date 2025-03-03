package alicloud

import (
	"fmt"
	"time"

	"github.com/PaesslerAG/jsonpath"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/helper/validation"
)

func dataSourceAlicloudCenTransitRouterRouteTableAssociations() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceAlicloudCenTransitRouterRouteTableAssociationsRead,
		Schema: map[string]*schema.Schema{
			"status": {
				Type:         schema.TypeString,
				Optional:     true,
				ForceNew:     true,
				ValidateFunc: validation.StringInSlice([]string{"Active", "Creating", "Deleting"}, false),
			},
			"ids": {
				Type:     schema.TypeList,
				Optional: true,
				ForceNew: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
				Computed: true,
			},
			"transit_router_route_table_id": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
			},
			"transit_router_attachment_id": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
			},
			"transit_router_attachment_resource_id": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
			},
			"transit_router_attachment_resource_type": {
				Type:         schema.TypeString,
				Optional:     true,
				ForceNew:     true,
				ValidateFunc: StringInSlice([]string{"VPC", "VBR", "TR", "VPN"}, false),
			},
			"output_file": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"associations": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"resource_id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"resource_type": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"status": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"transit_router_attachment_id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"transit_router_route_table_id": {
							Type:     schema.TypeString,
							Computed: true,
						},
					},
				},
			},
		},
	}
}

func dataSourceAlicloudCenTransitRouterRouteTableAssociationsRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)

	action := "ListTransitRouterRouteTableAssociations"
	request := make(map[string]interface{})
	if v, ok := d.GetOk("transit_router_route_table_id"); ok {
		request["TransitRouterRouteTableId"] = v
	}
	if v, ok := d.GetOk("transit_router_attachment_id"); ok {
		request["TransitRouterAttachmentId"] = v
	}
	if v, ok := d.GetOk("transit_router_attachment_resource_id"); ok {
		request["TransitRouterAttachmentResourceId"] = v
	}
	if v, ok := d.GetOk("transit_router_attachment_resource_type"); ok {
		request["TransitRouterAttachmentResourceType"] = v
	}
	request["MaxResults"] = PageSizeLarge
	var objects []map[string]interface{}

	idsMap := make(map[string]string)
	if v, ok := d.GetOk("ids"); ok {
		for _, vv := range v.([]interface{}) {
			if vv == nil {
				continue
			}
			idsMap[vv.(string)] = vv.(string)
		}
	}
	status, statusOk := d.GetOk("status")
	var response map[string]interface{}
	var err error
	for {
		wait := incrementalWait(3*time.Second, 3*time.Second)
		err = resource.Retry(5*time.Minute, func() *resource.RetryError {
			response, err = client.RpcPost("Cbn", "2017-09-12", action, nil, request, true)
			if err != nil {
				if NeedRetry(err) {
					wait()
					return resource.RetryableError(err)
				}
				return resource.NonRetryableError(err)
			}
			return nil
		})
		addDebug(action, response, request)
		if err != nil {
			return WrapErrorf(err, DataDefaultErrorMsg, "alicloud_cen_transit_router_route_table_associations", action, AlibabaCloudSdkGoERROR)
		}
		resp, err := jsonpath.Get("$.TransitRouterAssociations", response)
		if err != nil {
			return WrapErrorf(err, FailedGetAttributeMsg, action, "$.TransitRouterAssociations", response)
		}
		result, _ := resp.([]interface{})
		for _, v := range result {
			item := v.(map[string]interface{})
			if len(idsMap) > 0 {
				if _, ok := idsMap[fmt.Sprint(item["TransitRouterAttachmentId"])]; !ok {
					continue
				}
			}
			if statusOk && status.(string) != "" && status.(string) != item["Status"].(string) {
				continue
			}
			objects = append(objects, item)
		}
		if nextToken, ok := response["NextToken"].(string); ok && nextToken != "" {
			request["NextToken"] = nextToken
		} else {
			break
		}
	}
	ids := make([]string, 0)
	s := make([]map[string]interface{}, 0)
	for _, object := range objects {
		mapping := map[string]interface{}{
			"resource_id":                   object["ResourceId"],
			"resource_type":                 object["ResourceType"],
			"status":                        object["Status"],
			"id":                            object["TransitRouterAttachmentId"],
			"transit_router_attachment_id":  object["TransitRouterAttachmentId"],
			"transit_router_route_table_id": object["TransitRouterRouteTableId"],
		}
		ids = append(ids, fmt.Sprint(object["TransitRouterAttachmentId"]))
		s = append(s, mapping)
	}

	d.SetId(dataResourceIdHash(ids))
	if err := d.Set("ids", ids); err != nil {
		return WrapError(err)
	}

	if err := d.Set("associations", s); err != nil {
		return WrapError(err)
	}
	if output, ok := d.GetOk("output_file"); ok && output.(string) != "" {
		writeToFile(output.(string), s)
	}

	return nil
}
