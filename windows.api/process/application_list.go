package process

import (
	mapset "github.com/deckarep/golang-set"
	"golang.org/x/sys/windows/registry"
	"regexp"
	"strconv"
)

const (
	ApplicationWhiteList = `(Update for Microsoft .*)`+
					`|(Microsoft.+\(KB\d{1,}\).+Edition)`+
					`|(Security Update.+\(KB\d{1,}\).+Edition)`+
					`|(Microsoft Visual.+)`+
					`|(Microsoft .+ Update$)`+
					`|(Help 更新 \(KB.+)`+
					`|(Microsoft Office.+ MUI.+)`+
					`|(Microsoft Office .+Chinese|English)`+
					`|(Microsoft .+ MUI.+)`+
					`|( 语言包|简体中文|语言服务|管理对象|加载项.{0,})`+
					`|(Microsoft .{0,}\.(NET|Net) .+)`+
					`|( SDK.{0,})`+
					`|(^(SDK|WinRT|Kits|Universal|Windows App|Application Verifier|Visual|Microsoft Portable Library|Update for|Windows Software Development Kit) .+)`+
					`|(^(x86\_64|i686).+)`+
					`|(NVIDIA Display .+)`+
					`|(Python .+Test|Core|Executables|Path|Support|Documentation|Utility)`+
					`|( (Library|Libraries) .+)`+
					`|(Microsoft .+ Components .+)`+
					`|(Visual C\+\+.+)`
	Win32UninstllPath = `SOFTWARE\Wow6432Node\Microsoft\Windows\CurrentVersion\Uninstall\`
	win64UninstallPath = `SOFTWARE\Microsoft\Windows\CurrentVersion\Uninstall\`
)

type ApplicationInfo struct {
	Publisher 		string
	Version 		string
	DisplayName 	string
}


func GetApplicationList() (applicationInfoList []ApplicationInfo){
	nameSet := mapset.NewSet()
	if strconv.IntSize == 64 {
		applicationInfoList = append(applicationInfoList,queryRegistList(Win32UninstllPath,&nameSet) ...)
		applicationInfoList = append(applicationInfoList,queryRegistList(win64UninstallPath,&nameSet) ...)

	}else {
		applicationInfoList = append(applicationInfoList,queryRegistList(win64UninstallPath,&nameSet) ...)
	}

	return

}



func queryRegistList(path string,nameSet *mapset.Set) (appList []ApplicationInfo){
	appUinstallKey, err := registry.OpenKey(registry.LOCAL_MACHINE, path, registry.ALL_ACCESS)
	defer appUinstallKey.Close()
	if err != nil {
		panic(err)
	}
	keyNames, err := appUinstallKey.ReadSubKeyNames(-1)
	if err != nil {
		panic(err)
	}
	for _, keyName := range keyNames {
		subKey, err := registry.OpenKey(registry.LOCAL_MACHINE, path+keyName, registry.ALL_ACCESS)
		defer subKey.Close()
		if err != nil {
			continue
		}

		applicationInfo := ApplicationInfo{}
		if value, _, err := subKey.GetStringValue("DisplayName"); err == nil {
			applicationInfo.DisplayName = value
		}
		//应用名称为空则跳过
		if applicationInfo.DisplayName == ""{
			continue
		}


		//白名单不记录
		if mt,_ := regexp.MatchString(ApplicationWhiteList,applicationInfo.DisplayName);mt{
			continue
		}

		//如果名称重复，跳过
		if !(*nameSet).Add(applicationInfo.DisplayName){
			continue
		}

		if value, _, err := subKey.GetStringValue("Publisher"); err == nil {
			applicationInfo.Publisher =  value
		}
		if value, _, err := subKey.GetStringValue("DisplayVersion"); err == nil {
			applicationInfo.Version = value
		}

		appList = append(appList,applicationInfo)
	}
	return
}