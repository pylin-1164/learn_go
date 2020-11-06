package basic

import (
	"fmt"
	"github.com/shirou/gopsutil/cpu"
	"github.com/shirou/gopsutil/host"
	"strings"
	"time"
	"windows.api/api/iphlp"
)


func SysInfo() host.InfoStat{
	info, _ := host.Info()
	return *info
}

func CpuInfo() []cpu.InfoStat{
	info, err := cpu.Info()
	if err != nil{
		fmt.Println("get cpu data failed ... ")
	}
	return info
}


func CpuPercent() []float64{
	percent, err := cpu.Percent(time.Second, true)
	if err != nil{
		fmt.Println("get cpu percent failed ... ")
	}
	return percent
}

func systemIsWinXP() bool{
	sysInfo := SysInfo()
	if strings.Contains(sysInfo.Platform,"XP") {
		return true
	}
	return false
}


func IpAddress() []*iphlp.IpAddress{
	isWinXP := systemIsWinXP()
	addr := iphlp.GetAddresses(isWinXP) //XP下无法获取GateWway
	if isWinXP{
		ipv4List := iphlp.GetAllIpAddr()     //XP下无法获取DNS和Mac
		ipMap := make(map[string]iphlp.IpAddress)
		for _,ipv4 := range ipv4List {
			ipMap[ipv4.Ipv4] =  ipv4
		}
		for _,addr := range addr {
			xpAddress := ipMap[addr.Ipv4]
			addr.Gateway = xpAddress.Gateway
		}

	}
	return addr
}
