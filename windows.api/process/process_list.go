package process

import (
	"context"
	"fmt"
	"github.com/shirou/gopsutil/process"
	"golang.org/x/sys/windows"
)

func List(){

	processes, err := process.Processes()
	if err != nil{
		fmt.Println(err)
	}
	fmt.Println("process list size : ", len(processes))
	for _, p := range processes {
		if p.Pid != 0 {
			if name, err := p.Name();err == nil{
				fmt.Printf("%s \n",name)
			}
		}
	}
}

//根据WindowsAPI自实现,兼容XP环境
func PidsWithContext(ctx context.Context) []uint32  {
	// inspired by https://gist.github.com/henkman/3083408
	// and https://github.com/giampaolo/psutil/blob/1c3a15f637521ba5c0031283da39c733fda53e4c/psutil/arch/windows/process_info.c#L315-L329
	const dwordSize uint32 = 4
	var read uint32 = 0
	//只枚举1024个进程
	ps := make([]uint32, 1024)
	if err := windows.EnumProcesses(ps, &read); err != nil {
		return make([]uint32,0)
	}

	for _, pid := range ps {
		if pid != 0 {
			p, _ := process.NewProcess(int32(pid))
			if name, err := p.Name();err == nil{
				fmt.Printf("%s \n",name)
			}
		}
	}

	return ps
}