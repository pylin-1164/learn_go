package winservices

import (
	"fmt"
	"github.com/shirou/gopsutil/winservices"
	"golang.org/x/sys/windows"
	"golang.org/x/sys/windows/svc/mgr"
	"unsafe"
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

func List(){

	services, _ := winservices.ListServices()

	for _, service := range services {
		func(){
			defer func() {
				e := recover()
				if e != nil{
					fmt.Printf("error :  %s \n",service.Name)
				}

			}()
			newservice, _ := winservices.NewService(service.Name)
			newservice.GetServiceDetail()

			fmt.Println("Name:",newservice.Name,"DisplayName:", newservice.Config.DisplayName, "Binary Path:", newservice.Config.BinaryPathName, "State: ", newservice.Status.State)
		}()
	}
}


/**
 * xp系统兼容
 */
func QueryServiceStatus()(serviceList []*ServiceInfo){
	services, _ := winservices.ListServices()
	for _,s := range services {
		serviceInfo := queryServiceByName(s)
		if serviceInfo == nil{
			continue
		}
		serviceList = append(serviceList,serviceInfo)
	}
	return serviceList
}

func queryServiceByName(s winservices.Service) *ServiceInfo{
	defer func() {
		e := recover()
		if e != nil{
			//fmt.Printf("error :  %s \n",s.Name)
		}

	}()

	m, _ := mgr.Connect()
	defer m.Disconnect()

	service, _ := m.OpenService(s.Name)
	var bytesNeeded uint32 = 2048
	var buf []byte
	var configBuf []byte

	buf = make([]byte, bytesNeeded)
	configBuf = make([]byte, bytesNeeded)
	p := (*windows.SERVICE_STATUS_PROCESS)(unsafe.Pointer(&buf[0]))
	cf := (*windows.QUERY_SERVICE_CONFIG)(unsafe.Pointer(&configBuf[0]))

	windows.QueryServiceStatusEx(service.Handle, windows.SC_STATUS_PROCESS_INFO, &buf[0], uint32(len(buf)), &bytesNeeded)
	windows.QueryServiceConfig(service.Handle,cf,uint32(len(configBuf)),&bytesNeeded)
	return &ServiceInfo{
		ServiceName:        service.Name,
		ServiceDisplayName: windows.UTF16PtrToString(cf.DisplayName),
		ServiceStatus:      int(p.CurrentState),
		ServiceStatusCN:    ServiceStatusCNMap[int(p.CurrentState)],
	}
}