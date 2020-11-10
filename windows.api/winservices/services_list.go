package winservices

import (
	"fmt"
	wapi "github.com/iamacarpet/go-win64api"
)

const (
	//启动挂起中
	Service_Status_Continue_Pending  = 5
	//暂停挂起中
	Service_Status_Pause_Pending  = 6
	//暂停
	Service_Status_Paused  = 7
	//运行中
	Service_Status_Running  = 4
	//启动中
	Service_Status_Start_Pending  = 2
	//停止中
	Service_Status_Stop_Pending  = 3
	//已停止
	Service_Status_Stoped  = 1
)

var ServiceStatusCNMap = map[int]string{
		Service_Status_Continue_Pending:"启动挂起中",
		Service_Status_Pause_Pending:"暂停挂起中",
		Service_Status_Paused:"暂停",
		Service_Status_Running:"运行中",
		Service_Status_Start_Pending:"启动中",
		Service_Status_Stop_Pending:"停止中",
		Service_Status_Stoped:"已停止",
	}

type ServiceInfo struct {
	ServiceName 		string
	ServiceDisplayName	string
	ServiceStatus 		int
	ServiceStatusCN		string
}

//xp系统兼容版本请查看提交历史

func QueryServiceList()(serviceList []*ServiceInfo){
	svc, err := wapi.GetServices()
	if err != nil {
		fmt.Printf("%s\r\n", err.Error())
	}

	for _, v := range svc {
		serviceInfo := ServiceInfo{
			ServiceName:        v.SCName,
			ServiceDisplayName: v.DisplayName,
			ServiceStatus:      int(v.Status),
			ServiceStatusCN:    ServiceStatusCNMap[int(v.Status)],
		}
		serviceList = append(serviceList,&serviceInfo)
		//fmt.Printf("%-50s - %-75s - Status: %-20s - Accept Stop: %-5t, Running Pid: %d\r\n", v.SCName, v.DisplayName, v.StatusText, v.AcceptStop, v.RunningPid)
	}
	return
}