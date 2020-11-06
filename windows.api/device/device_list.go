package device

import (
	"fmt"
	"golang.org/x/sys/windows"
	"windows.api/api/setup"
)

/** 常用硬件设备GUID https://www.cnblogs.com/hester/p/7591876.html
	4D36E980-E325-11CE-BFC1-08002BE10318        软驱
	4D36E965-E325-11CE-BFC1-08002BE10318        光驱
	4D36E97D-E325-11CE-BFC1-08002BE10318       系统设备
	6D807884-7D21-11CF-801C-08002BE10318        磁带机
	36FC9E60-C465-11CF-8056-444553540000         USB
	4D36E964-E325-11CE-BFC1-08002BE10318       适配器
	D45B1C18-C8FA-11D1-9F77-0000F805F530        APMSUPPORT
	E0CBF06C-CD8B-4647-BB8A-253B43F0F974       蓝牙设备
	4D36E966-E325-11CE-BFC1-08002BE10318       电脑
	6BDD1FC2-810F-11D0-BEC7-08002BE2092F       解码器
	6BDD1FC3-810F-11D0-BEC7-08002BE2092F       GPS——Global Positioning System
	4D36E976-E325-11CE-BFC1-08002BE10318       No Driver
	8ECC055D-047F-11D1-A537-0000F8753ED1       Non-plug And Play Drivers
	4D36E97E-E325-11CE-BFC1-08002BE10318       Other Devices
	4D36E97A-E325-11CE-BFC1-08002BE10318       Printer Upgrade
	4D36E97C-E325-11CE-BFC1-08002BE10318       声音设备
	4D36E96F-E325-11CE-BFC1-08002BE10318       鼠标
	4D36E96B-E325-11CE-BFC1-08002BE10318       键盘
	4D36E96C-E325-11CE-BFC1-08002BE10318       音频视频设备
 */

var deviceClassNetGUID = windows.GUID{Data1: 0x4d36e972, Data2: 0xe325, Data3: 0x11ce, Data4: [8]byte{0xbf, 0xc1, 0x08, 0x00, 0x2b, 0xe1, 0x03, 0x18}}
var deviceInterfaceNetGUID = windows.GUID{Data1: 0xcac88484, Data2: 0x7515, Data3: 0x4c03, Data4: [8]byte{0x82, 0xe6, 0x71, 0xa8, 0x7a, 0xba, 0xc3, 0x61}}

//USB 36FC9E60-C465-11CF-8056-444553540000
var deviceUsbGUID = windows.GUID{
	Data1: 0x36fc9e60,
	Data2: 0xc465,
	Data3: 0x11cf,
	Data4: [8]byte{0x80,0x56,0x44,0x45,0x53,0x54,0x00,0x00},
}

//声音设备 4D36E97C-E325-11CE-BFC1-08002BE10318
var deviceVideoGUI = windows.GUID{
	Data1: 0x4d36e97c,
	Data2: 0xe325,
	Data3: 0x11ce,
	Data4: [8]byte{0xbf,0xc1,0x08,0x00,0x2b,0xe1,0x03,0x18},
}


//键盘驱动 4D36E96B-E325-11CE-BFC1-08002BE10318
var deviceKeyboardsGUI = windows.GUID{
	Data1: 0x4D36E96B,
	Data2: 0xE325,
	Data3: 0x11CE,
	Data4: [8]byte{0xBF,0xC1,0x08,0x00,0x2B,0xE1,0x03,0x18},
}

//蓝牙设备 E0CBF06C-CD8B-4647-BB8A-253B43F0F974
var deviceBlueGUI = windows.GUID{
	Data1: 0xE0CBF06C,
	Data2: 0xCD8B,
	Data3: 0x4647,
	Data4: [8]byte{0xBB,0x8A,0x25,0x3B,0x43,0xF0,0xF9,0x74},
}

//鼠标 Devices	4D36E97E-E325-11CE-BFC1-08002BE10318
var deviceMouseGUI = windows.GUID{
	Data1: 0x4D36E96F,
	Data2: 0xE325,
	Data3: 0x11CE,
	Data4: [8]byte{0xBF,0xC1,0x08,0x00,0x2B,0xE1,0x03,0x18},
}
//视频音频设备 Devices	4D36E96C-E325-11CE-BFC1-08002BE10318
var deviceMediaGUI = windows.GUID{
	Data1: 0x4D36E96C,
	Data2: 0xE325,
	Data3: 0x11CE,
	Data4: [8]byte{0xBF,0xC1,0x08,0x00,0x2B,0xE1,0x03,0x18},
}

//适配器 		4D36E964-E325-11CE-BFC1-08002BE10318
var deviceAdapterGUI = windows.GUID{
	Data1: 0x4d36e964,
	Data2: 0xe325,
	Data3: 0x11ce,
	Data4: [8]byte{0xbf,0xc1,0x08,0x00,0x2b,0xe1,0x03,0x18},
}

//系统设备 4D36E97D-E325-11CE-BFC1-08002BE10318
var deviceSystemGUI = windows.GUID{
	Data1: 0x4D36E97D,
	Data2: 0xe325,
	Data3: 0x11ce,
	Data4: [8]byte{0xbf,0xc1,0x08,0x00,0x2b,0xe1,0x03,0x18},
}


type DeviceInfo struct {
	DeviceDisplay 		string
}

func GetMouseDevice()(deviceList []*DeviceInfo){
	return list(&deviceMouseGUI)
}
func GetMediaDevice()(deviceList []*DeviceInfo){
	return list(&deviceMediaGUI)
}
func GetBlueDevice()(deviceList []*DeviceInfo){
	return list(&deviceBlueGUI)
}


func list(guid *windows.GUID) (deviceList []*DeviceInfo){
	devInfo, err := setup.SetupDiGetClassDevsEx(guid, "", 0, setup.DIGCF_PRESENT, setup.DevInfo(0), "")
	if err != nil {
		//fmt.Errorf("SetupDiGetClassDevsEx() failed: %v", err)
		return nil
	}
	defer devInfo.Close()
	for index := 0; ; index++ {
		b,deviceInfo := findProperty(devInfo, index)
		if deviceInfo != nil{
			deviceList = append(deviceList,deviceInfo)
		}
		if !b{
			break
		}
	}
	return
}

func findProperty(devInfo setup.DevInfo,index int) (bool,*DeviceInfo){
	defer func() {
		recover()
	}()
	devInfoData, err := devInfo.EnumDeviceInfo(index)
	if err != nil {
		if err == windows.ERROR_NO_MORE_ITEMS {
			return false,nil
		}
		return true,nil
	}
	property, err := devInfo.GetDeviceRegistryProperty(devInfoData, setup.SPDRP_PHYSICAL_DEVICE_OBJECT_NAME)
	if err != nil{
		return true,nil
	}
	property, _ = devInfo.GetDeviceRegistryProperty(devInfoData, setup.SPDRP_DEVICEDESC)
	return true,&DeviceInfo{DeviceDisplay:fmt.Sprintf("%s",property)}
}
