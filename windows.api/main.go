package main

import (
	"fmt"
	"windows.api/basic"
)

func main() {
	//winservices.QueryServiceStatus()
	ipAddresses := basic.IpAddress()
	for _,addr := range ipAddresses {
		fmt.Printf("%+v \n",*addr)
	}


}

