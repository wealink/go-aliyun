package main

import (
	"encoding/json"
	"fmt"
	"github.com/aliyun/alibaba-cloud-sdk-go/sdk/requests"
	"github.com/aliyun/alibaba-cloud-sdk-go/services/dts"
	"github.com/aliyun/alibaba-cloud-sdk-go/services/rds"
	"os"
)

var dtsclient *dts.Client
var rdsclient *rds.Client
var err error

type CreateResponse struct {
	RequestId  string
	InstanceId string
	Success    string
	JobId      string
}

func Init() {
	var (
		region string
		key    string
		secret string
	)
	region = os.Getenv("REGION")
	key = os.Getenv("ACCESS_KEY")
	secret = os.Getenv("ACCESS_SECRET")
	dtsclient, err = dts.NewClientWithAccessKey(region, key, secret)
	if err != nil {
		fmt.Print(err.Error())
		os.Exit(1)
	}
	rdsclient, err = rds.NewClientWithAccessKey(region, key, secret)
	if err != nil {
		fmt.Print(err.Error())
		os.Exit(1)
	}
}

func create_instance() string {
	var result CreateResponse
	request := requests.NewCommonRequest()
	request.Method = "POST"
	request.Scheme = "https" // https | http
	request.Domain = "dts.aliyuncs.com"
	request.Version = "2020-01-01"
	request.ApiName = "CreateDtsInstance"
	request.QueryParams["RegionId"] = "cn-shanghai"
	request.QueryParams["InstanceClass"] = "medium"
	request.QueryParams["PayType"] = "PostPaid"
	request.QueryParams["SyncArchitecture"] = "oneway"
	request.QueryParams["Quantity"] = "1"
	request.QueryParams["Type"] = "MIGRATION"
	request.QueryParams["SourceRegion"] = "cn-shanghai"
	request.QueryParams["DestinationRegion"] = "cn-shanghai"
	request.QueryParams["SourceEndpointEngineName"] = "MySQL"
	request.QueryParams["DestinationEndpointEngineName"] = "MySQL"
	response, err := dtsclient.ProcessCommonRequest(request)
	if err != nil {
		panic(err)
	}
	json.Unmarshal(response.GetHttpContentBytes(), &result)
	return result.InstanceId
}

func clean_database_pgsql(instanceid, dbname, chartset, accountname string) {
	request := rds.CreateDeleteDatabaseRequest()
	request.Scheme = "https"
	request.DBInstanceId = instanceid
	request.DBName = dbname
	_, err := rdsclient.DeleteDatabase(request)
	if err != nil {
		fmt.Println(err.Error())
	}

	request1 := rds.CreateCreateDatabaseRequest()
	request1.Scheme = "https"
	request1.DBInstanceId = instanceid
	request1.DBName = dbname
	request1.CharacterSetName = chartset
	_, err = rdsclient.CreateDatabase(request1)
	if err != nil {
		fmt.Println(err.Error())
	}

	request2 := rds.CreateGrantAccountPrivilegeRequest()
	request2.Scheme = "https"
	request2.DBInstanceId = instanceid
	request2.AccountName = accountname
	request2.DBName = dbname
	request2.AccountPrivilege = "DBOwner"
	_, err = rdsclient.GrantAccountPrivilege(request2)
	if err != nil {
		fmt.Print(err.Error())
	}
}

