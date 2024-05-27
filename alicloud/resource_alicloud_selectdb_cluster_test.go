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
	resource.AddTestSweepers("alicloud_selectdb_db_cluster", &resource.Sweeper{
		Name: "alicloud_selectdb_db_cluster",
		F:    testSweepSelectDBClusterAndInstance,
	})
}

func testSweepSelectDBClusterAndInstance(region string) error {
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

		log.Printf("[Info] Delete SelectDB DBInstance and cluster: %s", id)
		_, err := selectDBService.DeleteSelectDBInstance(id["DBInstanceId"].(string))
		if err != nil {
			return WrapErrorf(err, DataDefaultErrorMsg, "alicloud_selectdb_db_instances_clusters", AlibabaCloudSdkGoERROR)
		}
	}
	return nil
}

func TestAccAliCloudSelectDBCluster_basic_info(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_selectdb_db_cluster.default"
	ra := resourceAttrInit(resourceId, AliCloudSelectDBDBClusterMap0)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return &SelectDBService{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}, "DescribeSelectDBDbCluster")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := acctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tf-testacc%sselectdbcluster%d", defaultRegionToTest, rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, AliCloudSelectDBClusterBasicDependence0)
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
					"db_instance_class":      "selectdb.xlarge",
					"cache_size":             "200",
					"payment_type":           "PayAsYouGo",
					"db_cluster_description": name,
					"db_instance_id":         "${alicloud_selectdb_db_instance.default.db_instance_id}",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"db_instance_class":      "selectdb.xlarge",
						"cache_size":             "200",
						"payment_type":           "PayAsYouGo",
						"db_cluster_description": CHECKSET,
						"db_instance_id":         CHECKSET,
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"db_cluster_description": name + "_desc",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"db_cluster_description": name + "_desc",
					}),
				),
			},
			{
				ResourceName:      resourceId,
				ImportState:       true,
				ImportStateVerify: true, ImportStateVerifyIgnore: []string{"db_instance_id"},
			},
		},
	})
}

func TestAccAliCloudSelectDBCluster_config_list(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_selectdb_db_cluster.default"
	ra := resourceAttrInit(resourceId, AliCloudSelectDBDBClusterMap0)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return &SelectDBService{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}, "DescribeSelectDBDbCluster")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := acctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tf-testacc%sselectdbdbcluster%d", defaultRegionToTest, rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, AliCloudSelectDBClusterBasicDependence0)
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
					"db_instance_class":      "selectdb.xlarge",
					"cache_size":             "200",
					"payment_type":           "PayAsYouGo",
					"db_cluster_description": name,
					"db_instance_id":         "${alicloud_selectdb_db_instance.default.db_instance_id}",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"db_instance_class":      "selectdb.xlarge",
						"cache_size":             "200",
						"payment_type":           "PayAsYouGo",
						"db_cluster_description": name,
						"db_instance_id":         CHECKSET,
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"params": []map[string]interface{}{
						{
							"name":  "es_http_timeout_ms",
							"value": "6000",
						},
						{
							"name":  "sys_log_roll_num",
							"value": "12",
						},
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"params.#": "2",
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

func TestAccAliCloudSelectDBDBCluster_status(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_selectdb_db_cluster.default"
	ra := resourceAttrInit(resourceId, AliCloudSelectDBDBClusterMap0)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return &SelectDBService{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}, "DescribeSelectDBDbCluster")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := acctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tf-testacc%sselectdbdbcluster%d", defaultRegionToTest, rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, AliCloudSelectDBClusterBasicDependence0)
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
					"db_instance_class":      "selectdb.xlarge",
					"cache_size":             "200",
					"payment_type":           "PayAsYouGo",
					"db_cluster_description": name,
					"db_instance_id":         "${alicloud_selectdb_db_instance.default.db_instance_id}",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"db_instance_class":      "selectdb.xlarge",
						"cache_size":             "200",
						"payment_type":           "PayAsYouGo",
						"db_cluster_description": name,
						"db_instance_id":         CHECKSET,
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"status": "STOPPING",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"status": "STOPPED",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"status": "STARTING",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"status": "ACTIVATION",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"status": "RESTART",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"status": "ACTIVATION",
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

func TestAccAliCloudSelectDBDBCluster_modify(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_selectdb_db_cluster.default"
	ra := resourceAttrInit(resourceId, AliCloudSelectDBDBClusterMap0)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return &SelectDBService{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}, "DescribeSelectDBDbCluster")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := acctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tf-testacc%sselectdbdbcluster%d", defaultRegionToTest, rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, AliCloudSelectDBClusterBasicDependence0)
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
					"db_instance_class":      "selectdb.xlarge",
					"cache_size":             "200",
					"payment_type":           "PayAsYouGo",
					"db_cluster_description": name,
					"db_instance_id":         "${alicloud_selectdb_db_instance.default.db_instance_id}",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"db_instance_class":      "selectdb.xlarge",
						"cache_size":             "200",
						"payment_type":           "PayAsYouGo",
						"db_cluster_description": name,
						"db_instance_id":         CHECKSET,
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"cache_size": "400",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"cache_size": "400",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"db_instance_class": "selectdb.2xlarge",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"db_instance_class": "selectdb.2xlarge",
						"cache_size":        "800",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"db_instance_class": "selectdb.8xlarge",
					"cache_size":        "4000",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"db_instance_class": "selectdb.8xlarge",
						"cache_size":        "4000",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"db_instance_class": "selectdb.xlarge",
					"cache_size":        "1000",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"db_instance_class": "selectdb.xlarge",
						"cache_size":        "1000",
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

var AliCloudSelectDBDBClusterMap0 = map[string]string{
	"db_instance_id":         CHECKSET,
	"db_cluster_description": CHECKSET,
	"payment_type":           CHECKSET,
	"period":                 CHECKSET,
	"period_time":            CHECKSET,
	"db_cluster_class":       CHECKSET,
	"cache_size":             CHECKSET,
}

func AliCloudSelectDBClusterBasicDependence0(name string) string {
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
	resource "alicloud_vswitch" "this" {
		count        = length(data.alicloud_vswitches.default.ids) > 0 ? 0 : 1
		vswitch_name = var.name
		vpc_id       = data.alicloud_vpcs.default.ids.0
		zone_id      = data.alicloud_zones.default.ids.0
		cidr_block   = cidrsubnet(data.alicloud_vpcs.default.vpcs.0.cidr_block, 8, 4)
	}
	resource "alicloud_selectdb_db_instance" "default" {
		db_instance_class        = "selectdb.xlarge",
		db_instance_description  = var.name
		cache_size               = "200",
		payment_type             = "PayAsYouGo",
		engine_version           = "3.0",
		vpc_id                   = "${data.alicloud_vpcs.default.ids.0}"
		zone_id                  = "${data.alicloud_zones.default.ids.0}"
		vswitch_id               = "${data.alicloud_vswitches.default.ids.0}",
	}
`, name)
}
