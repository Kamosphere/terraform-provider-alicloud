package alicloud

import (
	"fmt"
	"log"
	"strings"
	"testing"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
)

func init() {
	resource.AddTestSweepers("alicloud_selectdb_db_instance", &resource.Sweeper{
		Name: "alicloud_selectdb_db_instance",
		F:    testSweepSelectDBDbInstance,
	})
}

func testSweepSelectDBDbInstance(region string) error {
	rawClient, err := sharedClientForRegion(region)
	if err != nil {
		return fmt.Errorf("error getting AliCloud client: %s", err)
	}
	client := rawClient.(*connectivity.AliyunClient)

	prefixes := []string{
		"tf-testAcc",
		"tf_testAcc",
	}
	selectDBService := SelectDBService{client}
	instanceResp, err := selectDBService.DescribeSelectDBDbInstances("")
	if err != nil {
		return WrapErrorf(err, DataDefaultErrorMsg, "alicloud_selectdb_db_instances", AlibabaCloudSdkGoERROR)
	}

	var objects []map[string]interface{}

	for _, item := range instanceResp {
		name := item["Description"].(string)
		skip := true
		if !sweepAll() {
			for _, prefix := range prefixes {
				if strings.HasPrefix(name, prefix) {
					skip = false
					break
				}
			}
			if skip {
				log.Printf("[INFO] Skipping DBinstance sweep: %s ", name)
				continue
			}
		}
		objects = append(objects, item)
	}

	for _, id := range objects {

		log.Printf("[Info] Delete SelectDB DBInstance: %s", id)
		_, err := selectDBService.DeleteSelectDBInstance(id["DBInstanceId"].(string))
		if err != nil {
			return WrapErrorf(err, DataDefaultErrorMsg, "alicloud_selectdb_db_instances", AlibabaCloudSdkGoERROR)
		}
	}
	return nil
}

func TestAccAliCloudSelectDBDBInstance_basic_info(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_selectdb_db_instance.default"
	ra := resourceAttrInit(resourceId, AliCloudSelectDBDBInstanceMap0)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return &SelectDBService{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}, "DescribeSelectDBDbInstance")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := acctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tf-testacc%sselectdbdbinstance%d", defaultRegionToTest, rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, AliCloudSelectDBDBInstanceBasicDependence0)
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
			testAccPreCheckWithRegions(t, true, connectivity.SelectDBSupportRegions)
		},
		IDRefreshName: resourceId,
		Providers:     testAccProviders,
		CheckDestroy:  rac.checkResourceDestroy(),
		Steps: []resource.TestStep{
			{

				Config: testAccConfig(map[string]interface{}{
					"db_instance_class": "selectdb.xlarge",
					"cache_size":        "200",
					"payment_type":      "PayAsYouGo",
					"engine_version":    "3.0",
					"zone_id":           "${alicloud_vswitch.default.zone_id}",
					"vpc_id":            "${alicloud_vswitch.default.vpc_id}",
					"vswitch_id":        "${alicloud_vswitch.default.id}",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"db_instance_class": "selectdb.xlarge",
						"cache_size":        "200",
						"payment_type":      "PayAsYouGo",
						"engine_version":    "3.0",
						"zone_id":           CHECKSET,
						"vpc_id":            CHECKSET,
						"vswitch_id":        CHECKSET,
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"db_instance_description": name + "_desc",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"db_instance_description": name + "_desc",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"maintain_start_time": "02:00",
					"maintain_end_time":   "03:00",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"maintain_start_time": "02:00",
						"maintain_end_time":   "03:00",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"db_instance_description": name + "_updateAll",
					"maintain_start_time":     "01:00",
					"maintain_end_time":       "02:00",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"db_instance_description": name + "_updateAll",
						"maintain_start_time":     "01:00",
						"maintain_end_time":       "02:00",
					}),
				),
			},
			{
				ResourceName:      resourceId,
				ImportState:       true,
				ImportStateVerify: true, ImportStateVerifyIgnore: []string{"db_cluster_class", "cache_size", "engine_version"},
			},
		},
	})
}

