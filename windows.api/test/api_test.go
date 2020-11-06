package test

import (
	"fmt"
	"net"
	"regexp"
	"testing"
	"time"
	"windows.api/basic"
	"windows.api/foreground"
	"windows.api/process"
	"windows.api/winservices"
)

func TestAllProcess(t *testing.T){
	//远程桌面使用进程: rdpclip.exe
	//屏保： Ribbons.scr  -> *.scr
	process.List()

}

func TestForeground(t *testing.T){
	time.Sleep(30*time.Second)
	foreground.LockStatusOpen()
}

func TestSystemBasic(t *testing.T){
	basic.SysInfo()
}


func TestService(t *testing.T){
	//防火墙： Name: MpsSvc DisplayName: Windows Firewall Binary Path: C:\Windows\system32\svchost.exe -k LocalServiceNoNetwork State:  1
	//远程桌面: Name: TermService DisplayName: Remote Desktop Services Binary Path: C:\Windows\System32\svchost.exe -k NetworkService State:  4


	winservices.List()
}


func TestIp(t *testing.T){
	basic.IpAddress()
}


func TestMac(t *testing.T){
	addrs, err := net.InterfaceAddrs()
	if err != nil{
		return
	}
	for _,addr := range addrs {
		ipNet, isValidIpNet := addr.(*net.IPNet)
		if isValidIpNet && !ipNet.IP.IsLoopback() {
			if ipNet.IP.To4() != nil {
				fmt.Printf("%s , %s \n",ipNet.IP,ipNet.Mask)
			}

		}

	}
}

func TestRegrex(t *testing.T){
	fmt.Println(regexp.MatchString(process.ApplicationWhiteList,
		" Tools for .Net 3.5"))
}