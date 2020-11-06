package main

import (
	"fmt"
	"windows.api/basic"
	"windows.api/device"
	"windows.api/process"
	"windows.api/winservices"
)

func main() {

	//登录账号
	fmt.Printf("当前使用人登录账号 : %s \n",basic.GetUserName())

	fmt.Printf("===========================%s=======================================\n","Windows系统基本信息")
	//系统信息
	fmt.Printf("%+v \n",basic.SysInfo())

	fmt.Printf("===========================%s=======================================\n","Windows系统CPU信息")
	cpuInfos := basic.CpuInfo()
	for _,cpuInfo := range cpuInfos {
		fmt.Printf("CPU modelName : %s   核心数: %d \n",cpuInfo.ModelName,cpuInfo.Cores)
	}
	cpuPercents := basic.CpuPercent()
	for i,cpuPercent := range cpuPercents {
		fmt.Printf("CPU Core[%d]: %.2f%s  ",i,cpuPercent,"%")
	}
	fmt.Println("")

	fmt.Printf("===========================%s=======================================\n","Windows系统磁盘信息")
	diskInfos := basic.DiskInfo()
	var  gbSize  int64 = 1024*1024*1024
	for _,diskInfo := range diskInfos {
		total := int64(diskInfo.Total)
		free := int64(diskInfo.Free)
		fmt.Printf("disk 序号:%s  总空间:%dGB  剩余空间%dGB 使用率%.2f \n",diskInfo.Path,total/gbSize,free/gbSize,diskInfo.UsedPercent)
	}



	fmt.Printf("===========================%s=======================================\n","Windows服务列表")
	//windows服务列表
	serviceList := winservices.QueryServiceStatus()
	for _,serviceInfo := range serviceList {
		fmt.Printf("Service:%s  Name:%s Status:%s \n",serviceInfo.ServiceName,serviceInfo.ServiceDisplayName,serviceInfo.ServiceStatusCN)
	}

	fmt.Printf("===========================%s=======================================\n","Windows系统IP网络信息")
	//IPV4
	ipAddresses := basic.IpAddress()
	for _,addr := range ipAddresses {
		fmt.Printf("%+v \n",*addr)
	}

	fmt.Printf("===========================%s=======================================\n","Windows系统进程信息")
	//进程列表
	processList := process.List()
	for _,processInfo := range processList {
		fmt.Printf("%s    Version:%s    Company:%s \n",processInfo.ProcessName,processInfo.VersionInfo.Version,processInfo.VersionInfo.Company)
	}


	fmt.Printf("===========================%s=======================================\n","Windows鼠标驱动列表")
	deviceList := device.GetMouseDevice()
	for _,deviceInfo := range deviceList {
		fmt.Printf("Mouse device : %s \n",deviceInfo.DeviceDisplay)
	}
	fmt.Printf("===========================%s=======================================\n","WindowsMedia驱动列表")
	deviceList = device.GetMediaDevice()
	for _,deviceInfo := range deviceList {
		fmt.Printf("Media device : %s \n",deviceInfo.DeviceDisplay)
	}

	fmt.Printf("===========================%s=======================================\n","Windows蓝牙驱动列表")
	deviceList = device.GetBlueDevice()
	for _,deviceInfo := range deviceList {
		fmt.Printf("Blue device : %s \n",deviceInfo.DeviceDisplay)
	}

	fmt.Printf("===========================%s=======================================\n","Windows安装程序列表")
	//查询安装程序列表
	applicationInfoList := process.GetApplicationList()
	for _,applicationInfo := range applicationInfoList {
		fmt.Printf("安装程序： %s 发行商：%s \n",applicationInfo.DisplayName,applicationInfo.Publisher)
	}
}

