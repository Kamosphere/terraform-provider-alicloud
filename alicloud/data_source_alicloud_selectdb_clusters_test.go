package alicloud

import (
	"fmt"
	"strings"
	"testing"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"

	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"
)

func TestAccAlicloudSelectDBDbClusterDataSource(t *testing.T) {
	rand := acctest.RandInt()
	dbClusterIdsConf := dataSourceTestAccConfig{
		existConfig: testAccCheckAlicloudSelectDBDbClusterDataSourceName(rand, map[string]string{
			"db_cluster_ids": `["${alicloud_selectdb_db_cluster.default.id}"]`,
		}),
		fakeConfig: testAccCheckAlicloudSelectDBDbClusterDataSourceName(rand, map[string]string{
			"db_cluster_ids": `["${alicloud_selectdb_db_cluster.default.id}_fake"]`,
		}),
	}
	dbInstanceIdConf := dataSourceTestAccConfig{
		existConfig: testAccCheckAlicloudSelectDBDbClusterDataSourceName(rand, map[string]string{
			"db_instance_id": `"${alicloud_selectdb_db_cluster.default.db_instance_id}"`,
		}),
		fakeConfig: testAccCheckAlicloudSelectDBDbClusterDataSourceName(rand, map[string]string{
			"db_instance_id": `"${alicloud_selectdb_db_cluster.default.db_instance_id}_fake"`,
		}),
	}

	allConf := dataSourceTestAccConfig{
		existConfig: testAccCheckAlicloudSelectDBDbClusterDataSourceName(rand, map[string]string{
			"db_cluster_ids": `["${alicloud_selectdb_db_cluster.default.id}"]`,
			"db_instance_id": `"${alicloud_selectdb_db_cluster.default.db_instance_id}"`,
		}),
		fakeConfig: testAccCheckAlicloudSelectDBDbClusterDataSourceName(rand, map[string]string{
			"db_cluster_ids": `["${alicloud_selectdb_db_cluster.default.id}_fake"]`,
			"db_instance_id": `"${alicloud_selectdb_db_cluster.default.db_instance_id}_fake"`,
		}),
	}
	var existAlicloudSelectDBDbClusterDataSourceNameMapFunc = func(rand int) map[string]string {
		return map[string]string{
			"ids.#":                                       "2",
			"clusters.#":                                  "2",
			"clusters.0.payment_type":                     "PayAsYouGo",
			"clusters.0.ali_uid":                          CHECKSET,
			"clusters.0.bid":                              CHECKSET,
			"clusters.0.commodity_code":                   CHECKSET,
			"clusters.0.db_cluster_id":                    CHECKSET,
			"clusters.0.db_instance_id":                   CHECKSET,
			"clusters.0.db_cluster_class":                 "selectdb.xlarge",
			"clusters.0.db_cluster_description":           "new_cluster",
			"clusters.0.engine":                           "selectdb",
			"clusters.0.engine_version":                   CHECKSET,
			"clusters.0.create_time":                      CHECKSET,
			"clusters.0.status":                           "ACTIVATION",
			"clusters.0.cpu":                              "4",
			"clusters.0.memory":                           "32",
			"clusters.0.cache_size":                       "200",
			"clusters.0.region_id":                        "default",
			"clusters.0.vswitch_id":                       CHECKSET,
			"clusters.0.vpc_id":                           CHECKSET,
			"clusters.0.zone_id":                          CHECKSET,
			"clusters.0.params.#":                         "1",
			"clusters.0.params.0.name":                    "test",
			"clusters.0.params.0.optional":                "1",
			"clusters.0.params.0.comment":                 "test config",
			"clusters.0.params.0.value":                   "true",
			"clusters.0.params.0.param_category":          "test",
			"clusters.0.params.0.default_value":           "false",
			"clusters.0.params.0.is_dynamic":              "1",
			"clusters.0.params.0.is_user_modifiable":      "1",
			"clusters.0.param_change_logs.#":              "1",
			"clusters.0.param_change_logs.0.name":         "test",
			"clusters.0.param_change_logs.0.old_value":    "false",
			"clusters.0.param_change_logs.0.new_value":    "true",
			"clusters.0.param_change_logs.0.gmt_created":  "2000-01-01",
			"clusters.0.param_change_logs.0.gmt_modified": "2000-01-01",
			"clusters.0.param_change_logs.0.config_id":    "1",
			"clusters.0.param_change_logs.0.is_applied":   "1",
			"clusters.1.payment_type":                     "PayAsYouGo",
			"clusters.1.ali_uid":                          CHECKSET,
			"clusters.1.bid":                              CHECKSET,
			"clusters.1.commodity_code":                   CHECKSET,
			"clusters.1.db_cluster_id":                    CHECKSET,
			"clusters.1.db_instance_id":                   CHECKSET,
			"clusters.1.db_cluster_class":                 "selectdb.2xlarge",
			"clusters.1.db_cluster_description":           fmt.Sprintf("tf-testAccSelectDBDbCluster-%d", rand),
			"clusters.1.engine":                           "selectdb",
			"clusters.1.engine_version":                   CHECKSET,
			"clusters.1.create_time":                      CHECKSET,
			"clusters.1.status":                           "ACTIVATION",
			"clusters.1.cpu":                              "8",
			"clusters.1.memory":                           "64",
			"clusters.1.cache_size":                       "400",
			"clusters.1.region_id":                        "default",
			"clusters.1.vswitch_id":                       CHECKSET,
			"clusters.1.vpc_id":                           CHECKSET,
			"clusters.1.zone_id":                          CHECKSET,
			"clusters.1.params.#":                         "1",
			"clusters.1.params.0.name":                    "test",
			"clusters.1.params.0.optional":                "1",
			"clusters.1.params.0.comment":                 "test config",
			"clusters.1.params.0.value":                   "true",
			"clusters.1.params.0.param_category":          "test",
			"clusters.1.params.0.default_value":           "false",
			"clusters.1.params.0.is_dynamic":              "1",
			"clusters.1.params.0.is_user_modifiable":      "1",
			"clusters.1.param_change_logs.#":              "1",
			"clusters.1.param_change_logs.0.name":         "test",
			"clusters.1.param_change_logs.0.old_value":    "false",
			"clusters.1.param_change_logs.0.new_value":    "true",
			"clusters.1.param_change_logs.0.gmt_created":  "2000-01-01",
			"clusters.1.param_change_logs.0.gmt_modified": "2000-01-01",
			"clusters.1.param_change_logs.0.config_id":    "1",
			"clusters.1.param_change_logs.0.is_applied":   "1",
		}
	}
	var fakeAlicloudSelectDBDbClusterDataSourceNameMapFunc = func(rand int) map[string]string {
		return map[string]string{
			"ids.#":   "0",
			"names.#": "0",
		}
	}
	var alicloudSelectDBDbClusterCheckInfo = dataSourceAttr{
		resourceId:   "data.alicloud_selectdb_db_clusters.default",
		existMapFunc: existAlicloudSelectDBDbClusterDataSourceNameMapFunc,
		fakeMapFunc:  fakeAlicloudSelectDBDbClusterDataSourceNameMapFunc,
	}
	preCheck := func() {
		testAccPreCheckWithRegions(t, true, connectivity.SelectDBSupportRegions)
	}

	alicloudSelectDBDbClusterCheckInfo.dataSourceTestCheckWithPreCheck(t, rand, preCheck, dbClusterIdsConf, dbInstanceIdConf, allConf)
}