func mysql_java_test_mulan_db_v56() {
	//create
	id := create_instance()
	//config && start
	request := requests.NewCommonRequest()
	request.Method = "POST"
	request.Scheme = "https" // https | http
	request.Domain = "dts.aliyuncs.com"
	request.Version = "2020-01-01"
	request.ApiName = "ConfigureMigrationJob"
	request.QueryParams["RegionId"] = "cn-shanghai"
	request.QueryParams["MigrationJobName"] = "mysql_java_test_mulan_db_v56"
	request.QueryParams["SourceEndpoint.InstanceType"] = "RDS"
	request.QueryParams["DestinationEndpoint.InstanceType"] = "RDS"
	request.QueryParams["MigrationMode.StructureIntialization"] = "true"
	request.QueryParams["MigrationMode.DataIntialization"] = "true"
	request.QueryParams["MigrationMode.DataSynchronization"] = "false"
	request.QueryParams["MigrationObject"] = "[{\"DBName\":\"wwc_announcement_production\",\"NewDBName\":\"wwc_announcement_test\",\"AllTable\": true},{\"DBName\":\"wwc_base_production\",\"NewDBName\":\"wwc_base_test\",\"AllTable\": true},{\"DBName\":\"wwc_comment_production\",\"NewDBName\":\"wwc_comment_test\",\"AllTable\": true},{\"DBName\":\"wwc_datasync_production\",\"NewDBName\":\"wwc_datasync_test\",\"AllTable\": true},{\"DBName\":\"wwc_door_production\",\"NewDBName\":\"wwc_door_test\",\"AllTable\": true},{\"DBName\":\"wwc_event_production\",\"NewDBName\":\"wwc_event_test\",\"AllTable\": true},{\"DBName\":\"wwc_face_production\",\"NewDBName\":\"wwc_face_test\",\"AllTable\": true},{\"DBName\":\"wwc_feed_production\",\"NewDBName\":\"wwc_feed_test\",\"AllTable\": true},{\"DBName\":\"wwc_im_production\",\"NewDBName\":\"wwc_im_test\",\"AllTable\": true},{\"DBName\":\"wwc_notification_production\",\"NewDBName\":\"wwc_notification_test\",\"AllTable\": true},{\"DBName\":\"wwc_relation_production\",\"NewDBName\":\"wwc_relation_test\",\"AllTable\": true},{\"DBName\":\"wwc_space_production\",\"NewDBName\":\"wwc_space_test\",\"AllTable\": true},{\"DBName\":\"wwc_user_production\",\"NewDBName\":\"wwc_user_test\",\"AllTable\": true}]"
	request.QueryParams["MigrationJobId"] = id
	request.QueryParams["SourceEndpoint.InstanceID"] = "rm-uf6prcst021us4q3p"
	request.QueryParams["SourceEndpoint.EngineName"] = "MySQL"
	request.QueryParams["SourceEndpoint.Region"] = "cn-shanghai"
	request.QueryParams["SourceEndpoint.UserName"] = "dms"
	request.QueryParams["SourceEndpoint.Password"] = "*v*#lgiqMuEKjheW"
	request.QueryParams["DestinationEndpoint.InstanceID"] = "rm-uf6g6uhcktm12y1ik"
	request.QueryParams["DestinationEndpoint.EngineName"] = "MySQL"
	request.QueryParams["DestinationEndpoint.UserName"] = "dms"
	request.QueryParams["DestinationEndpoint.Password"] = "q0OzU4B^bpnEqkS4"
	request.QueryParams["MigrationReserved"] = "{ 	\"autoStartModulesAfterConfig\": \"all\", 	\"targetTableMode\": 2 }"
	response, err := dtsclient.ProcessCommonRequest(request)
	if err != nil {
		panic(err)
	}
	fmt.Println(response.GetHttpStatus())
}

func mysql_java_test_mulan_db_v57() {
	//create
	id := create_instance()

	//config && start
	request := requests.NewCommonRequest()
	request.Method = "POST"
	request.Scheme = "https" // https | http
	request.Domain = "dts.aliyuncs.com"
	request.Version = "2020-01-01"
	request.ApiName = "ConfigureMigrationJob"
	request.QueryParams["RegionId"] = "cn-shanghai"
	request.QueryParams["MigrationJobName"] = "mysql_java_test_mulan_db_v57"
	request.QueryParams["SourceEndpoint.InstanceType"] = "RDS"
	request.QueryParams["DestinationEndpoint.InstanceType"] = "RDS"
	request.QueryParams["MigrationMode.StructureIntialization"] = "true"
	request.QueryParams["MigrationMode.DataIntialization"] = "true"
	request.QueryParams["MigrationMode.DataSynchronization"] = "false"
	request.QueryParams["MigrationObject"] = "[{\"DBName\":\"mulan_billing_production\",\"NewDBName\":\"mulan_billing_test\",\"AllTable\": true},{\"DBName\":\"mulan_bis_production\",\"NewDBName\":\"mulan_bis_production\",\"AllTable\": true},{\"DBName\":\"mulan_credits_production\",\"NewDBName\":\"mulan_credits_production\",\"AllTable\": true},{\"DBName\":\"mulan_inventory_production\",\"NewDBName\":\"mulan_inventory_test\",\"AllTable\": true},{\"DBName\":\"mulan_order_production\",\"NewDBName\":\"mulan_order_test\",\"AllTable\": true},{\"DBName\":\"mulan_print_production\",\"NewDBName\":\"mulan_print_test\",\"AllTable\": true},{\"DBName\":\"mulan_workday_production\",\"NewDBName\":\"mulan_workday_test\",\"AllTable\": true},{\"DBName\":\"wwc_ads_production\",\"NewDBName\":\"wwc_ads_test\",\"AllTable\": true},{\"DBName\":\"wwc_auth_production\",\"NewDBName\":\"wwc_auth_test\",\"AllTable\": true},{\"DBName\":\"wwc_order_production\",\"NewDBName\":\"wwc_order_test\",\"AllTable\": true},{\"DBName\":\"wwc_paidevent_production\",\"NewDBName\":\"wwc_paidevent_test\",\"AllTable\": true},{\"DBName\":\"wwc_pricing_production\",\"NewDBName\":\"wwc_pricing_test\",\"AllTable\": true},{\"DBName\":\"wwc_product_production\",\"NewDBName\":\"wwc_product_test\",\"AllTable\": true},{\"DBName\":\"wwc_support_production\",\"NewDBName\":\"wwc_support_test\",\"AllTable\": true}]"
	request.QueryParams["MigrationJobId"] = id
	request.QueryParams["SourceEndpoint.InstanceID"] = "rm-uf65zox8h1mi44j67"
	request.QueryParams["SourceEndpoint.EngineName"] = "MySQL"
	request.QueryParams["SourceEndpoint.Region"] = "cn-shanghai"
	request.QueryParams["SourceEndpoint.UserName"] = "dms"
	request.QueryParams["SourceEndpoint.Password"] = "*v*#lgiqMuEKjheW"
	request.QueryParams["DestinationEndpoint.InstanceID"] = "rm-uf65m74z317qid2i8"
	request.QueryParams["DestinationEndpoint.EngineName"] = "MySQL"
	request.QueryParams["DestinationEndpoint.UserName"] = "dms"
	request.QueryParams["DestinationEndpoint.Password"] = "q0OzU4B^bpnEqkS4"
	request.QueryParams["MigrationReserved"] = "{ 	\"autoStartModulesAfterConfig\": \"all\", 	\"targetTableMode\": 2 }"
	response, err := dtsclient.ProcessCommonRequest(request)
	if err != nil {
		panic(err)
	}
	fmt.Println(response.GetHttpStatus())
}

