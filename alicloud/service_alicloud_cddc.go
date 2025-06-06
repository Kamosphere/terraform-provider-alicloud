package alicloud

import (
	"fmt"
	"time"

	"github.com/PaesslerAG/jsonpath"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

type CddcService struct {
	client *connectivity.AliyunClient
}

func (s *CddcService) DescribeCddcDedicatedHostGroup(id string) (object map[string]interface{}, err error) {
	var response map[string]interface{}
	client := s.client
	action := "DescribeDedicatedHostGroups"
	request := map[string]interface{}{
		"RegionId":             s.client.RegionId,
		"DedicatedHostGroupId": id,
	}
	wait := incrementalWait(3*time.Second, 3*time.Second)
	err = resource.Retry(5*time.Minute, func() *resource.RetryError {
		response, err = client.RpcPost("cddc", "2020-03-20", action, nil, request, true)
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
		return object, WrapErrorf(err, DefaultErrorMsg, id, action, AlibabaCloudSdkGoERROR)
	}
	v, err := jsonpath.Get("$.DedicatedHostGroups.DedicatedHostGroups", response)
	if err != nil {
		return object, WrapErrorf(err, FailedGetAttributeMsg, id, "$.DedicatedHostGroups.DedicatedHostGroups", response)
	}
	if len(v.([]interface{})) < 1 {
		return object, WrapErrorf(NotFoundErr("CDDC", id), NotFoundWithResponse, response)
	} else {
		if fmt.Sprint(v.([]interface{})[0].(map[string]interface{})["DedicatedHostGroupId"]) != id {
			return object, WrapErrorf(NotFoundErr("CDDC", id), NotFoundWithResponse, response)
		}
	}
	object = v.([]interface{})[0].(map[string]interface{})
	return object, nil
}

func (s *CddcService) ListTagResources(id string, resourceType string) (object interface{}, err error) {
	client := s.client
	parts, err := ParseResourceId(id, 2)
	if err != nil {
		err = WrapError(err)
		return
	}
	action := "ListTagResources"
	request := map[string]interface{}{
		"RegionId":     s.client.RegionId,
		"ResourceType": resourceType,
		"ResourceId.1": parts[1],
	}
	tags := make([]interface{}, 0)
	var response map[string]interface{}

	for {
		wait := incrementalWait(3*time.Second, 5*time.Second)
		err = resource.Retry(5*time.Minute, func() *resource.RetryError {
			response, err := client.RpcPost("cddc", "2020-03-20", action, nil, request, true)
			if err != nil {
				if NeedRetry(err) {
					wait()
					return resource.RetryableError(err)
				}
				return resource.NonRetryableError(err)
			}
			addDebug(action, response, request)
			v, err := jsonpath.Get("$.TagResources", response)
			if err != nil {
				return resource.NonRetryableError(WrapErrorf(err, FailedGetAttributeMsg, id, "$.TagResources", response))
			}
			if v != nil {
				tags = append(tags, v.([]interface{})...)
			}
			return nil
		})
		if err != nil {
			err = WrapErrorf(err, DefaultErrorMsg, id, action, AlibabaCloudSdkGoERROR)
			return
		}
		if response["NextToken"] == nil {
			break
		}
		request["NextToken"] = response["NextToken"]
	}

	return tags, nil
}

func (s *CddcService) SetResourceTags(d *schema.ResourceData, resourceType string) error {

	parts, err := ParseResourceId(d.Id(), 2)
	if err != nil {
		return WrapError(err)
	}
	if d.HasChange("tags") {
		client := s.client
		added, removed := parsingTags(d)
		removedTagKeys := make([]string, 0)
		for _, v := range removed {
			if !ignoredTags(v, "") {
				removedTagKeys = append(removedTagKeys, v)
			}
		}
		if len(removedTagKeys) > 0 {
			action := "UntagResources"
			request := map[string]interface{}{
				"RegionId":     s.client.RegionId,
				"ResourceType": resourceType,
				"ResourceId.1": parts[1],
			}
			for i, key := range removedTagKeys {
				request[fmt.Sprintf("TagKey.%d", i+1)] = key
			}
			wait := incrementalWait(2*time.Second, 1*time.Second)
			err := resource.Retry(10*time.Minute, func() *resource.RetryError {
				response, err := client.RpcPost("cddc", "2020-03-20", action, nil, request, false)
				if err != nil {
					if NeedRetry(err) {
						wait()
						return resource.RetryableError(err)

					}
					return resource.NonRetryableError(err)
				}
				addDebug(action, response, request)
				return nil
			})
			if err != nil {
				return WrapErrorf(err, DefaultErrorMsg, d.Id(), action, AlibabaCloudSdkGoERROR)
			}
		}
		if len(added) > 0 {
			action := "TagResources"
			request := map[string]interface{}{
				"RegionId":     s.client.RegionId,
				"ResourceType": resourceType,
				"ResourceId.1": parts[1],
			}
			count := 1
			for key, value := range added {
				request[fmt.Sprintf("Tag.%d.Key", count)] = key
				request[fmt.Sprintf("Tag.%d.Value", count)] = value
				count++
			}

			wait := incrementalWait(2*time.Second, 1*time.Second)
			err := resource.Retry(10*time.Minute, func() *resource.RetryError {
				response, err := client.RpcPost("cddc", "2020-03-20", action, nil, request, false)
				if err != nil {
					if NeedRetry(err) {
						wait()
						return resource.RetryableError(err)

					}
					return resource.NonRetryableError(err)
				}
				addDebug(action, response, request)
				return nil
			})
			if err != nil {
				return WrapErrorf(err, DefaultErrorMsg, d.Id(), action, AlibabaCloudSdkGoERROR)
			}
		}
		d.SetPartial("tags")
	}
	return nil
}

func (s *CddcService) DescribeCddcDedicatedHost(id string) (object map[string]interface{}, err error) {
	var response map[string]interface{}
	client := s.client
	parts, err := ParseResourceId(id, 2)
	if err != nil {
		err = WrapError(err)
		return
	}
	action := "DescribeDedicatedHostAttribute"
	request := map[string]interface{}{
		"RegionId":             s.client.RegionId,
		"DedicatedHostGroupId": parts[0],
		"DedicatedHostId":      parts[1],
	}
	wait := incrementalWait(3*time.Second, 3*time.Second)
	err = resource.Retry(5*time.Minute, func() *resource.RetryError {
		response, err = client.RpcPost("cddc", "2020-03-20", action, nil, request, true)
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
		if IsExpectedErrors(err, []string{"InvalidDedicatedHostGroup.NotFound"}) {
			return object, WrapErrorf(NotFoundErr("DedicatedHost", id), NotFoundMsg, ProviderERROR, fmt.Sprint(response["RequestId"]))
		}
		return object, WrapErrorf(err, DefaultErrorMsg, id, action, AlibabaCloudSdkGoERROR)
	}
	v, err := jsonpath.Get("$", response)
	if err != nil {
		return object, WrapErrorf(err, FailedGetAttributeMsg, id, "$", response)
	}
	if fmt.Sprint(v.(map[string]interface{})["DedicatedHostId"]) != parts[1] {
		return object, WrapErrorf(NotFoundErr("CDDC", id), NotFoundWithResponse, response)
	}
	object = v.(map[string]interface{})
	return object, nil
}

func (s *CddcService) CddcDedicatedHostStateRefreshFunc(id string, failStates []string) resource.StateRefreshFunc {
	return func() (interface{}, string, error) {
		object, err := s.DescribeCddcDedicatedHost(id)
		if err != nil {
			if NotFoundError(err) {
				// Set this to nil as if we didn't find anything.
				return nil, "", nil
			}
			return nil, "", WrapError(err)
		}

		for _, failState := range failStates {
			if fmt.Sprint(object["HostStatus"]) == failState {
				return object, fmt.Sprint(object["HostStatus"]), WrapError(Error(FailedToReachTargetStatus, fmt.Sprint(object["HostStatus"])))
			}
		}
		return object, fmt.Sprint(object["HostStatus"]), nil
	}
}
func (s *CddcService) DescribeCddcDedicatedHostAccount(id string) (object map[string]interface{}, err error) {
	var response map[string]interface{}
	client := s.client
	action := "DescribeDedicatedHosts"
	parts, err := ParseResourceId(id, 2)
	if err != nil {
		err = WrapError(err)
		return
	}
	request := map[string]interface{}{
		"RegionId":        s.client.RegionId,
		"DedicatedHostId": parts[0],
		"PageSize":        PageSizeMedium,
	}
	wait := incrementalWait(3*time.Second, 3*time.Second)
	err = resource.Retry(5*time.Minute, func() *resource.RetryError {
		response, err = client.RpcPost("cddc", "2020-03-20", action, nil, request, true)
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
		return object, WrapErrorf(err, DefaultErrorMsg, id, action, AlibabaCloudSdkGoERROR)
	}
	v, err := jsonpath.Get("$.DedicatedHosts.DedicatedHosts", response)
	if err != nil {
		return object, WrapErrorf(err, FailedGetAttributeMsg, id, "$.DedicatedHosts.DedicatedHosts", response)
	}
	if len(v.([]interface{})) < 1 {
		return object, WrapErrorf(NotFoundErr("CDDC", id), NotFoundWithResponse, response)
	}
	for _, v := range v.([]interface{}) {
		if fmt.Sprint(v.(map[string]interface{})["AccountName"]) == parts[1] {
			return v.(map[string]interface{}), nil
		}
	}
	return object, WrapErrorf(NotFoundErr("CDDC", id), NotFoundWithResponse, response)
}

func (s *CddcService) CddcDedicatedHostAccountStateRefreshFunc(d *schema.ResourceData, failStates []string) resource.StateRefreshFunc {
	return func() (interface{}, string, error) {
		object, err := s.DescribeCddcDedicatedHostAccount(d.Id())
		if err != nil {
			if NotFoundError(err) {
				// Set this to nil as if we didn't find anything.
				return nil, "", nil
			}
			return nil, "", WrapError(err)
		}

		for _, failState := range failStates {
			if fmt.Sprint(object["AccountName"]) == failState {
				return object, fmt.Sprint(object["AccountName"]), WrapError(Error(FailedToReachTargetStatus, fmt.Sprint(object["AccountName"])))
			}
		}
		return object, fmt.Sprint(object["AccountName"]), nil
	}
}
