package winservices

import (
	"fmt"
	"github.com/shirou/gopsutil/winservices"
	"golang.org/x/sys/windows"
	"golang.org/x/sys/windows/svc/mgr"
	"unsafe"
)

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
func QueryServiceStatus(){
	services, _ := winservices.ListServices()
	for _,s := range services {
		queryServiceByName(s)

	}

}

func queryServiceByName(s winservices.Service) {
	defer func() {
		e := recover()
		if e != nil{
			fmt.Printf("error :  %s \n",s.Name)
		}

	}()

	m, _ := mgr.Connect()
	defer m.Disconnect()

	service, _ := m.OpenService(s.Name)
	var bytesNeeded uint32 = 2048
	var buf []byte

	buf = make([]byte, bytesNeeded)
	p := (*windows.SERVICE_STATUS_PROCESS)(unsafe.Pointer(&buf[0]))
	windows.QueryServiceStatusEx(service.Handle, windows.SC_STATUS_PROCESS_INFO, &buf[0], uint32(len(buf)), &bytesNeeded)
	fmt.Println(service.Name,p.CurrentState)
}