func TestAccAliCloudSelectDBDBInstance_basic_prepaid(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_selectdb_db_instance.default"
	ra := resourceAttrInit(resourceId, AliCloudSelectDBDBInstanceMap0)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return &SelectDBService{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}, "DescribeSelectDBDbInstance")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := acctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tf-testacc%sselectdbdbinstance%d", defaultRegionToTest, rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, AliCloudSelectDBDBInstanceBasicDependence0)
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
			testAccPreCheckWithRegions(t, true, connectivity.SelectDBSupportRegions)
		},
		IDRefreshName: resourceId,
		Providers:     testAccProviders,
		CheckDestroy:  rac.checkResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccConfig(map[string]interface{}{
					"db_instance_class": "selectdb.xlarge",
					"cache_size":        "200",
					"payment_type":      "Subscription",
					"period":            "Month",
					"period_time":       "1",
					"engine_version":    "3.0",
					"zone_id":           "${alicloud_vswitch.default.zone_id}",
					"vpc_id":            "${alicloud_vswitch.default.vpc_id}",
					"vswitch_id":        "${alicloud_vswitch.default.id}",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"db_instance_class": "selectdb.xlarge",
						"cache_size":        "200",
						"payment_type":      "Subscription",
						"period":            "Month",
						"period_time":       "1",
						"engine_version":    "3.0",
						"zone_id":           CHECKSET,
						"vpc_id":            CHECKSET,
						"vswitch_id":        CHECKSET,
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"db_instance_description": name + "_desc",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"db_instance_description": name + "_desc",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"maintain_start_time": "02:00",
					"maintain_end_time":   "03:00",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"maintain_start_time": "02:00",
						"maintain_end_time":   "03:00",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"db_instance_description": name + "_updateAll",
					"maintain_start_time":     "01:00",
					"maintain_end_time":       "02:00",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"db_instance_description": name + "_updateAll",
						"maintain_start_time":     "01:00",
						"maintain_end_time":       "02:00",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"payment_type": "PayAsYouGo",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"payment_type": "PayAsYouGo",
					}),
				),
			},
			{
				ResourceName:      resourceId,
				ImportState:       true,
				ImportStateVerify: true, ImportStateVerifyIgnore: []string{"db_cluster_class", "cache_size", "engine_version"},
			},
		},
	})
}

func TestAccAliCloudSelectDBDBInstance_security_list(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_selectdb_db_instance.default"
	ra := resourceAttrInit(resourceId, AliCloudSelectDBDBInstanceMap0)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return &SelectDBService{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}, "DescribeSelectDBDbInstance")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := acctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tf-testacc%sselectdbdbinstance%d", defaultRegionToTest, rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, AliCloudSelectDBDBInstanceBasicDependence0)
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
			testAccPreCheckWithRegions(t, true, connectivity.SelectDBSupportRegions)
		},
		IDRefreshName: resourceId,
		Providers:     testAccProviders,
		CheckDestroy:  rac.checkResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccConfig(map[string]interface{}{
					"db_instance_class": "selectdb.xlarge",
					"cache_size":        "200",
					"payment_type":      "PayAsYouGo",
					"engine_version":    "3.0",
					"zone_id":           "${alicloud_vswitch.default.zone_id}",
					"vpc_id":            "${alicloud_vswitch.default.vpc_id}",
					"vswitch_id":        "${alicloud_vswitch.default.id}",
					"security_ip_lists": []map[string]interface{}{
						{
							"group_name":       "test1",
							"security_ip_list": "192.168.0.1",
						},
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"db_instance_class":   "selectdb.xlarge",
						"cache_size":          "200",
						"payment_type":        "PayAsYouGo",
						"engine_version":      "3.0",
						"zone_id":             CHECKSET,
						"vpc_id":              CHECKSET,
						"vswitch_id":          CHECKSET,
						"security_ip_lists.#": "1",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"security_ip_lists": []map[string]interface{}{
						{
							"group_name":       "test2",
							"security_ip_list": "192.168.0.3",
						},
						{
							"group_name":       "test1",
							"security_ip_list": "192.168.0.2",
						},
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"security_ip_lists.#": "2",
					}),
				),
			},
			{
				ResourceName:            resourceId,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"db_cluster_class", "db_node_group_count", "db_cluster_version"},
			},
		},
	})
}

