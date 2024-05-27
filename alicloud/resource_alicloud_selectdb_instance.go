package alicloud

import (
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

func resourceAlicloudSelectDBInstance() *schema.Resource {
	return &schema.Resource{
		Create: resourceAlicloudSelectDBInstanceCreate,
		Read:   resourceAlicloudSelectDBInstanceRead,
		Update: resourceAlicloudSelectDBInstanceUpdate,
		Delete: resourceAlicloudSelectDBInstanceDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},

		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(30 * time.Minute),
			Update: schema.DefaultTimeout(30 * time.Minute),
			Delete: schema.DefaultTimeout(5 * time.Minute),
		},
		Schema: map[string]*schema.Schema{
			"db_instance_class": {
				Type:     schema.TypeString,
				Required: true,
			},
			"cache_size": {
				Type:     schema.TypeInt,
				Required: true,
			},
			"payment_type": {
				Type:         schema.TypeString,
				ValidateFunc: StringInSlice([]string{"PayAsYouGo", "Subscription"}, false),
				Required:     true,
				Default:      "PayAsYouGo",
			},
			"period": {
				Type:         schema.TypeString,
				ValidateFunc: StringInSlice([]string{string(Year), string(Month)}, false),
				Optional:     true,
			},
			"period_time": {
				Type:             schema.TypeInt,
				ValidateFunc:     IntInSlice([]int{1, 2, 3, 4, 5, 6, 7, 8, 9, 12, 24, 36}),
				Optional:         true,
				DiffSuppressFunc: selectdbPostPaidDiffSuppressFunc,
			},
			"engine_version": {
				Type:     schema.TypeString,
				Required: true,
			},
			"zone_id": {
				Type:     schema.TypeString,
				Required: true,
			},
			"vpc_id": {
				Type:     schema.TypeString,
				Required: true,
			},
			"vswitch_id": {
				Type:     schema.TypeString,
				Required: true,
			},

			// flag for public network and update
			"enable_public_network": {
				Type:     schema.TypeBool,
				Required: true,
				Default:  false,
			},
			"region_id": {
				Type:     schema.TypeString,
				Computed: true,
			},

			"db_instance_id": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"maintain_start_time": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"maintain_end_time": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"engine": {
				Type:     schema.TypeString,
				Required: false,
			},
			"engine_minor_version": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"db_instance_description": {
				Type:     schema.TypeString,
				Required: false,
			},
			"status": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"cpu_prepaid": {
				Type:     schema.TypeInt,
				Computed: true,
			},
			"memory_prepaid": {
				Type:     schema.TypeInt,
				Computed: true,
			},
			"cache_size_prepaid": {
				Type:     schema.TypeInt,
				Computed: true,
			},
			"cluster_count_prepaid": {
				Type:     schema.TypeInt,
				Computed: true,
			},
			"cpu_postpaid": {
				Type:     schema.TypeInt,
				Computed: true,
			},
			"memory_postpaid": {
				Type:     schema.TypeInt,
				Computed: true,
			},
			"cache_size_postpaid": {
				Type:     schema.TypeInt,
				Computed: true,
			},
			"cluster_count_postpaid": {
				Type:     schema.TypeInt,
				Computed: true,
			},
			"sub_domain": {
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
			"gmt_expired": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"lock_mode": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"lock_reason": {
				Type:     schema.TypeString,
				Computed: true,
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
						"db_cluster_class": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"db_cluster_description": {
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
					},
				},
			},

			// Computed values
			"cluster_net_infos": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"db_cluster_id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"db_cluster_ip": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"vpc_instance_id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"connection_string": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"net_type": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"vswitch_id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"port_list": {
							Type:     schema.TypeList,
							Computed: true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"port": {
										Type:     schema.TypeString,
										Computed: true,
									},
									"protocal": {
										Type:     schema.TypeString,
										Computed: true,
									},
								},
							},
						},
					},
				},
			},
			"instance_net_infos": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"db_cluster_id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"db_cluster_ip": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"vpc_instance_id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"connection_string": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"net_type": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"vswitch_id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"port_list": {
							Type:     schema.TypeList,
							Computed: true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"port": {
										Type:     schema.TypeString,
										Computed: true,
									},
									"protocal": {
										Type:     schema.TypeString,
										Computed: true,
									},
								},
							},
						},
					},
				},
			},
			"security_ip_lists": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"group_name": {
							Type:     schema.TypeString,
							Required: false,
						},
						"group_tag": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"security_ip_type": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"security_ip_list": {
							Type:     schema.TypeString,
							Required: false,
						},
						"list_net_type": {
							Type:     schema.TypeString,
							Computed: true,
						},
					},
				},
			},
		},
	}
}

func resourceAlicloudSelectDBInstanceCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	selectDBService := SelectDBService{client}
	request, err := buildSelectDBCreateInstanceRequest(d, meta)
	if err != nil {
		return WrapError(err)
	}
	action := "CreateDBInstance"
	resp, err := selectDBService.RequestProcessForSelectDB(request, action, "POST")
	if err != nil {
		return WrapError(err)
	}
	instanceId := resp["Data"].(map[string]interface{})["DBInstanceId"].(string)

	// wait status change from RESOURCE_PREPARING to CREATING
	stateConfPreparing := BuildStateConf([]string{"RESOURCE_PREPARING"}, []string{"CREATING"}, d.Timeout(schema.TimeoutCreate), 5*time.Minute, selectDBService.SelectDBDbClusterStateRefreshFunc(instanceId, []string{"Deleting"}))
	if _, err := stateConfPreparing.WaitForState(); err != nil {
		return WrapErrorf(err, IdMsg, d.Id())
	}
	// wait status change from CREATING to ACTIVATION
	stateConf := BuildStateConf([]string{"CREATING"}, []string{"ACTIVATION"}, d.Timeout(schema.TimeoutCreate), 5*time.Minute, selectDBService.SelectDBDbClusterStateRefreshFunc(instanceId, []string{"Deleting"}))
	if _, err := stateConf.WaitForState(); err != nil {
		return WrapErrorf(err, IdMsg, d.Id())
	}

	d.SetId(fmt.Sprint(instanceId))

	// enable_public_network

	return resourceAlicloudSelectDBInstanceUpdate(d, meta)
}

func resourceAlicloudSelectDBInstanceUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	selectDBService := SelectDBService{client}
	d.Partial(true)

	if d.IsNewResource() {
		d.Partial(false)
		return resourceAlicloudSelectDBInstanceRead(d, meta)
	}

	// 计费转换
	if d.HasChange("payment_type") {
		_, newPayment := d.GetChange("payment_type")

		request := make(map[string]string)
		payment := convertPaymentTypeToChargeType(newPayment)

		if payment == string(PostPaid) {
			request["payment_type"] = string("POSTPAY")
		} else if payment == string(PrePaid) {
			request["payment_type"] = string("PREPAY")
			request["period"] = d.Get("period").(string)
			request["period_time"] = d.Get("period_time").(string)
		}

		if _, err := selectDBService.ModifySelectDBInstancePaymentType(d.Id(), request); err != nil {
			return WrapError(err)
		}
		d.SetPartial("payment_type")

	}

	// 公网
	if !d.IsNewResource() && d.HasChange("enable_public_network") {
		oldNetStatus, newNetStatus := d.GetChange("enable_public_network")
		if oldNetStatus == true && newNetStatus == false {
			if _, err := selectDBService.ReleaseSelectDBInstancePublicConnection(d.Id()); err != nil {
				return WrapError(err)
			}
			// wait status change from NET_DELETING to ACTIVATION
			stateConf := BuildStateConf([]string{"NET_DELETING"}, []string{"ACTIVATION"}, d.Timeout(schema.TimeoutCreate), 5*time.Minute, selectDBService.SelectDBDbClusterStateRefreshFunc(d.Id(), []string{"Deleting"}))
			if _, err := stateConf.WaitForState(); err != nil {
				return WrapErrorf(err, IdMsg, d.Id())
			}
			d.SetPartial("enable_public_network")
			d.SetPartial("instance_net_infos")
		} else if oldNetStatus == false && newNetStatus == true {
			if _, err := selectDBService.AllocateSelectDBInstancePublicConnection(d.Id()); err != nil {
				return WrapError(err)
			}
			// wait status change from NET_CREATING to ACTIVATION
			stateConf := BuildStateConf([]string{"NET_CREATING"}, []string{"ACTIVATION"}, d.Timeout(schema.TimeoutCreate), 5*time.Minute, selectDBService.SelectDBDbClusterStateRefreshFunc(d.Id(), []string{"Deleting"}))
			if _, err := stateConf.WaitForState(); err != nil {
				return WrapErrorf(err, IdMsg, d.Id())
			}
			d.SetPartial("enable_public_network")
			d.SetPartial("instance_net_infos")
		}

	}

	// 升级
	if !d.IsNewResource() && d.HasChange("engine_minor_version") {
		// minor version
		_, newVersion := d.GetChange("engine_minor_version")
		instanceId := fmt.Sprint(d.Id())
		instanceResp, err := selectDBService.DescribeSelectDBDbInstance(instanceId)
		if err != nil {
			return WrapError(err)
		}
		upgradeTargetVersion := ""
		canUpgradeVersion := instanceResp["CanUpgradeVersions"].(map[string]interface{})
		for _, version := range canUpgradeVersion {
			if newVersion.(string) == version.(string) {
				upgradeTargetVersion = newVersion.(string)
				break
			}
		}
		if upgradeTargetVersion == "" {
			return WrapErrorf(err, "Invalid upgrade version for %s, cannot upgrade to %s", d.Id(), newVersion.(string), AlibabaCloudSdkGoERROR)
		}

		// todo maintaintime update
		_, err = selectDBService.UpgradeSelectDBInstanceEngineVersion(d.Id(), upgradeTargetVersion, false)
		if err != nil {
			return WrapErrorf(err, DefaultErrorMsg, d.Id(), "UpgradeSelectDBInstanceEngineVersion", AlibabaCloudSdkGoERROR)
		}
		d.SetPartial("engine_minor_version")
		d.SetPartial("engine_version")

	} else if !d.IsNewResource() && d.HasChange("engine_version") {
		// major version
		_, newVersion := d.GetChange("engine_version")
		instanceId := fmt.Sprint(d.Id())
		instanceResp, err := selectDBService.DescribeSelectDBDbInstance(instanceId)
		if err != nil {
			return WrapError(err)
		}
		upgradeTargetVersion := ""
		canUpgradeVersion := instanceResp["CanUpgradeVersions"].(map[string]interface{})
		for _, version := range canUpgradeVersion {
			if strings.HasPrefix(version.(string), newVersion.(string)) {
				upgradeTargetVersion = version.(string)
				break
			}
		}
		if upgradeTargetVersion == "" {
			return WrapErrorf(err, "Invalid upgrade version for %s, cannot upgrade to %s", d.Id(), newVersion.(string), AlibabaCloudSdkGoERROR)
		}

		// todo maintaintime update
		_, err = selectDBService.UpgradeSelectDBInstanceEngineVersion(d.Id(), upgradeTargetVersion, false)
		if err != nil {
			return WrapErrorf(err, DefaultErrorMsg, d.Id(), "UpgradeSelectDBInstanceEngineVersion", AlibabaCloudSdkGoERROR)
		}
		d.SetPartial("engine_minor_version")
		d.SetPartial("engine_version")
	}

	// 维护时间
	if !d.IsNewResource() && (d.HasChange("maintain_start_time") || d.HasChange("maintain_end_time")) {
		_, newStartTime := d.GetChange("maintain_start_time")
		_, newEndTime := d.GetChange("maintain_end_time")
		_, err := selectDBService.ModifySelectDBInstanceMaintainTime(d.Id(), newStartTime.(string), newEndTime.(string))
		if err != nil {
			return WrapErrorf(err, DefaultErrorMsg, d.Id(), "ModifySelectDBInstanceMaintainTime", AlibabaCloudSdkGoERROR)
		}
		d.SetPartial("maintain_start_time")
		d.SetPartial("maintain_end_time")
	}

	// 描述
	if !d.IsNewResource() && d.HasChange("db_instance_description") {
		_, newDesc := d.GetChange("db_instance_description")
		_, err := selectDBService.ModifySelectDBInstanceDescription(d.Id(), newDesc.(string))
		if err != nil {
			return WrapErrorf(err, DefaultErrorMsg, d.Id(), "ModifySelectDBInstanceDescription", AlibabaCloudSdkGoERROR)
		}
		d.SetPartial("db_instance_description")
	}

	// 安全组
	if !d.IsNewResource() && d.HasChange("security_ip_lists") {
		_, newDesc := d.GetChange("security_ip_lists")
		for _, v := range newDesc.(map[string]interface{}) {
			item := v.(map[string]interface{})
			_, err := selectDBService.ModifySelectDBDbInstanceSecurityIPList(d.Id(), item["group_name"].(string), item["security_ip_list"].(string))
			if err != nil {
				return WrapErrorf(err, DefaultErrorMsg, d.Id(), "ModifySecurityIPList", AlibabaCloudSdkGoERROR)
			}
		}
		d.SetPartial("security_ip_lists")
	}

	d.Partial(false)
	return resourceAlicloudPolarDBClusterRead(d, meta)
}

func resourceAlicloudSelectDBInstanceRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	selectDBService := SelectDBService{client}

	instanceId := fmt.Sprint(d.Id())
	instanceResp, err := selectDBService.DescribeSelectDBDbInstance(instanceId)
	if err != nil {
		if !d.IsNewResource() && NotFoundError(err) {
			d.SetId("")
			return nil
		}
		return WrapError(err)
	}

	d.Set("engine", instanceResp["Engine"])
	d.Set("engine_version", instanceResp["EngineVersion"])
	d.Set("engine_minor_version", instanceResp["EngineMinorVersion"])

	d.Set("region_id", instanceResp["RegionId"])
	d.Set("zone_id", instanceResp["ZoneId"])
	d.Set("vpc_id", instanceResp["VpcId"])
	// d.Set("vswitch_id", instanceResp["VpcId"])

	d.Set("payment_type", convertChargeTypeToPaymentType(instanceResp["ChargeType"]))
	d.Set("cache_size", "")

	d.Set("db_instance_id", fmt.Sprint(instanceId))
	d.Set("maintain_start_time", instanceResp["MaintainStarttime"])
	d.Set("maintain_end_time", instanceResp["MaintainEndtime"])
	d.Set("db_instance_description", instanceResp["Description"])
	d.Set("status", instanceResp["Status"])
	d.Set("sub_domain", instanceResp["SubDomain"])
	d.Set("gmt_created", instanceResp["CreateTime"])
	d.Set("gmt_modified", instanceResp["GmtModified"])
	d.Set("gmt_expired", instanceResp["ExpiredTime"])
	d.Set("lock_mode", instanceResp["LockMode"])
	d.Set("lock_reason", instanceResp["LockReason"])

	// clusters
	clustersArray := make([]map[string]interface{}, 0)
	clusterResp := instanceResp["DBClusterList"]
	result, _ := clusterResp.([]interface{})
	for _, v := range result {
		item := v.(map[string]interface{})
		mapping := map[string]interface{}{
			"status":                 item["Status"],
			"create_time":            item["CreatedTime"],
			"payment_type":           convertChargeTypeToPaymentType(item["ChargeType"]),
			"cpu":                    item["CpuCores"],
			"memory":                 item["Memory"],
			"cache_size":             fmt.Sprint(item["CacheStorageSizeGB"]),
			"db_cluster_id":          fmt.Sprint(item["DBClusterId"]),
			"db_cluster_class":       item["DbClusterClass"],
			"db_cluster_description": item["DbClusterName"],
		}
		clustersArray = append(clustersArray, mapping)
	}
	d.Set("clusters", clustersArray)

	// cpu,mem,cache
	cpuPrepaid := 0
	cpuPostpaid := 0
	memPrepaid := 0
	memPostpaid := 0
	cachePrepaid := 0
	cachePostpaid := 0

	clusterPrepaidCount := 0
	clusterPostpaidCount := 0

	for _, v := range clusterResp.(map[string]interface{}) {
		item := v.(map[string]interface{})
		if item["ChargeType"].(string) == "Postpaid" {
			cpuPostpaid += item["CpuCores"].(int)
			memPostpaid += item["Memory"].(int)
			cachePostpaid += item["CacheStorageSizeGB"].(int)
			clusterPostpaidCount += 1
		}
		if item["ChargeType"].(string) == "Prepaid" {
			cpuPrepaid += item["CpuCores"].(int)
			memPrepaid += item["Memory"].(int)
			cachePrepaid += item["CacheStorageSizeGB"].(int)
			clusterPrepaidCount += 1
		}
	}
	d.Set("cpu_prepaid", cpuPrepaid)
	d.Set("memory_prepaid", memPrepaid)
	d.Set("cache_size_prepaid", cachePrepaid)
	d.Set("cpu_postpaid", cpuPostpaid)
	d.Set("memory_postpaid", memPostpaid)
	d.Set("cache_size_postpaid", cachePostpaid)

	d.Set("cluster_count_prepaid", clusterPrepaidCount)
	d.Set("cluster_count_postpaid", clusterPostpaidCount)

	// net_infos
	netResp, err := selectDBService.DescribeSelectDBDbInstanceNetInfo(d.Id())
	clustersNetArray := make([]map[string]interface{}, 0)
	resultClusterNet, _ := netResp["DBClustersNetInfos"].(map[string]interface{})
	for _, v := range resultClusterNet {
		item := v.(map[string]interface{})
		port_list := make([]map[string]interface{}, 0)
		for _, vv := range item["PortList"].(map[string]interface{}) {
			port_map := map[string]interface{}{
				"port":     vv.(map[string]interface{})["Port"],
				"protocal": vv.(map[string]interface{})["Protocol"],
			}
			port_list = append(port_list, port_map)
		}
		mapping := map[string]interface{}{
			"db_cluster_id":     item["ClusterId"],
			"db_cluster_ip":     item["Ip"],
			"vpc_instance_id":   item["VpcInstanceId"],
			"connection_string": item["ConnectionString"],
			"net_type":          item["NetType"],
			"vswitch_id":        item["VswitchId"],
			"port_list":         port_list,
		}
		clustersNetArray = append(clustersNetArray, mapping)
	}
	d.Set("cluster_net_infos", clustersNetArray)

	instanceNetArray := make([]map[string]interface{}, 0)
	resultInstanceNet, _ := netResp["DBInstanceNetInfos"].(map[string]interface{})
	for _, v := range resultInstanceNet {
		item := v.(map[string]interface{})
		port_list := make([]map[string]interface{}, 0)
		for _, vv := range item["PortList"].(map[string]interface{}) {
			port_map := map[string]interface{}{
				"port":     vv.(map[string]interface{})["Port"],
				"protocal": vv.(map[string]interface{})["Protocol"],
			}
			port_list = append(port_list, port_map)
		}
		mapping := map[string]interface{}{
			"db_cluster_id":     item["ClusterId"],
			"db_cluster_ip":     item["Ip"],
			"vpc_instance_id":   item["VpcInstanceId"],
			"connection_string": item["ConnectionString"],
			"net_type":          item["NetType"],
			"vswitch_id":        item["VswitchId"],
			"port_list":         port_list,
		}
		instanceNetArray = append(instanceNetArray, mapping)
	}
	d.Set("instance_net_infos", instanceNetArray)

	// security ip lists
	securityIpArrayList, err := selectDBService.DescribeSelectDBDbInstanceSecurityIPList(d.Id())
	if err != nil {
		return WrapError(err)
	}
	securityIpArray := make([]map[string]interface{}, 0)
	for _, v := range securityIpArrayList {
		if m1, ok := v.(map[string]interface{}); ok {
			temp1 := map[string]interface{}{
				"group_name":       m1["GroupName"],
				"group_tag":        m1["GroupTag"],
				"security_ip_type": m1["AecurityIPType"],
				"security_ip_list": m1["SecurityIPList"],
				"list_net_type":    m1["WhitelistNetType"],
			}
			securityIpArray = append(securityIpArray, temp1)
		}
	}
	d.Set("security_ip_lists", securityIpArray)

	return nil
}

func resourceAlicloudSelectDBInstanceDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	selectDBService := SelectDBService{client}

	// 计费转换,需要把包月转成按量才能删除
	payment := convertPaymentTypeToChargeType(d.Get("payment_type"))

	if payment == string(PrePaid) {
		request := make(map[string]string)
		request["payment_type"] = string("POSTPAY")
		if _, err := selectDBService.ModifySelectDBInstancePaymentType(d.Id(), request); err != nil {
			return WrapError(err)
		}
	}

	_, err := selectDBService.DeleteSelectDBCluster(d.Id())
	if err != nil {
		return WrapErrorf(err, DefaultErrorMsg, d.Id(), "DeleteDBInstance", AlibabaCloudSdkGoERROR)
	}
	return nil
}

func buildSelectDBCreateInstanceRequest(d *schema.ResourceData, meta interface{}) (map[string]interface{}, error) {
	client := meta.(*connectivity.AliyunClient)

	request := map[string]interface{}{

		"Engine":          "SelectDB",
		"EngineVersion":   Trim(d.Get("engine_version").(string)),
		"DBInstanceClass": d.Get("db_instance_class").(string),
		"RegionId":        client.RegionId,
		"ZoneId":          d.Get("zone_id").(string),
		"VpcId":           d.Get("vpc_id").(string),
		"VSwitchId":       d.Get("vswitch_id").(string),
		"CacheSize":       d.Get("cache_size").(string),
	}

	payType := convertPaymentTypeToChargeType(d.Get("payment_type"))

	if payType == string(PostPaid) {
		request["ChargeType"] = string("POSTPAY")
	} else if payType == string(PrePaid) {
		request["ChargeType"] = string("PREPAY")
		request["Period"] = d.Get("period").(string)
		request["UsedTime"] = strconv.Itoa(d.Get("period_time").(int))
	}

	return request, nil
}

func convertSelectDBInstanceStatusActionFinal(source string) string {
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
