package alicloud

import (
	"encoding/json"
	"strings"
	"time"

	"github.com/PaesslerAG/jsonpath"
	util "github.com/alibabacloud-go/tea-utils/service"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
)

type SelectDBService struct {
	client *connectivity.AliyunClient
}

func (s *SelectDBService) RequestProcessForSelectDB(request map[string]interface{}, action string, method string) (object map[string]interface{}, err error) {
	var response map[string]interface{}
	conn, err := s.client.NewSelectDBClient()
	if err != nil {
		return nil, WrapError(err)
	}
	runtime := util.RuntimeOptions{}
	runtime.SetAutoretry(true)
	wait := incrementalWait(3*time.Second, 3*time.Second)
	err = resource.Retry(5*time.Minute, func() *resource.RetryError {
		response, err = conn.DoRequest(StringPointer(action), nil, StringPointer(method), StringPointer("2023-05-22"), StringPointer("AK"), nil, request, &runtime)
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
		if IsExpectedErrors(err, []string{"InvalidDBInstanceId.NotFound"}) {
			return nil, WrapErrorf(err, NotFoundMsg, AlibabaCloudSdkGoERROR)
		}
		return object, WrapErrorf(err, DefaultErrorMsg, request, action, AlibabaCloudSdkGoERROR)
	}
	v, err := jsonpath.Get("$", response)
	if err != nil {
		return object, WrapErrorf(err, FailedGetAttributeMsg, request, "$", response)
	}
	object = v.(map[string]interface{})
	return object, nil
}

func (s *SelectDBService) RequestProcessPageableForSelectDB(request map[string]interface{}, action string, method string, pageItemJsonpath string) (object []map[string]interface{}, err error) {
	var response map[string]interface{}
	conn, err := s.client.NewSelectDBClient()
	if err != nil {
		return nil, WrapError(err)
	}
	request["PageSize"] = PageSizeLarge
	request["PageNumber"] = 1
	var objects []map[string]interface{}

	for {
		runtime := util.RuntimeOptions{}
		runtime.SetAutoretry(true)
		wait := incrementalWait(3*time.Second, 3*time.Second)
		err = resource.Retry(5*time.Minute, func() *resource.RetryError {
			response, err = conn.DoRequest(StringPointer(action), nil, StringPointer(method), StringPointer("2023-05-22"), StringPointer("AK"), nil, request, &runtime)
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
			return objects, WrapErrorf(err, DefaultErrorMsg, request, action, AlibabaCloudSdkGoERROR)
		}
		resp, err := jsonpath.Get(pageItemJsonpath, response)
		if err != nil {
			return objects, WrapErrorf(err, FailedGetAttributeMsg, action, pageItemJsonpath, response)
		}
		result, _ := resp.([]interface{})
		for _, v := range result {
			item := v.(map[string]interface{})
			objects = append(objects, item)
		}
		if len(result) < PageSizeLarge {
			break
		}
		request["PageNumber"] = request["PageNumber"].(int) + 1
	}
	return objects, nil
}

func (s *SelectDBService) DescribeSelectDBDbCluster(id string) (object map[string]interface{}, err error) {

	action := "DescribeDBInstanceAttribute"
	parts, err := ParseResourceId(id, 2)
	if err != nil {
		return nil, WrapError(err)
	}
	request := map[string]interface{}{
		"DBInstanceId": parts[0],
		"RegionId":     s.client.RegionId,
	}

	v, err := s.RequestProcessForSelectDB(request, action, "GET")
	if err != nil {
		return v, err
	}
	clusterInfo := make(map[string]interface{})
	clusterIndex := v["DBClusterList"].(map[string]interface{})
	for _, w := range clusterIndex {
		ws := w.(map[string]interface{})
		if ws["DbClusterId"] == parts[1] {
			clusterInfo = ws
		}
	}
	if clusterInfo == nil {
		return object, WrapErrorf(err, FailedGetAttributeMsg, id, "$.DbClusterId", v)
	}
	return clusterInfo, nil
}

func (s *SelectDBService) SelectDBDbClusterStateRefreshFunc(id string, failStates []string) resource.StateRefreshFunc {
	return func() (interface{}, string, error) {
		object, err := s.DescribeSelectDBDbCluster(id)
		if err != nil {
			if NotFoundError(err) {
				// Set this to nil as if we didn't find anything.
				return nil, "", nil
			}
			return nil, "", WrapError(err)
		}

		for _, failState := range failStates {
			if object["Status"].(string) == failState {
				return object, object["Status"].(string), WrapError(Error(FailedToReachTargetStatus, object["Status"].(string)))
			}
		}
		return object, object["Status"].(string), nil
	}
}

func (s *SelectDBService) DescribeSelectDBDbInstance(id string) (object map[string]interface{}, err error) {

	action := "DescribeDBInstanceAttribute"

	request := map[string]interface{}{
		"DBInstanceId": id,
		"RegionId":     s.client.RegionId,
	}

	return s.RequestProcessForSelectDB(request, action, "GET")

}

func (s *SelectDBService) SelectDBDbInstanceStateRefreshFunc(id string, failStates []string) resource.StateRefreshFunc {
	return func() (interface{}, string, error) {
		object, err := s.DescribeSelectDBDbInstance(id)
		if err != nil {
			if NotFoundError(err) {
				// Set this to nil as if we didn't find anything.
				return nil, "", nil
			}
			return nil, "", WrapError(err)
		}

		for _, failState := range failStates {
			if object["Status"].(string) == failState {
				return object, object["Status"].(string), WrapError(Error(FailedToReachTargetStatus, object["Status"].(string)))
			}
		}
		return object, object["Status"].(string), nil
	}
}

func (s *SelectDBService) DescribeSelectDBDbInstances(ids string) (objects []map[string]interface{}, err error) {

	action := "DescribeDBInstances"

	request := map[string]interface{}{
		"RegionId": s.client.RegionId,
	}
	if len(ids) > 0 {
		instance_ids := strings.Replace(ids, ":", ",", -1)
		request["DBInstanceIds"] = instance_ids
	}
	pageItemJsonpath := "$.Items"

	return s.RequestProcessPageableForSelectDB(request, action, "GET", pageItemJsonpath)

}

func (s *SelectDBService) DescribeDBInstanceAccessWhiteList(id string) (object map[string]interface{}, err error) {

	action := "DescribeSecurityIPList"
	request := map[string]interface{}{
		"DBInstanceId": id,
		"RegionId":     s.client.RegionId,
	}
	v, err := s.RequestProcessForSelectDB(request, action, "GET")
	if err != nil {
		return v, err
	}
	clusterInfo := v["GroupItems"].(map[string]interface{})
	return clusterInfo, nil
}

func (s *SelectDBService) UpdateSelectDBDbClusterConfig(id string, config map[string]string) (object map[string]interface{}, err error) {

	action := "ModifyDBClusterConfig"
	parts, err := ParseResourceId(id, 2)
	if err != nil {
		return nil, WrapError(err)
	}
	configKey := ""
	if parts[0]+"-fe" == parts[1] {
		configKey = "fe.conf"
	} else {
		configKey = "be.conf"
	}
	param, _ := json.Marshal(config)
	request := map[string]interface{}{
		"DBInstanceId": parts[0],
		"DBClusterId":  parts[1],
		"RegionId":     s.client.RegionId,
		"ConfigKey":    configKey,
		"Parameters":   param,
	}

	v, err := s.RequestProcessForSelectDB(request, action, "POST")
	if err != nil {
		return v, err
	}
	return v, nil

}

func (s *SelectDBService) DescribeSelectDBDbClusterConfig(id string) (object map[string]interface{}, err error) {

	action := "DescribeDBClusterConfig"
	parts, err := ParseResourceId(id, 2)
	if err != nil {
		return nil, WrapError(err)
	}
	configKey := ""
	if parts[0]+"-fe" == parts[1] {
		configKey = "fe.conf"
	} else {
		configKey = "be.conf"
	}
	request := map[string]interface{}{
		"DBInstanceId": parts[0],
		"DBClusterId":  parts[1],
		"RegionId":     s.client.RegionId,
		"ConfigKey":    configKey,
	}

	v, err := s.RequestProcessForSelectDB(request, action, "GET")
	if err != nil {
		return v, err
	}
	clusterInfo := v["Data"].(map[string]interface{})["Params"].(map[string]interface{})
	return clusterInfo, nil

}

func (s *SelectDBService) DescribeSelectDBDbClusterConfigChangeLog(id string) (object map[string]interface{}, err error) {

	action := "DescribeDBClusterConfigChangeLogs"
	parts, err := ParseResourceId(id, 2)
	if err != nil {
		return nil, WrapError(err)
	}
	configKey := ""
	if parts[0]+"-fe" == parts[1] {
		configKey = "fe.conf"
	} else {
		configKey = "be.conf"
	}
	currentTime := time.Now()
	currentTimeString := currentTime.Format("1970-01-01 10:00:00")
	// time not support.
	request := map[string]interface{}{
		"DBInstanceId": parts[0],
		"DBClusterId":  parts[1],
		"RegionId":     s.client.RegionId,
		"ConfigKey":    configKey,
		"StartTime":    "1970-01-01 10:00:00",
		"EndTime":      currentTimeString,
	}
	v, err := s.RequestProcessForSelectDB(request, action, "GET")
	if err != nil {
		return v, err
	}
	clusterInfo := v["Data"].(map[string]interface{})["Params"].(map[string]interface{})
	return clusterInfo, nil
}

func (s *SelectDBService) DeleteSelectDBCluster(id string) (object map[string]interface{}, err error) {

	action := "DeleteDBCluster"

	parts, err := ParseResourceId(id, 2)
	if err != nil {
		return nil, WrapError(err)
	}

	request := map[string]interface{}{
		"DBInstanceId": parts[0],
		"DBClusterId":  parts[1],
		"RegionId":     s.client.RegionId,
	}

	return s.RequestProcessForSelectDB(request, action, "POST")

}

func (s *SelectDBService) DeleteSelectDBInstance(id string) (object map[string]interface{}, err error) {

	action := "DeleteSelectDBInstance"

	request := map[string]interface{}{
		"DBInstanceId": id,
		"RegionId":     s.client.RegionId,
	}

	return s.RequestProcessForSelectDB(request, action, "POST")

}

func (s *SelectDBService) ModifySelectDBClusterDescription(id string, newDescription string) (object map[string]interface{}, err error) {

	action := "ModifyBEClusterAttribute"

	parts, err := ParseResourceId(id, 2)
	if err != nil {
		return nil, WrapError(err)
	}

	request := map[string]interface{}{
		"DBInstanceId":          parts[0],
		"DBClusterId":           parts[1],
		"RegionId":              s.client.RegionId,
		"InstanceAttributeType": "DBInstanceDescription",
		"Value":                 newDescription,
	}

	return s.RequestProcessForSelectDB(request, action, "POST")

}

func (s *SelectDBService) ModifySelectDBCluster(id string, newClass string, newCacheSize string) (object map[string]interface{}, err error) {

	action := "ModifyDBCluster"

	parts, err := ParseResourceId(id, 2)
	if err != nil {
		return nil, WrapError(err)
	}

	request := map[string]interface{}{
		"DBInstanceId":   parts[0],
		"DBClusterId":    parts[1],
		"RegionId":       s.client.RegionId,
		"DBClusterClass": newClass,
		"CacheSize":      newCacheSize,
		"Engine":         "selectdb",
	}

	return s.RequestProcessForSelectDB(request, action, "POST")

}

func (s *SelectDBService) UpdateSelectDBClusterStatus(id string, targetStatus string) (object map[string]interface{}, err error) {
	action := ""
	switch targetStatus {
	case "STOPPING":
		action = "StopBECluster"
	case "STARTING":
		action = "StartBECluster"
	case "RESTART":
		action = "RestartDBCluster"
	}
	if action == "" {
		return nil, WrapError(Error(FailedToReachTargetStatus, targetStatus))
	}

	parts, err := ParseResourceId(id, 2)
	if err != nil {
		return nil, WrapError(err)
	}
	request := map[string]interface{}{
		"DBInstanceId": parts[0],
		"DBClusterId":  parts[1],
		"RegionId":     s.client.RegionId,
	}

	return s.RequestProcessForSelectDB(request, action, "POST")

}

func (s *SelectDBService) DescribeSelectDBDbInstanceNetInfo(id string) (object map[string]interface{}, err error) {

	action := "DescribeDBInstanceNetInfo"

	request := map[string]interface{}{
		"DBInstanceId": id,
		"RegionId":     s.client.RegionId,
	}

	v, err := s.RequestProcessForSelectDB(request, action, "GET")
	if err != nil {
		return v, err
	}
	response := map[string]interface{}{
		"DBClustersNetInfos": v["DBClustersNetInfos"].(map[string]interface{}),
		"DBInstanceNetInfos": v["DBInstanceNetInfos"].(map[string]interface{}),
	}
	return response, nil

}

func (s *SelectDBService) DescribeSelectDBDbInstanceSecurityIPList(id string) (object map[string]interface{}, err error) {

	action := "DescribeSecurityIPList"

	request := map[string]interface{}{
		"DBInstanceId": id,
		"RegionId":     s.client.RegionId,
	}

	v, err := s.RequestProcessForSelectDB(request, action, "GET")
	if err != nil {
		return v, err
	}
	clusterInfo := v["GroupItems"].(map[string]interface{})
	return clusterInfo, nil

}

func (s *SelectDBService) ModifySelectDBDbInstanceSecurityIPList(id string, groupName string, newIpList string) (object map[string]interface{}, err error) {

	action := "ModifySecurityIPList"

	request := map[string]interface{}{
		"DBInstanceId":   id,
		"RegionId":       s.client.RegionId,
		"GroupName":      groupName,
		"SecurityIPList": newIpList,
		"ModifyMode":     0,
	}

	return s.RequestProcessForSelectDB(request, action, "POST")

}

func (s *SelectDBService) ModifySelectDBInstanceDescription(id string, newDescription string) (object map[string]interface{}, err error) {

	action := "ModifyDBInstanceAttribute"

	request := map[string]interface{}{
		"DBInstanceId":          id,
		"RegionId":              s.client.RegionId,
		"InstanceAttributeType": "DBInstanceDescription",
		"Value":                 newDescription,
	}

	return s.RequestProcessForSelectDB(request, action, "POST")

}

func (s *SelectDBService) ModifySelectDBInstanceMaintainTime(id string, maintainStartTime string, maintainEndTime string) (object map[string]interface{}, err error) {

	action := "ModifyDBInstanceAttribute"

	// Format: 1970-01-01T02:00Z, only need 02:00Z
	start_time := strings.Split(maintainStartTime, "T")
	end_time := strings.Split(maintainEndTime, "T")
	time_scope := start_time[1] + "-" + end_time[1]

	request := map[string]interface{}{
		"DBInstanceId":          id,
		"RegionId":              s.client.RegionId,
		"InstanceAttributeType": "MaintainTime",
		"Value":                 time_scope,
	}

	return s.RequestProcessForSelectDB(request, action, "POST")

}

func (s *SelectDBService) ModifySelectDBInstancePaymentType(id string, paymentRequest map[string]string) (object map[string]interface{}, err error) {

	action := "ModifyDBInstancePayType"

	request := map[string]interface{}{
		"DBInstanceId": id,
		"RegionId":     s.client.RegionId,
		"ChargeType":   paymentRequest["payment_type"],
	}
	period, exist := paymentRequest["period"]
	if exist {
		request["Period"] = period
	}
	usedTime, exist := paymentRequest["period_time"]
	if exist {
		request["usedTime"] = usedTime
	}

	return s.RequestProcessForSelectDB(request, action, "POST")

}

func (s *SelectDBService) UpgradeSelectDBInstanceEngineVersion(id string, version string, upgradeInMaintainTime bool) (object map[string]interface{}, err error) {

	action := "UpgradeDBInstanceEngineVersion"
	switchTimeMode := 0

	if upgradeInMaintainTime {
		switchTimeMode = 1
	}

	request := map[string]interface{}{
		"DBInstanceId":   id,
		"RegionId":       s.client.RegionId,
		"EngineVersion":  version,
		"SwitchTimeMode": switchTimeMode,
	}

	return s.RequestProcessForSelectDB(request, action, "POST")

}

func (s *SelectDBService) AllocateSelectDBInstancePublicConnection(id string) (object map[string]interface{}, err error) {

	action := "AllocateInstancePublicConnection"

	request := map[string]interface{}{
		"DBInstanceId":           id,
		"RegionId":               s.client.RegionId,
		"ConnectionStringPrefix": id + "-public",
		"NetType":                "Public",
	}

	return s.RequestProcessForSelectDB(request, action, "POST")

}

func (s *SelectDBService) ReleaseSelectDBInstancePublicConnection(id string) (object map[string]interface{}, err error) {

	action := "ReleaseInstancePublicConnection"

	request := map[string]interface{}{
		"DBInstanceId":           id,
		"RegionId":               s.client.RegionId,
		"ConnectionStringPrefix": id + "-public",
	}

	return s.RequestProcessForSelectDB(request, action, "POST")

}
