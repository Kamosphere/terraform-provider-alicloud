package alicloud

import (
	"fmt"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

func dataSourceAlicloudSelectDBClusters() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceAlicloudSelectDBClustersRead,

		Schema: map[string]*schema.Schema{
			"db_instance_id": {
				Type:     schema.TypeString,
				Optional: false,
			},
			"db_cluster_ids": {
				Type:     schema.TypeList,
				Optional: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},
			"tags": tagsSchema(),
			"output_file": {
				Type:     schema.TypeString,
				Optional: true,
			},

			// Computed values
			"clusters": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"db_cluster_id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"db_instance_id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"db_cluster_class": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"db_cluster_description": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"engine": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"engine_version": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"create_time": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"status": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"payment_type": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"cpu": {
							Type:     schema.TypeInt,
							Computed: true,
						},
						"memory": {
							Type:     schema.TypeInt,
							Computed: true,
						},
						"cache_size": {
							Type:     schema.TypeInt,
							Computed: true,
						},
						"region_id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"zone_id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"vpc_id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"params": {
							Type:     schema.TypeList,
							Computed: true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"name": {
										Type:     schema.TypeString,
										Computed: true,
									},
									"optional": {
										Type:     schema.TypeInt,
										Computed: true,
									},
									"comment": {
										Type:     schema.TypeString,
										Computed: true,
									},
									"value": {
										Type:     schema.TypeString,
										Computed: true,
									},
									"param_category": {
										Type:     schema.TypeString,
										Computed: true,
									},
									"default_value": {
										Type:     schema.TypeString,
										Computed: true,
									},
									"is_dynamic": {
										Type:     schema.TypeInt,
										Computed: true,
									},
									"is_user_modifiable": {
										Type:     schema.TypeInt,
										Computed: true,
									},
								},
							},
						},
						"param_change_logs": {
							Type:     schema.TypeList,
							Computed: true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"name": {
										Type:     schema.TypeString,
										Computed: true,
									},
									"old_value": {
										Type:     schema.TypeString,
										Computed: true,
									},
									"new_value": {
										Type:     schema.TypeString,
										Computed: true,
									},
									"gmt_created": {
										Type:     schema.TypeString,
										Computed: true,
									},
									"gmt_modified": {
										Type:     schema.TypeString,
										Computed: true,
									},
									"config_id": {
										Type:     schema.TypeInt,
										Computed: true,
									},
									"is_applied": {
										Type:     schema.TypeInt,
										Computed: true,
									},
								},
							},
						},
					},
				},
			},
		},
	}
}

func dataSourceAlicloudSelectDBClustersRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	selectDBService := SelectDBService{client}

	instanceId := fmt.Sprint(d.GetOk("db_instance_id"))
	instanceResp, err := selectDBService.DescribeSelectDBDbInstance(instanceId)
	if err != nil {
		return WrapErrorf(err, DataDefaultErrorMsg, "alicloud_selectdb_db_clusters", AlibabaCloudSdkGoERROR)
	}
	resp := instanceResp["DBClusterList"]

	var objects []map[string]interface{}

	idsMap := make(map[string]string)
	if v, ok := d.GetOk("db_cluster_ids"); ok {
		for _, vv := range v.([]interface{}) {
			if vv == nil {
				continue
			}
			idsMap[vv.(string)] = vv.(string)
		}
	}

	result, _ := resp.([]interface{})
	if len(idsMap) > 0 {
		for _, v := range result {
			item := v.(map[string]interface{})
			if _, ok := idsMap[fmt.Sprint(item["DBClusterId"])]; !ok {
				continue
			}
			objects = append(objects, item)
		}
	} else {
		for _, v := range result {
			item := v.(map[string]interface{})
			objects = append(objects, item)
		}
	}

	ids := make([]string, 0)
	s := make([]map[string]interface{}, 0)
	for _, object := range objects {
		mapping := map[string]interface{}{
			"status":                 object["Status"],
			"create_time":            object["CreatedTime"],
			"memory":                 object["Memory"],
			"db_cluster_description": object["DbClusterName"],
			"payment_type":           object["ChargeType"].(string),
			"cpu":                    object["CpuCores"],
			"db_instance_id":         fmt.Sprint(instanceId),
			"db_cluster_id":          fmt.Sprint(object["DBClusterId"]),
			"db_cluster_class":       object["DbClusterClass"],
			"cache_size":             fmt.Sprint(object["CacheStorageSizeGB"]),

			"engine":         fmt.Sprint(instanceResp["Engine"]),
			"engine_version": fmt.Sprint(instanceResp["EngineVersion"]),
			"vpc_id":         fmt.Sprint(instanceResp["VpcId"]),
			"zone_id":        fmt.Sprint(instanceResp["ZoneId"]),
			"region_id":      fmt.Sprint(instanceResp["RegionId"]),
		}

		id := fmt.Sprint(instanceId) + ":" + fmt.Sprint(object["DBClusterId"])
		selectDBService := SelectDBService{client}

		configArrayList, err := selectDBService.DescribeSelectDBDbClusterConfig(id)
		if err != nil {
			return WrapError(err)
		}
		configArray := make([]map[string]interface{}, 0)
		for _, v := range configArrayList {
			if m1, ok := v.(map[string]interface{}); ok {
				temp1 := map[string]interface{}{
					"comment":            m1["Comment"],
					"default_value":      m1["DefaultValue"],
					"optional":           m1["Optional"],
					"param_category":     m1["ParamCategory"],
					"value":              m1["Value"],
					"is_user_modifiable": m1["IsUserModifiable"],
					"is_dynamic":         m1["IsDynamic"],
					"name":               m1["Name"],
				}
				// config with default value will not be updated
				if m1["DefaultValue"].(string) != m1["Value"].(string) {
					configArray = append(configArray, temp1)
				}
			}
		}
		mapping["params"] = configArray

		configChangeArrayList, err := selectDBService.DescribeSelectDBDbClusterConfigChangeLog(id)
		if err != nil {
			return WrapError(err)
		}
		configChangeArray := make([]map[string]interface{}, 0)
		for _, v := range configChangeArrayList {
			if m1, ok := v.(map[string]interface{}); ok {
				temp1 := map[string]interface{}{
					"name":         m1["Name"],
					"old_value":    m1["OldValue"],
					"new_value":    m1["NewValue"],
					"is_applied":   m1["isApplied"],
					"gmt_created":  m1["GmtCreated"],
					"gmt_modified": m1["GmtModified"],
					"config_id":    m1["ConfigId"],
				}
				configChangeArray = append(configChangeArray, temp1)
			}
		}
		mapping["param_change_logs"] = configChangeArray

		s = append(s, mapping)
	}

	d.SetId(dataResourceIdHash(ids))
	if err := d.Set("db_cluster_ids", ids); err != nil {
		return WrapError(err)
	}

	if err := d.Set("clusters", s); err != nil {
		return WrapError(err)
	}
	if output, ok := d.GetOk("output_file"); ok && output.(string) != "" {
		writeToFile(output.(string), s)
	}

	return nil
}
