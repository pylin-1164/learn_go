package main

import (
	"fmt"
	"regexp"
	"windows.api/basic"
	"windows.api/device"
	"windows.api/foreground"
	"windows.api/net"
	"windows.api/process"
	"windows.api/winservices"
)

func main() {


	var RemoteServiceName = "Remote Desktop Services"
	var FirewallServiceName = "(Windows Firewall)|(Windows Defender Firewall)"
	var RemoteProcessName = `Windows\\System32\\rdpclip.exe`
	var ScreenSaverName = `\.scr$`
	//远程桌面服务状态，防火墙状态，正在使用远程桌面，正在使用屏保,锁屏中
	var remoteServiceStatus,
	firewallServiceStatus,
	remoteProcessStatus,
	ScreenSaverStatus,
	foregroundStatus = false,false,false,false,false


	//登录账号
	fmt.Printf("当前使用人登录账号 : %s \n",basic.GetUserName())

	fmt.Printf("===========================%s=======================================\n","Windows系统基本信息")
	//系统信息
	fmt.Printf("%+v \n",basic.SysInfo())

	fmt.Printf("===========================%s=======================================\n","Windows系统IP网络信息")
	//IPV4
	ipAddresses := basic.IpAddress()
	for _,addr := range ipAddresses {
		fmt.Printf("%+v \n",*addr)
	}

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
	serviceList := winservices.QueryServiceList()
	for _,serviceInfo := range serviceList {
		fmt.Printf("Service:%-40s  Name:%-50s Status:%-10s \n",serviceInfo.ServiceName,serviceInfo.ServiceDisplayName,serviceInfo.ServiceStatusCN)

		//判断远程桌面是否打开
		if RemoteServiceName == serviceInfo.ServiceDisplayName && serviceInfo.ServiceStatus == winservices.Service_Status_Running{
			remoteServiceStatus = true
		}
		//判断防火墙是否打开
		firewallMatch, _ := regexp.MatchString(FirewallServiceName, serviceInfo.ServiceDisplayName)
		if  firewallMatch && serviceInfo.ServiceStatus == winservices.Service_Status_Running{
			firewallServiceStatus = true
		}
	}


	fmt.Printf("===========================%s=======================================\n","Windows系统进程信息")
	//进程列表
	processList := process.List()
	for _,processInfo := range processList {
		fmt.Printf("%-120s    Version:%-20s    Company:%-70s \n",processInfo.ProcessName,processInfo.VersionInfo.Version,processInfo.VersionInfo.Company)
		if match, _ := regexp.MatchString(RemoteProcessName, processInfo.ProcessName);match{
			remoteProcessStatus = true
		}
		if match, _ := regexp.MatchString(ScreenSaverName, processInfo.ProcessName);match{
			ScreenSaverStatus = true
		}
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
		fmt.Printf("安装程序： %-70s 发行商：%-40s \n",applicationInfo.DisplayName,applicationInfo.Publisher)
	}

	foregroundStatus = foreground.LockStatusOpen()

	fmt.Printf("===========================%s=======================================\n","Windows防火墙、远程桌面状态")

	if firewallServiceStatus{
		fmt.Printf("防火墙状态：   开启   \n")
	}else{
		fmt.Printf("防火墙状态：   关闭   \n")
	}

	if remoteServiceStatus{
		fmt.Printf("远程桌面状态：   开启   \n")
	}else{
		fmt.Printf("远程桌面状态：   关闭   \n")
	}
	if remoteProcessStatus{
		fmt.Printf("正在使用远程桌面：  是   \n")
	}else{
		fmt.Printf("正在使用远程桌面：  否   \n")
	}
	if ScreenSaverStatus{
		fmt.Printf("屏保打开状态：  是   \n")
	}else{
		fmt.Printf("屏保打开状态：  否   \n")
	}
	if foregroundStatus{
		fmt.Printf("锁屏状态：  是   \n")
	}else{
		fmt.Printf("锁屏状态：  否   \n")
	}


	netStatus := net.NetWorkStatus()
	if netStatus{
		fmt.Printf("可以访问外网：  是   \n")
	}else{
		fmt.Printf("可以访问外网：  否   \n")
	}
}