func mysql_java_test_sales_wizard() {
	//create
	id := create_instance()
	//config && start
	request := requests.NewCommonRequest()
	request.Method = "POST"
	request.Scheme = "https" // https | http
	request.Domain = "dts.aliyuncs.com"
	request.Version = "2020-01-01"
	request.ApiName = "ConfigureMigrationJob"
	request.QueryParams["RegionId"] = "cn-shanghai"
	request.QueryParams["MigrationJobName"] = "mysql_java_test_sales_wizard"
	request.QueryParams["SourceEndpoint.InstanceType"] = "RDS"
	request.QueryParams["DestinationEndpoint.InstanceType"] = "RDS"
	request.QueryParams["MigrationMode.StructureIntialization"] = "true"
	request.QueryParams["MigrationMode.DataIntialization"] = "true"
	request.QueryParams["MigrationMode.DataSynchronization"] = "false"
	request.QueryParams["MigrationObject"] = "[{\"DBName\":\"sales_wizard_java_production\",\"NewDBName\":\"sales_wizard_java_test\",\"AllTable\": true}]"
	request.QueryParams["MigrationJobId"] = id
	request.QueryParams["SourceEndpoint.InstanceID"] = "rm-uf6294i9i0oip39yt"
	request.QueryParams["SourceEndpoint.EngineName"] = "MySQL"
	request.QueryParams["SourceEndpoint.Region"] = "cn-shanghai"
	request.QueryParams["SourceEndpoint.UserName"] = "dms"
	request.QueryParams["SourceEndpoint.Password"] = "*v*#lgiqMuEKjheW"
	request.QueryParams["DestinationEndpoint.InstanceID"] = "rm-uf67w36764okf9c10"
	request.QueryParams["DestinationEndpoint.EngineName"] = "MySQL"
	request.QueryParams["DestinationEndpoint.UserName"] = "dms"
	request.QueryParams["DestinationEndpoint.Password"] = "q0OzU4B^bpnEqkS4"
	request.QueryParams["MigrationReserved"] = "{ 	\"autoStartModulesAfterConfig\": \"all\", 	\"targetTableMode\": 2 }"
	response, err := dtsclient.ProcessCommonRequest(request)
	if err != nil {
		panic(err)
	}
	fmt.Println(response.GetHttpStatus())
}

