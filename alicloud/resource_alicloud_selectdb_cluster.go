package alicloud

import (
	"fmt"
	"strconv"
	"time"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

func resourceAlicloudSelectDBCluster() *schema.Resource {
	return &schema.Resource{
		Create: resourceAlicloudSelectDBClusterCreate,
		Read:   resourceAlicloudSelectDBClusterRead,
		Update: resourceAlicloudSelectDBClusterUpdate,
		Delete: resourceAlicloudSelectDBClusterDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},

		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(30 * time.Minute),
			Update: schema.DefaultTimeout(30 * time.Minute),
			Delete: schema.DefaultTimeout(5 * time.Minute),
		},
		Schema: map[string]*schema.Schema{
			"db_instance_id": {
				Type:     schema.TypeString,
				Required: true,
			},
			"db_cluster_description": {
				Type:     schema.TypeString,
				Required: true,
			},
			"db_cluster_class": {
				Type:     schema.TypeString,
				Required: true,
			},
			"payment_type": {
				Type:         schema.TypeString,
				ValidateFunc: StringInSlice([]string{"PayAsYouGo", "Subscription"}, false),
				Required:     true,
				ForceNew:     true,
				Default:      "PayAsYouGo",
			},
			"cache_size": {
				Type:     schema.TypeInt,
				Required: false,
			},

			"params": {
				Type: schema.TypeList,
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
							Required: false,
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

			// computed
			"db_cluster_id": {
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
				Required: false,
			},
			"cpu": {
				Type:     schema.TypeInt,
				Computed: true,
			},
			"memory": {
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
							Type:     schema.TypeInt,
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
	}
}

func resourceAlicloudSelectDBClusterCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	selectDBService := SelectDBService{client}
	request, instanceBeforeCreateResp, err := buildSelectDBCreateClusterRequest(d, meta)
	if err != nil {
		return WrapError(err)
	}
	action := "CreateDBCluster"
	_, err = selectDBService.RequestProcessForSelectDB(request, action, "POST")
	if err != nil {
		return WrapError(err)
	}

	instanceAfterResp, err := selectDBService.DescribeSelectDBDbInstance(d.Get("db_instance_id").(string))
	if err != nil {
		return WrapErrorf(err, DefaultErrorMsg, d.Id())
	}

	clustersBefore := instanceBeforeCreateResp["DBClusterList"].(map[string]interface{})
	clustersAfter := instanceAfterResp["DBClusterList"].(map[string]interface{})

	clustersBeforeIndex := make(map[interface{}]bool)
	for _, item := range clustersBefore {
		clustersBeforeIndex[item.(map[string]interface{})["DBClusterId"]] = true
	}

	clusterId := ""
	for _, item := range clustersAfter {
		if !clustersBeforeIndex[item.(map[string]interface{})["DBClusterId"]] {
			clusterId = item.(map[string]interface{})["DBClusterId"].(string)
		}
	}
	d.SetId(fmt.Sprint(d.Get("db_instance_id").(string) + ":" + clusterId))

	// wait cluster status change from Creating to running
	stateConf := BuildStateConf([]string{"CREATING"}, []string{"ACTIVATION"}, d.Timeout(schema.TimeoutCreate), 5*time.Minute, selectDBService.SelectDBDbClusterStateRefreshFunc(d.Id(), []string{"Deleting"}))
	if _, err := stateConf.WaitForState(); err != nil {
		return WrapErrorf(err, IdMsg, d.Id())
	}

	return resourceAlicloudSelectDBClusterUpdate(d, meta)
}

func resourceAlicloudSelectDBClusterUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	selectDBService := SelectDBService{client}
	d.Partial(true)

	if d.IsNewResource() {
		d.Partial(false)
		return resourceAlicloudPolarDBClusterRead(d, meta)
	}

	// 参数
	if d.HasChange("params") {
		oldConfig, newConfig := d.GetChange("params")
		oldConfigMap := oldConfig.(map[string]interface{})
		newConfigMap := newConfig.(map[string]interface{})
		oldConfigMapIndex := make(map[string]string)
		for _, v := range oldConfigMap {
			item := v.(map[string]interface{})
			oldConfigMapIndex[item["Name"].(string)] = item["Value"].(string)
		}
		newConfigMapIndex := make(map[string]string)
		for _, v := range newConfigMap {
			item := v.(map[string]interface{})
			newConfigMapIndex[item["Name"].(string)] = item["Value"].(string)
		}

		diffConfig := make(map[string]string)
		for k, v := range newConfigMapIndex {
			if oldConfigMapIndex[k] != v {
				diffConfig[k] = v
			}
		}

		if _, err := selectDBService.UpdateSelectDBDbClusterConfig(d.Id(), diffConfig); err != nil {
			return WrapError(err)
		}
		d.SetPartial("params")

		// check params whether is applied
		stateConf := BuildStateConf([]string{"RESTART"}, []string{"ACTIVATION"}, d.Timeout(schema.TimeoutUpdate), 10*time.Minute, selectDBService.SelectDBDbClusterStateRefreshFunc(d.Id(), []string{}))
		if _, err := stateConf.WaitForState(); err != nil {
			return WrapErrorf(err, IdMsg, d.Id())
		}

		// param change must update param change log
		configChangeArrayList, err := selectDBService.DescribeSelectDBDbClusterConfigChangeLog(d.Id())
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
		d.Set("param_change_logs", configChangeArray)
		d.SetPartial("param_change_logs")

	}

	// 启停状态
	if d.HasChange("status") {
		_, newStatus := d.GetChange("status")
		_, err := selectDBService.UpdateSelectDBClusterStatus(d.Id(), newStatus.(string))
		if err != nil {
			return WrapErrorf(err, DefaultErrorMsg, d.Id(), "UpdateSelectDBClusterStatus", AlibabaCloudSdkGoERROR)
		}
		newStatusFinal := convertSelectDBClusterStatusActionFinal(newStatus.(string))
		if newStatusFinal == "" {
			return WrapErrorf(err, DefaultErrorMsg, d.Id(), "UpdateSelectDBClusterStatus", AlibabaCloudSdkGoERROR)
		}
		stateConf := BuildStateConf([]string{newStatus.(string)}, []string{newStatusFinal}, d.Timeout(schema.TimeoutUpdate), 5*time.Minute, selectDBService.SelectDBDbClusterStateRefreshFunc(d.Id(), []string{}))
		if _, err := stateConf.WaitForState(); err != nil {
			return WrapErrorf(err, IdMsg, d.Id())
		}
		d.Set("status", newStatusFinal)
		d.SetPartial("status")
	}

	// 扩缩容
	if !d.IsNewResource() && (d.HasChange("db_cluster_class") || d.HasChange("cache_size")) {
		_, newClass := d.GetChange("db_cluster_class")
		_, newCacheSize := d.GetChange("cache_size")
		_, err := selectDBService.ModifySelectDBCluster(d.Id(), newClass.(string), newCacheSize.(string))
		if err != nil {
			return WrapErrorf(err, DefaultErrorMsg, d.Id(), "ModifyDBCluster", AlibabaCloudSdkGoERROR)
		}
		stateConf := BuildStateConf([]string{"RESOURCE_PREPARING"}, []string{"ACTIVATION"}, d.Timeout(schema.TimeoutCreate), 5*time.Minute, selectDBService.SelectDBDbClusterStateRefreshFunc(d.Id(), []string{"DELETING"}))
		if _, err := stateConf.WaitForState(); err != nil {
			return WrapErrorf(err, IdMsg, d.Id())
		}
		d.SetPartial("db_cluster_class")
	}

	// 描述
	if d.HasChange("db_cluster_description") {
		_, newDesc := d.GetChange("db_cluster_description")
		_, err := selectDBService.ModifySelectDBClusterDescription(d.Id(), newDesc.(string))
		if err != nil {
			return WrapErrorf(err, DefaultErrorMsg, d.Id(), "ModifyBEClusterAttribute", AlibabaCloudSdkGoERROR)
		}
		d.SetPartial("db_cluster_description")
	}

	d.Partial(false)
	return resourceAlicloudPolarDBClusterRead(d, meta)
}

func resourceAlicloudSelectDBClusterRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	selectDBService := SelectDBService{client}

	instanceId := fmt.Sprint(d.GetOk("db_instance_id"))
	instanceResp, err := selectDBService.DescribeSelectDBDbInstance(instanceId)
	if err != nil {
		return WrapErrorf(err, DataDefaultErrorMsg, "alicloud_selectdb_db_clusters", AlibabaCloudSdkGoERROR)
	}
	resp := instanceResp["DBClusterList"]

	var object map[string]interface{}
	result, _ := resp.([]interface{})
	for _, v := range result {
		item := v.(map[string]interface{})
		itemId := fmt.Sprint(instanceId + ":" + item["DbClusterId"].(string))
		if itemId == fmt.Sprint(d.Id()) {
			object = item
			break
		}
	}

	if len(object) == 0 {
		if !d.IsNewResource() {
			d.SetId("")
			return nil
		}
		return WrapErrorf(err, DataDefaultErrorMsg, "alicloud_selectdb_db_clusters", AlibabaCloudSdkGoERROR)
	}
	d.Set("status", object["Status"])
	d.Set("create_time", object["CreatedTime"])
	d.Set("memory", object["Memory"])
	d.Set("db_cluster_description", object["DbClusterName"])
	d.Set("payment_type", convertChargeTypeToPaymentType(object["ChargeType"]))
	d.Set("cpu", object["CpuCores"])
	d.Set("db_instance_id", fmt.Sprint(instanceId))
	d.Set("db_cluster_class", object["DbClusterClass"])
	d.Set("cache_size", fmt.Sprint(object["CacheStorageSizeGB"]))

	d.Set("engine", fmt.Sprint(instanceResp["Engine"]))
	d.Set("engine_version", fmt.Sprint(instanceResp["EngineVersion"]))
	d.Set("vpc_id", fmt.Sprint(instanceResp["VpcId"]))
	d.Set("zone_id", fmt.Sprint(instanceResp["ZoneId"]))
	d.Set("region_id", fmt.Sprint(instanceResp["RegionId"]))

	configArrayList, err := selectDBService.DescribeSelectDBDbClusterConfig(d.Id())
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
	d.Set("params", configArray)

	configChangeArrayList, err := selectDBService.DescribeSelectDBDbClusterConfigChangeLog(d.Id())
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
	d.Set("param_change_logs", configChangeArray)

	return nil
}

func resourceAlicloudSelectDBClusterDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	selectDBService := SelectDBService{client}

	_, err := selectDBService.DescribeSelectDBDbCluster(d.Id())
	if err != nil {
		if NotFoundError(err) {
			return nil
		}
		return WrapError(err)
	}

	_, err = selectDBService.DeleteSelectDBCluster(d.Id())
	if err != nil {
		return WrapErrorf(err, DefaultErrorMsg, d.Id(), "DeleteDBCluster", AlibabaCloudSdkGoERROR)
	}
	return nil
}

func buildSelectDBCreateClusterRequest(d *schema.ResourceData, meta interface{}) (map[string]interface{}, map[string]interface{}, error) {
	client := meta.(*connectivity.AliyunClient)
	selectDBService := SelectDBService{client}

	instanceResp, err := selectDBService.DescribeSelectDBDbInstance(d.Get("db_instance_id").(string))
	if err != nil {
		return nil, nil, WrapErrorf(err, DefaultErrorMsg, d.Id())
	}

	request := map[string]interface{}{

		"DBInstanceId":         d.Get("db_instance_id").(string),
		"Engine":               "SelectDB",
		"EngineVersion":        Trim(d.Get("engine_version").(string)),
		"DBClusterClass":       d.Get("db_cluster_class").(string),
		"RegionId":             client.RegionId,
		"ZoneId":               instanceResp["ZoneId"],
		"VpcId":                instanceResp["VpcId"],
		"VSwitchId":            instanceResp["VSwitchId"],
		"CacheSize":            d.Get("cache_size").(string),
		"DBClusterDescription": Trim(d.Get("db_cluster_name").(string)),
	}

	payType := convertPaymentTypeToChargeType(d.Get("payment_type"))

	if payType == string(PostPaid) {
		request["ChargeType"] = string("POSTPAY")
	} else if payType == string(PrePaid) {
		request["ChargeType"] = string("PREPAY")
		request["Period"] = d.Get("period").(string)
		request["UsedTime"] = strconv.Itoa(d.Get("period_time").(int))
	}

	return request, instanceResp, nil
}

func convertSelectDBClusterStatusActionFinal(source string) string {
	action := ""
	switch source {
	case "STOPPING":
		action = "STOPPED"
	case "STARTING":
		action = "ACTIVATION"
	case "RESTART":
		action = "ACTIVATION"
	}
	return action
}