func TestAccAliCloudSelectDBDBInstance_upgrade(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_selectdb_db_instance.default"
	ra := resourceAttrInit(resourceId, AliCloudSelectDBDBInstanceMap0)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return &SelectDBService{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}, "DescribeSelectDBDbInstance")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := acctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tf-testacc%sselectdbdbinstance%d", defaultRegionToTest, rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, AliCloudSelectDBDBInstanceBasicDependence0)
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
			testAccPreCheckWithRegions(t, true, connectivity.SelectDBSupportRegions)
		},
		IDRefreshName: resourceId,
		Providers:     testAccProviders,
		CheckDestroy:  rac.checkResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccConfig(map[string]interface{}{
					"db_instance_class": "selectdb.xlarge",
					"cache_size":        "200",
					"payment_type":      "PayAsYouGo",
					"engine_version":    "2.4",
					"zone_id":           "${alicloud_vswitch.default.zone_id}",
					"vpc_id":            "${alicloud_vswitch.default.vpc_id}",
					"vswitch_id":        "${alicloud_vswitch.default.id}",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"db_instance_class": "selectdb.xlarge",
						"cache_size":        "200",
						"payment_type":      "PayAsYouGo",
						"engine_version":    "2.4",
						"zone_id":           CHECKSET,
						"vpc_id":            CHECKSET,
						"vswitch_id":        CHECKSET,
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"engine_version": "3.0",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"engine_version": "3.0",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"db_instance_description": name + "_updateAll",
					"maintain_start_time":     "01:00",
					"maintain_end_time":       "02:00",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"db_instance_description": name + "_updateAll",
						"maintain_start_time":     "01:00",
						"maintain_end_time":       "02:00",
					}),
				),
			},
			{
				ResourceName:      resourceId,
				ImportState:       true,
				ImportStateVerify: true, ImportStateVerifyIgnore: []string{"db_cluster_class", "cache_size", "engine_version"},
			},
		},
	})
}

func TestAccAliCloudSelectDBDBInstance_network(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_selectdb_db_instance.default"
	ra := resourceAttrInit(resourceId, AliCloudSelectDBDBInstanceMap0)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return &SelectDBService{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}, "DescribeSelectDBDbInstance")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := acctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tf-testacc%sselectdbdbinstance%d", defaultRegionToTest, rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, AliCloudSelectDBDBInstanceBasicDependence0)
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
			testAccPreCheckWithRegions(t, true, connectivity.SelectDBSupportRegions)
		},
		IDRefreshName: resourceId,
		Providers:     testAccProviders,
		CheckDestroy:  rac.checkResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccConfig(map[string]interface{}{
					"db_instance_class": "selectdb.xlarge",
					"cache_size":        "200",
					"payment_type":      "PayAsYouGo",
					"engine_version":    "3.0",
					"zone_id":           "${alicloud_vswitch.default.zone_id}",
					"vpc_id":            "${alicloud_vswitch.default.vpc_id}",
					"vswitch_id":        "${alicloud_vswitch.default.id}",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"db_instance_class":     "selectdb.xlarge",
						"cache_size":            "200",
						"payment_type":          "PayAsYouGo",
						"engine_version":        "3.0",
						"zone_id":               CHECKSET,
						"vpc_id":                CHECKSET,
						"vswitch_id":            CHECKSET,
						"instance_net_infos.#":  "1",
						"enable_public_network": "false",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"enable_public_network": "true",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"enable_public_network": "true",
						"instance_net_infos.#":  "3",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"enable_public_network": "false",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"enable_public_network": "false",
						"instance_net_infos.#":  "1",
					}),
				),
			},
			{
				ResourceName:      resourceId,
				ImportState:       true,
				ImportStateVerify: true, ImportStateVerifyIgnore: []string{"db_cluster_class", "cache_size", "engine_version"},
			},
		},
	})
}

var AliCloudSelectDBDBInstanceMap0 = map[string]string{
	"db_instance_class": CHECKSET,
	"cache_size":        CHECKSET,
	"payment_type":      CHECKSET,
	"period":            CHECKSET,
	"period_time":       CHECKSET,
	"engine_version":    CHECKSET,
	"zone_id":           CHECKSET,
	"vpc_id":            CHECKSET,
	"vswitch_id":        CHECKSET,
}

func AliCloudSelectDBDBInstanceBasicDependence0(name string) string {
	return fmt.Sprintf(`
	variable "name" {
  		default = "%s"
	}

	data "alicloud_zones" "default" {
		available_resource_creation = "VSwitch"
	}
	data "alicloud_vpcs" "default" {
		name_regex = "^default-NODELETING$"
	}
	data "alicloud_vswitches" "default" {
		vpc_id  = data.alicloud_vpcs.default.ids.0
	}
	resource "alicloud_vswitch" "default" {
		count        = length(data.alicloud_vswitches.default.ids) > 0 ? 0 : 1
		vswitch_name = var.name
		vpc_id       = data.alicloud_vpcs.default.ids.0
		zone_id      = data.alicloud_zones.default.ids.0
		cidr_block   = cidrsubnet(data.alicloud_vpcs.default.vpcs.0.cidr_block, 8, 4)
	}
`, name)
}