func mysql_java_test_wwcnapi() {
	//create
	id := create_instance()
	//config && start
	request := requests.NewCommonRequest()
	request.Method = "POST"
	request.Scheme = "https" // https | http
	request.Domain = "dts.aliyuncs.com"
	request.Version = "2020-01-01"
	request.ApiName = "ConfigureMigrationJob"
	request.QueryParams["RegionId"] = "cn-shanghai"
	request.QueryParams["MigrationJobName"] = "mysql_java_test_wwcnapi"
	request.QueryParams["SourceEndpoint.InstanceType"] = "RDS"
	request.QueryParams["DestinationEndpoint.InstanceType"] = "RDS"
	request.QueryParams["MigrationMode.StructureIntialization"] = "true"
	request.QueryParams["MigrationMode.DataIntialization"] = "true"
	request.QueryParams["MigrationMode.DataSynchronization"] = "false"
	request.QueryParams["MigrationObject"] = "[{\"DBName\":\"wwcnapi_db_prod\",\"NewDBName\":\"wwcnapi_db_test\",\"AllTable\": true}]"
	request.QueryParams["MigrationJobId"] = id
	request.QueryParams["SourceEndpoint.InstanceID"] = "rm-uf625bi5skv69igu8"
	request.QueryParams["SourceEndpoint.EngineName"] = "MySQL"
	request.QueryParams["SourceEndpoint.Region"] = "cn-shanghai"
	request.QueryParams["SourceEndpoint.UserName"] = "dms"
	request.QueryParams["SourceEndpoint.Password"] = "*v*#lgiqMuEKjheW"
	request.QueryParams["DestinationEndpoint.InstanceID"] = "rm-uf6wm1h60mzoyx7gv"
	request.QueryParams["DestinationEndpoint.EngineName"] = "MySQL"
	request.QueryParams["DestinationEndpoint.UserName"] = "dms"
	request.QueryParams["DestinationEndpoint.Password"] = "q0OzU4B^bpnEqkS4"
	request.QueryParams["MigrationReserved"] = "{ 	\"autoStartModulesAfterConfig\": \"all\", 	\"targetTableMode\": 2 }"
	response, err := dtsclient.ProcessCommonRequest(request)
	if err != nil {
		panic(err)
	}
	fmt.Println(response.GetHttpStatus())
}

func settlement_reports_test() {
	//create
	id := create_instance()
	//config && start
	request := requests.NewCommonRequest()
	request.Method = "POST"
	request.Scheme = "https" // https | http
	request.Domain = "dts.aliyuncs.com"
	request.Version = "2020-01-01"
	request.ApiName = "ConfigureMigrationJob"
	request.QueryParams["RegionId"] = "cn-shanghai"
	request.QueryParams["MigrationJobName"] = "settlement_reports_test"
	request.QueryParams["SourceEndpoint.InstanceType"] = "RDS"
	request.QueryParams["DestinationEndpoint.InstanceType"] = "RDS"
	request.QueryParams["MigrationMode.StructureIntialization"] = "true"
	request.QueryParams["MigrationMode.DataIntialization"] = "true"
	request.QueryParams["MigrationMode.DataSynchronization"] = "false"
	request.QueryParams["MigrationObject"] = "[{\"DBName\":\"settlement_reports_production\",\"NewDBName\":\"settlement_reports_test\",\"AllTable\": true}]"
	request.QueryParams["MigrationJobId"] = id
	request.QueryParams["SourceEndpoint.InstanceID"] = "rm-uf64dy2k6x9w6pg72"
	request.QueryParams["SourceEndpoint.EngineName"] = "MySQL"
	request.QueryParams["SourceEndpoint.Region"] = "cn-shanghai"
	request.QueryParams["SourceEndpoint.UserName"] = "dms"
	request.QueryParams["SourceEndpoint.Password"] = "*v*#lgiqMuEKjheW"
	request.QueryParams["DestinationEndpoint.InstanceID"] = "rm-uf6n3x710st980mv6"
	request.QueryParams["DestinationEndpoint.EngineName"] = "MySQL"
	request.QueryParams["DestinationEndpoint.UserName"] = "dms"
	request.QueryParams["DestinationEndpoint.Password"] = "q0OzU4B^bpnEqkS4"
	request.QueryParams["MigrationReserved"] = "{ 	\"autoStartModulesAfterConfig\": \"all\", 	\"targetTableMode\": 2 }"
	response, err := dtsclient.ProcessCommonRequest(request)
	if err != nil {
		panic(err)
	}
	fmt.Println(response.GetHttpStatus())
}