func testAccCheckAlicloudSelectDBDbClusterDataSourceName(rand int, attrMap map[string]string) string {
	var pairs []string
	for k, v := range attrMap {
		pairs = append(pairs, k+" = "+v)
	}

	config := fmt.Sprintf(`
variable "name" {	
	default = "tf-testAccSelectDBDbCluster-%d"
}
data "alicloud_zones" "default" {
	available_resource_creation = "VSwitch"
}
data "alicloud_vpcs" "default" {
	name_regex = "^default-NODELETING$"
  }
data "alicloud_vswitches" "default" {
	vpc_id = "${data.alicloud_vpcs.default.ids.0}"
	zone_id = "${data.alicloud_zones.default.ids.0}"
}
resource "alicloud_selectdb_db_instance" "default" {
	db_instance_class  =  "selectdb.xlarge",
	cache_size         =  "200",
	payment_type       =  "PayAsYouGo",
	engine_version     =  "3.0",
	vpc_id             =  "${data.alicloud_vpcs.default.ids.0}"
    zone_id            =  "${data.alicloud_zones.default.ids.0}"
	vswitch_id         =  "${data.alicloud_vswitches.default.ids.0}",
  }
resource "alicloud_selectdb_db_cluster" "default" {
	db_instance_id          =  "${alicloud_selectdb_db_instance.default.id}"
	db_cluster_description  =  var.name
	db_cluster_class        =  "selectdb.2xlarge"
	cache_size              =  "400",
	payment_type            =  "PayAsYouGo",
}

data "alicloud_selectdb_db_clusters" "default" {	
	%s
}
`, rand, strings.Join(pairs, " \n "))
	return config
}
