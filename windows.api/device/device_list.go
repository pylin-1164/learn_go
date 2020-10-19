package device

import (
	"fmt"
	"golang.org/x/sys/windows"
	"windows.api/api/setup"
)

var deviceClassNetGUID = windows.GUID{Data1: 0x4d36e972, Data2: 0xe325, Data3: 0x11ce, Data4: [8]byte{0xbf, 0xc1, 0x08, 0x00, 0x2b, 0xe1, 0x03, 0x18}}
var deviceInterfaceNetGUID = windows.GUID{Data1: 0xcac88484, Data2: 0x7515, Data3: 0x4c03, Data4: [8]byte{0x82, 0xe6, 0x71, 0xa8, 0x7a, 0xba, 0xc3, 0x61}}

func List(){

	devInfo, err := setup.SetupDiGetClassDevsEx(nil, "", 0, setup.DIGCF_ALLCLASSES, setup.DevInfo(0), "")
	if err != nil {
		fmt.Errorf("SetupDiGetClassDevsEx() failed: %v", err)
		return
	}
	defer devInfo.Close()
	for index := 0; ; index++ {
		b := findProperty(devInfo, index)

		if !b{
			break
		}
	}

}

func findProperty(devInfo setup.DevInfo,index int) bool{
	defer func() {
		recover()
	}()
	devInfoData, err := devInfo.EnumDeviceInfo(index)
	if err != nil {
		if err == windows.ERROR_NO_MORE_ITEMS {
			return false
		}
		return true
	}
	property, err := devInfo.GetDeviceRegistryProperty(devInfoData, setup.SPDRP_PHYSICAL_DEVICE_OBJECT_NAME)
	if err != nil{
		return true
	}
	property, _ = devInfo.GetDeviceRegistryProperty(devInfoData, setup.SPDRP_DEVICEDESC)
	fmt.Printf("%s \n",property)
	return true
}