func hotdesk_test() {
	//create
	id := create_instance()
	//config && start
	request := requests.NewCommonRequest()
	request.Method = "POST"
	request.Scheme = "https" // https | http
	request.Domain = "dts.aliyuncs.com"
	request.Version = "2020-01-01"
	request.ApiName = "ConfigureMigrationJob"
	request.QueryParams["RegionId"] = "cn-shanghai"
	request.QueryParams["MigrationJobName"] = "hotdesk_test"
	request.QueryParams["SourceEndpoint.InstanceType"] = "RDS"
	request.QueryParams["DestinationEndpoint.InstanceType"] = "RDS"
	request.QueryParams["MigrationMode.StructureIntialization"] = "true"
	request.QueryParams["MigrationMode.DataIntialization"] = "true"
	request.QueryParams["MigrationMode.DataSynchronization"] = "false"
	request.QueryParams["MigrationObject"] = "[{\"DBName\":\"hotdesk_production\",\"NewDBName\":\"hotdesk_test\",\"AllTable\": true}]"
	request.QueryParams["MigrationJobId"] = id
	request.QueryParams["SourceEndpoint.InstanceID"] = "rm-uf66m36jo800f1fa9"
	request.QueryParams["SourceEndpoint.EngineName"] = "MySQL"
	request.QueryParams["SourceEndpoint.Region"] = "cn-shanghai"
	request.QueryParams["SourceEndpoint.UserName"] = "dms"
	request.QueryParams["SourceEndpoint.Password"] = "*v*#lgiqMuEKjheW"
	request.QueryParams["DestinationEndpoint.InstanceID"] = "rm-uf618wj9qonhu178p"
	request.QueryParams["DestinationEndpoint.EngineName"] = "MySQL"
	request.QueryParams["DestinationEndpoint.UserName"] = "dms"
	request.QueryParams["DestinationEndpoint.Password"] = "q0OzU4B^bpnEqkS4"
	request.QueryParams["MigrationReserved"] = "{ 	\"autoStartModulesAfterConfig\": \"all\", 	\"targetTableMode\": 2 }"
	response, err := dtsclient.ProcessCommonRequest(request)
	if err != nil {
		panic(err)
	}
	fmt.Println(response.GetHttpStatus())
}

func translation_test() {
	//create
	id := create_instance()
	//config && start
	request := requests.NewCommonRequest()
	request.Method = "POST"
	request.Scheme = "https" // https | http
	request.Domain = "dts.aliyuncs.com"
	request.Version = "2020-01-01"
	request.ApiName = "ConfigureMigrationJob"
	request.QueryParams["RegionId"] = "cn-shanghai"
	request.QueryParams["MigrationJobName"] = "translation_test"
	request.QueryParams["SourceEndpoint.InstanceType"] = "RDS"
	request.QueryParams["DestinationEndpoint.InstanceType"] = "RDS"
	request.QueryParams["MigrationMode.StructureIntialization"] = "true"
	request.QueryParams["MigrationMode.DataIntialization"] = "true"
	request.QueryParams["MigrationMode.DataSynchronization"] = "false"
	request.QueryParams["MigrationObject"] = "[{\"DBName\":\"translation_production\",\"NewDBName\":\"translation_test\",\"AllTable\": true}]"
	request.QueryParams["MigrationJobId"] = id
	request.QueryParams["SourceEndpoint.InstanceID"] = "rm-uf69v9049c8ffm8tn"
	request.QueryParams["SourceEndpoint.EngineName"] = "MySQL"
	request.QueryParams["SourceEndpoint.Region"] = "cn-shanghai"
	request.QueryParams["SourceEndpoint.UserName"] = "dms"
	request.QueryParams["SourceEndpoint.Password"] = "*v*#lgiqMuEKjheW"
	request.QueryParams["DestinationEndpoint.InstanceID"] = "rm-uf6c27gmkwf91l376"
	request.QueryParams["DestinationEndpoint.EngineName"] = "MySQL"
	request.QueryParams["DestinationEndpoint.UserName"] = "dms"
	request.QueryParams["DestinationEndpoint.Password"] = "q0OzU4B^bpnEqkS4"
	request.QueryParams["MigrationReserved"] = "{ 	\"autoStartModulesAfterConfig\": \"all\", 	\"targetTableMode\": 2 }"
	response, err := dtsclient.ProcessCommonRequest(request)
	if err != nil {
		panic(err)
	}
	fmt.Println(response.GetHttpStatus())
}

