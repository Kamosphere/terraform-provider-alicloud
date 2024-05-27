package alicloud

import (
	"fmt"
	"strings"
	"testing"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"

	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"
)

func TestAccAlicloudSelectDBDbInstanceDataSource(t *testing.T) {
	rand := acctest.RandInt()
	dbInstanceIdsConf := dataSourceTestAccConfig{
		existConfig: testAccCheckAlicloudSelectDBDbInstanceDataSourceName(rand, map[string]string{
			"db_instance_ids": `["${alicloud_selectdb_db_instance.default.id}"]`,
		}),
		fakeConfig: testAccCheckAlicloudSelectDBDbInstanceDataSourceName(rand, map[string]string{
			"db_instance_ids": `["${alicloud_selectdb_db_instance.default.id}_fake"]`,
		}),
	}
	var existAlicloudSelectDBDbInstanceDataSourceNameMapFunc = func(rand int) map[string]string {
		return map[string]string{
			"ids.#":                                       "1",
			"instances.#":                                  "1",
			"instances.0.payment_type":                     "PayAsYouGo",
			"instances.0.ali_uid":                          CHECKSET,
			"instances.0.bid":                              CHECKSET,
			"instances.0.commodity_code":                   CHECKSET,
			"instances.0.db_cluster_id":                    CHECKSET,
			"instances.0.db_instance_id":                   CHECKSET,
			"instances.0.db_cluster_class":                 "selectdb.xlarge",
			"instances.0.db_cluster_description":           "new_cluster",
			"instances.0.engine":                           "selectdb",
			"instances.0.engine_version":                   CHECKSET,
			"instances.0.create_time":                      CHECKSET,
			"instances.0.status":                           "ACTIVATION",
			"instances.0.cpu":                              "4",
			"instances.0.memory":                           "32",
			"instances.0.cache_size":                       "200",
			"instances.0.region_id":                        "default",
			"instances.0.vswitch_id":                       CHECKSET,
			"instances.0.vpc_id":                           CHECKSET,
			"instances.0.zone_id":                          CHECKSET,
			"instances.0.params.#":                         "1",
			"instances.0.params.0.name":                    "test",
			"instances.0.params.0.optional":                "1",
			"instances.0.params.0.comment":                 "test config",
			"instances.0.params.0.value":                   "true",
			"instances.0.params.0.param_category":          "test",
			"instances.0.params.0.default_value":           "false",
			"instances.0.params.0.is_dynamic":              "1",
			"instances.0.params.0.is_user_modifiable":      "1",
			"instances.0.param_change_logs.#":              "1",
			"instances.0.param_change_logs.0.name":         "test",
			"instances.0.param_change_logs.0.old_value":    "false",
			"instances.0.param_change_logs.0.new_value":    "true",
			"instances.0.param_change_logs.0.gmt_created":  "2000-01-01",
			"instances.0.param_change_logs.0.gmt_modified": "2000-01-01",
			"instances.0.param_change_logs.0.config_id":    "1",
			"instances.0.param_change_logs.0.is_applied":   "1",
		}
	}
	var fakeAlicloudSelectDBDbInstanceDataSourceNameMapFunc = func(rand int) map[string]string {
		return map[string]string{
			"ids.#":   "0",
			"names.#": "0",
		}
	}
	var alicloudSelectDBDbInstanceCheckInfo = dataSourceAttr{
		resourceId:   "data.alicloud_selectdb_db_instances.default",
		existMapFunc: existAlicloudSelectDBDbInstanceDataSourceNameMapFunc,
		fakeMapFunc:  fakeAlicloudSelectDBDbInstanceDataSourceNameMapFunc,
	}
	preCheck := func() {
		testAccPreCheckWithRegions(t, true, connectivity.SelectDBSupportRegions)
	}

	alicloudSelectDBDbInstanceCheckInfo.dataSourceTestCheckWithPreCheck(t, rand, preCheck, dbInstanceIdsConf)
}

func testAccCheckAlicloudSelectDBDbInstanceDataSourceName(rand int, attrMap map[string]string) string {
	var pairs []string
	for k, v := range attrMap {
		pairs = append(pairs, k+" = "+v)
	}

	config := fmt.Sprintf(`
variable "name" {	
    default = "tf-testAccSelectDBDbInstance-%d"
}
data "alicloud_zones" "default" {
    available_resource_creation = "VSwitch"
}
data "alicloud_vpcs" "default" {
    name_regex = "^default-NODELETING$"
  }
data "alicloud_vswitches" "default" {
    vpc_id  = "${data.alicloud_vpcs.default.ids.0}"
    zone_id = "${data.alicloud_zones.default.ids.0}"
}
resource "alicloud_selectdb_db_instance" "default" {
	db_instance_class        =  "selectdb.xlarge",
	db_instance_description  =  var.name
	cache_size               =  "200",
	payment_type             =  "PayAsYouGo",
	engine_version           =  "3.0",
	vpc_id                   =  "${data.alicloud_vpcs.default.ids.0}"
    zone_id                  =  "${data.alicloud_zones.default.ids.0}"
	vswitch_id               =  "${data.alicloud_vswitches.default.ids.0}",
  }

data "alicloud_selectdb_db_instances" "default" {	
	%s
}
`, rand, strings.Join(pairs, " \n "))
	return config
}
