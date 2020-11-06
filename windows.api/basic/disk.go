package basic

import (
	"fmt"
	"github.com/shirou/gopsutil/disk"
)

func DiskInfo() (usages []*disk.UsageStat){
	partitions, err := disk.Partitions(false)
	if err != nil{
		fmt.Println("Query Disk Info failed ...")
	}
	for _,partition := range partitions {
		usage, err := disk.Usage(partition.Mountpoint)
		if err != nil{
			continue
		}
		usages = append(usages,usage)
	}
	return
}