func contracts_test() {
	//create
	id := create_instance()
	//config && start
	request := requests.NewCommonRequest()
	request.Method = "POST"
	request.Scheme = "https" // https | http
	request.Domain = "dts.aliyuncs.com"
	request.Version = "2020-01-01"
	request.ApiName = "ConfigureMigrationJob"
	request.QueryParams["RegionId"] = "cn-shanghai"
	request.QueryParams["MigrationJobName"] = "contracts_test"
	request.QueryParams["SourceEndpoint.InstanceType"] = "RDS"
	request.QueryParams["DestinationEndpoint.InstanceType"] = "RDS"
	request.QueryParams["MigrationMode.StructureIntialization"] = "true"
	request.QueryParams["MigrationMode.DataIntialization"] = "true"
	request.QueryParams["MigrationMode.DataSynchronization"] = "false"
	request.QueryParams["MigrationObject"] = "[{\"DBName\":\"contracts_production\",\"NewDBName\":\"contracts_test\",\"AllTable\": true}]"
	request.QueryParams["MigrationJobId"] = id
	request.QueryParams["SourceEndpoint.InstanceID"] = "rm-uf6b5523yvawl90ah"
	request.QueryParams["SourceEndpoint.EngineName"] = "MySQL"
	request.QueryParams["SourceEndpoint.Region"] = "cn-shanghai"
	request.QueryParams["SourceEndpoint.UserName"] = "dms"
	request.QueryParams["SourceEndpoint.Password"] = "*v*#lgiqMuEKjheW"
	request.QueryParams["DestinationEndpoint.InstanceID"] = "rm-uf6x1q40zhky67bi0"
	request.QueryParams["DestinationEndpoint.EngineName"] = "MySQL"
	request.QueryParams["DestinationEndpoint.UserName"] = "dms"
	request.QueryParams["DestinationEndpoint.Password"] = "q0OzU4B^bpnEqkS4"
	request.QueryParams["MigrationReserved"] = "{ 	\"autoStartModulesAfterConfig\": \"all\", 	\"targetTableMode\": 2 }"
	response, err := dtsclient.ProcessCommonRequest(request)
	if err != nil {
		panic(err)
	}
	fmt.Println(response.GetHttpStatus())
}

func account_overview_test() {
	//create
	id := create_instance()
	//config && start
	request := requests.NewCommonRequest()
	request.Method = "POST"
	request.Scheme = "https" // https | http
	request.Domain = "dts.aliyuncs.com"
	request.Version = "2020-01-01"
	request.ApiName = "ConfigureMigrationJob"
	request.QueryParams["RegionId"] = "cn-shanghai"
	request.QueryParams["MigrationJobName"] = "account_overview_test"
	request.QueryParams["SourceEndpoint.InstanceType"] = "RDS"
	request.QueryParams["DestinationEndpoint.InstanceType"] = "RDS"
	request.QueryParams["MigrationMode.StructureIntialization"] = "true"
	request.QueryParams["MigrationMode.DataIntialization"] = "true"
	request.QueryParams["MigrationMode.DataSynchronization"] = "false"
	request.QueryParams["MigrationObject"] = "[{\"DBName\":\"account_overview_production\",\"NewDBName\":\"account_overview_test\",\"AllTable\": true}]"
	request.QueryParams["MigrationJobId"] = id
	request.QueryParams["SourceEndpoint.InstanceID"] = "rm-uf6641878c00e4wj1"
	request.QueryParams["SourceEndpoint.EngineName"] = "MySQL"
	request.QueryParams["SourceEndpoint.Region"] = "cn-shanghai"
	request.QueryParams["SourceEndpoint.UserName"] = "dms"
	request.QueryParams["SourceEndpoint.Password"] = "*v*#lgiqMuEKjheW"
	request.QueryParams["DestinationEndpoint.InstanceID"] = "rm-uf66apv43f5818iro"
	request.QueryParams["DestinationEndpoint.EngineName"] = "MySQL"
	request.QueryParams["DestinationEndpoint.UserName"] = "dms"
	request.QueryParams["DestinationEndpoint.Password"] = "q0OzU4B^bpnEqkS4"
	request.QueryParams["MigrationReserved"] = "{ 	\"autoStartModulesAfterConfig\": \"all\", 	\"targetTableMode\": 2 }"
	response, err := dtsclient.ProcessCommonRequest(request)
	if err != nil {
		panic(err)
	}
	fmt.Println(response.GetHttpStatus())
}

func spacecowboy_test() {
	//create
	id := create_instance()
	//config && start
	request := requests.NewCommonRequest()
	request.Method = "POST"
	request.Scheme = "https" // https | http
	request.Domain = "dts.aliyuncs.com"
	request.Version = "2020-01-01"
	request.ApiName = "ConfigureMigrationJob"
	request.QueryParams["RegionId"] = "cn-shanghai"
	request.QueryParams["MigrationJobName"] = "spacecowboy_test"
	request.QueryParams["SourceEndpoint.InstanceType"] = "RDS"
	request.QueryParams["DestinationEndpoint.InstanceType"] = "RDS"
	request.QueryParams["MigrationMode.StructureIntialization"] = "true"
	request.QueryParams["MigrationMode.DataIntialization"] = "true"
	request.QueryParams["MigrationMode.DataSynchronization"] = "false"
	request.QueryParams["MigrationObject"] = "[{\"DBName\":\"spacecowboy_prod\",\"NewDBName\":\"spacecowboy_test\",\"AllTable\": true}]"
	request.QueryParams["MigrationJobId"] = id
	request.QueryParams["SourceEndpoint.InstanceID"] = "rm-uf6zh4h890u6o2385"
	request.QueryParams["SourceEndpoint.EngineName"] = "MySQL"
	request.QueryParams["SourceEndpoint.Region"] = "cn-shanghai"
	request.QueryParams["SourceEndpoint.UserName"] = "dms"
	request.QueryParams["SourceEndpoint.Password"] = "*v*#lgiqMuEKjheW"
	request.QueryParams["DestinationEndpoint.InstanceID"] = "rm-uf64a69n5h5g1p1n4"
	request.QueryParams["DestinationEndpoint.EngineName"] = "MySQL"
	request.QueryParams["DestinationEndpoint.UserName"] = "dms"
	request.QueryParams["DestinationEndpoint.Password"] = "q0OzU4B^bpnEqkS4"
	request.QueryParams["MigrationReserved"] = "{ 	\"autoStartModulesAfterConfig\": \"all\", 	\"targetTableMode\": 2 }"
	response, err := dtsclient.ProcessCommonRequest(request)
	if err != nil {
		panic(err)
	}
	fmt.Println(response.GetHttpStatus())
}

func china_pos_payments_service_test() {
	//create
	id := create_instance()
	//config && start
	request := requests.NewCommonRequest()
	request.Method = "POST"
	request.Scheme = "https" // https | http
	request.Domain = "dts.aliyuncs.com"
	request.Version = "2020-01-01"
	request.ApiName = "ConfigureMigrationJob"
	request.QueryParams["RegionId"] = "cn-shanghai"
	request.QueryParams["MigrationJobName"] = "china_pos_payments_service_test"
	request.QueryParams["SourceEndpoint.InstanceType"] = "RDS"
	request.QueryParams["DestinationEndpoint.InstanceType"] = "RDS"
	request.QueryParams["MigrationMode.StructureIntialization"] = "true"
	request.QueryParams["MigrationMode.DataIntialization"] = "true"
	request.QueryParams["MigrationMode.DataSynchronization"] = "false"
	request.QueryParams["MigrationObject"] = "[{\"DBName\":\"china_pos_payments_service_production\",\"NewDBName\":\"china_pos_payments_service_test\",\"AllTable\": true}]"
	request.QueryParams["MigrationJobId"] = id
	request.QueryParams["SourceEndpoint.InstanceID"] = "rm-uf6y47s55r5mui871"
	request.QueryParams["SourceEndpoint.EngineName"] = "MySQL"
	request.QueryParams["SourceEndpoint.Region"] = "cn-shanghai"
	request.QueryParams["SourceEndpoint.UserName"] = "dms"
	request.QueryParams["SourceEndpoint.Password"] = "*v*#lgiqMuEKjheW"
	request.QueryParams["DestinationEndpoint.InstanceID"] = "rm-uf6cak131q04c6pn7"
	request.QueryParams["DestinationEndpoint.EngineName"] = "MySQL"
	request.QueryParams["DestinationEndpoint.UserName"] = "dms"
	request.QueryParams["DestinationEndpoint.Password"] = "q0OzU4B^bpnEqkS4"
	request.QueryParams["MigrationReserved"] = "{ 	\"autoStartModulesAfterConfig\": \"all\", 	\"targetTableMode\": 2 }"
	response, err := dtsclient.ProcessCommonRequest(request)
	if err != nil {
		panic(err)
	}
	fmt.Println(response.GetHttpStatus())
}

//rooms-cn实例
func keycards_test() {
	//clean_database
	clean_database_pgsql("pgm-uf6s09cqwe5710u5", "keycards_test", "UTF8", "api")
	//create
	id := create_instance()
	//config && start
	request := requests.NewCommonRequest()
	request.Method = "POST"
	request.Scheme = "https" // https | http
	request.Domain = "dts.aliyuncs.com"
	request.Version = "2020-01-01"
	request.ApiName = "ConfigureMigrationJob"
	request.QueryParams["RegionId"] = "cn-shanghai"
	request.QueryParams["MigrationJobName"] = "keycards_test"
	request.QueryParams["SourceEndpoint.InstanceType"] = "RDS"
	request.QueryParams["DestinationEndpoint.InstanceType"] = "RDS"
	request.QueryParams["MigrationMode.StructureIntialization"] = "true"
	request.QueryParams["MigrationMode.DataIntialization"] = "true"
	request.QueryParams["MigrationMode.DataSynchronization"] = "false"
	request.QueryParams["MigrationObject"] = "[{\"DBName\":\"public\"}]"
	request.QueryParams["MigrationJobId"] = id
	request.QueryParams["SourceEndpoint.InstanceID"] = "pgm-uf6507s5zam3f23y"
	request.QueryParams["SourceEndpoint.EngineName"] = "PostgreSQL"
	request.QueryParams["SourceEndpoint.Region"] = "cn-shanghai"
	request.QueryParams["SourceEndpoint.UserName"] = "wework_root"
	request.QueryParams["SourceEndpoint.Password"] = "cFk2H28Vz0gUJwUasPvcyuL5xl8wKkcC"
	request.QueryParams["SourceEndpoint.DatabaseName"] = "keycards_production"
	request.QueryParams["DestinationEndpoint.InstanceID"] = "pgm-uf6s09cqwe5710u5"
	request.QueryParams["DestinationEndpoint.EngineName"] = "PostgreSQL"
	request.QueryParams["DestinationEndpoint.UserName"] = "api"
	request.QueryParams["DestinationEndpoint.Password"] = "sD7C4yBvs9vPBja6"
	request.QueryParams["DestinationEndpoint.DataBaseName"] = "keycards_test"
	request.QueryParams["MigrationReserved"] = "{ 	\"autoStartModulesAfterConfig\": \"all\", 	\"targetTableMode\": 2 }"
	response, err := dtsclient.ProcessCommonRequest(request)
	if err != nil {
		panic(err)
	}
	fmt.Println(response.GetHttpStatus())
}

//bp-pg实例
func fapiao_test() {
	//clean_database
	clean_database_pgsql("pgm-uf6ecw5006vsi4z9", "fapiao_test", "UTF8", "api")
	//create
	id := create_instance()
	//config && start
	request := requests.NewCommonRequest()
	request.Method = "POST"
	request.Scheme = "https" // https | http
	request.Domain = "dts.aliyuncs.com"
	request.Version = "2020-01-01"
	request.ApiName = "ConfigureMigrationJob"
	request.QueryParams["RegionId"] = "cn-shanghai"
	request.QueryParams["MigrationJobName"] = "fapiao_test"
	request.QueryParams["SourceEndpoint.InstanceType"] = "RDS"
	request.QueryParams["DestinationEndpoint.InstanceType"] = "RDS"
	request.QueryParams["MigrationMode.StructureIntialization"] = "true"
	request.QueryParams["MigrationMode.DataIntialization"] = "true"
	request.QueryParams["MigrationMode.DataSynchronization"] = "false"
	request.QueryParams["MigrationObject"] = "[{\"DBName\":\"public\"}]"
	request.QueryParams["MigrationJobId"] = id
	request.QueryParams["SourceEndpoint.InstanceID"] = "pgm-uf68a9116jv634nk"
	request.QueryParams["SourceEndpoint.EngineName"] = "PostgreSQL"
	request.QueryParams["SourceEndpoint.Region"] = "cn-shanghai"
	request.QueryParams["SourceEndpoint.UserName"] = "api"
	request.QueryParams["SourceEndpoint.Password"] = "qu9xnfdrZgXyBkHhkLyk"
	request.QueryParams["SourceEndpoint.DatabaseName"] = "fapiao_production"
	request.QueryParams["DestinationEndpoint.InstanceID"] = "pgm-uf6ecw5006vsi4z9"
	request.QueryParams["DestinationEndpoint.EngineName"] = "PostgreSQL"
	request.QueryParams["DestinationEndpoint.UserName"] = "dms"
	request.QueryParams["DestinationEndpoint.Password"] = "q0OzU4B^bpnEqkS4"
	request.QueryParams["DestinationEndpoint.DataBaseName"] = "fapiao_test"
	request.QueryParams["MigrationReserved"] = "{ 	\"autoStartModulesAfterConfig\": \"all\", 	\"targetTableMode\": 2 }"
	response, err := dtsclient.ProcessCommonRequest(request)
	if err != nil {
		panic(err)
	}
	fmt.Println(response.GetHttpStatus())
}

func main() {
	Init()
	//mysql_java_test_mulan_db_v56()
	//mysql_java_test_mulan_db_v57()
	//mysql_java_test_sales_wizard()
	//mysql_java_test_wwcnapi()
	//settlement_reports_test()
	//hotdesk_test()
	//translation_test()
	//contracts_test()
	//account_overview_test()
	//china_pos_payments_service_test()
	//spacecowboy_test()
	//keycards_test()
	fapiao_test()
}
