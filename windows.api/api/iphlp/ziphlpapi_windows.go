package iphlp

import (
	"fmt"
	"golang.org/x/sys/windows"
	"os"
	"strconv"
	"strings"
	"syscall"
	"unsafe"
)

/**
 * IP地址相关 windows  api 库
 * 参考 ： https://github.com/kbinani/win
 */

var(
	libiphlpapi = windows.NewLazySystemDLL("iphlpapi.dll")
	getIpAddrTable = libiphlpapi.NewProc("GetIpAddrTable")
	getAdaptersAddresses  = libiphlpapi.NewProc("GetAdaptersAddresses")
)

const (
	ANY_SIZE = 20

	/**Return both IPv4 and IPv6 addresses associated with adapters with IPv4 or IPv6 enabled.*/
	Family_AF_UNSPEC = 0
	/**Return only IPv4 addresses associated with adapters with IPv4 enabled.*/
	Family_AF_INET = 2
	/**Return only IPv6 addresses associated with adapters with IPv6 enabled.*/
	Family_AF_INET6 = 23

	/**Return a list of IP address prefixes on this adapter. When this flag is set, IP address prefixes are returned for both IPv6 and IPv4 addresses.
	This flag is supported on Windows XP with SP1 and later.*/
	GAA_FLAG_INCLUDE_PREFIX = 0x0010
)

type MIB_IPADDRTABLE struct {
	DwNumEntries uint32
	Table        [ANY_SIZE]MIB_IPADDRROW_XP
}

type MIB_IPADDRROW_XP struct {
	DwAddr      uint32
	DwIndex     uint32
	DwMask      uint32
	DwBCastAddr uint32
	DwReasmSize uint32
	Unused1     uint16
	WType       uint16
}

type IpAddress struct {
	Ipv4		string
	Ipv6 		string
	Mask		string
	Gateway 	string
	Mac 		string
	DNS			[]string
	Type 		int
}

func GetAddresses(isWinXP bool) (adds []*IpAddress){
	addresses, _ := adapterAddressesEx()
	for _,addr := range addresses {
		ipaddr := &IpAddress{}
		getIP(addr,ipaddr)
		if !isWinXP{
			getGateway(addr,ipaddr)
		}

		getDns(addr, ipaddr)

		getMac(addr, ipaddr)

		getType(ipaddr, addr)

		if ipaddr.Ipv4 == "127.0.0.1"{
			continue
		}

		if isWinXP{
			adds = append(adds, ipaddr)

		}else if ipaddr.Gateway != "" {
			/*if addr.IfType == windows.IF_TYPE_ETHERNET_CSMACD {
				fmt.Printf("Index: %d  Type:%s	 add: %+v \n", addr.IfIndex, "本地连接", ipaddr)
			}
			if addr.IfType == windows.IF_TYPE_IEEE80211 {
				fmt.Printf("Index: %d  Type:%s	 add: %+v \n", addr.IfIndex, "无线连接", ipaddr)
			}*/
			adds = append(adds, ipaddr)
		}

	}
	return adds
}

func getType(ipaddr *IpAddress, addr *IpAdapterAddresses) {
	ipaddr.Type = int(addr.IfType)
}

func getMac(addr *IpAdapterAddresses, ipaddr *IpAddress) {
	mac := fmt.Sprintf("%.2X-%.2X-%.2X-%.2X-%.2X-%.2X", addr.PhysicalAddress[0], addr.PhysicalAddress[1], addr.PhysicalAddress[2], addr.PhysicalAddress[3], addr.PhysicalAddress[4], addr.PhysicalAddress[5])
	ipaddr.Mac = mac
}

func getDns(addr *IpAdapterAddresses, ipaddr *IpAddress) {
	dnsServer := addr.FirstDnsServerAddress
	for {
		if dnsServer != nil {
			dns := dnsServer.Address.IP()
			dnsServer = dnsServer.Next
			ipaddr.DNS = append(ipaddr.DNS, dns.String())
		} else {
			break
		}
	}
}

func getGateway(addr *IpAdapterAddresses,ipaddr *IpAddress){
	gateway := addr.FirstGatewayAddress
	if gateway == nil || gateway.Address.IP() == nil{
		return
	}
	ipaddr.Gateway = gateway.Address.IP().String()
}

func getIP(addr *IpAdapterAddresses,ipaddr *IpAddress) {
	//fmt.Printf("%s \n",strconv.FormatUint(uint64(*(addr.FriendlyName)), 10))
	unicastAddress := addr.FirstUnicastAddress
	if unicastAddress == nil{
		return
	}
	for {
		ip := unicastAddress.Address.IP()
		if strings.Contains(ip.String(), ":") {
			ipaddr.Ipv6 = ip.String()
		} else {
			ipaddr.Ipv4 = ip.String()
			if ip.DefaultMask() != nil {
				parseUint, _ := strconv.ParseUint(ip.DefaultMask().String(), 16, 32)
				mask := uint32(parseUint)
				ipaddr.Mask = fmt.Sprintf("%d.%d.%d.%d", byte(mask>>24), byte(mask>>16), byte(mask>>8), byte(mask))
			}

		}
		if next := unicastAddress.Next; next != nil {
			unicastAddress = next
		} else {
			break
		}

	}
}

func GetAllIpAddr() (ips []IpAddress) {

	var b []byte
	l := uint32(1024) // recommended initial size
	b = make([]byte, l)
	a1, _, err := syscall.Syscall(getIpAddrTable.Addr(), 3,
		uintptr(unsafe.Pointer(&b[0])),
		uintptr(unsafe.Pointer(&l)),
		uintptr(1))
	if err != windows.NO_ERROR{
		return
	}
	if a1 != windows.NO_ERROR{
		fmt.Println("执行失败 : CODE = ",a1)
	}

	pIpAddrTable := (*MIB_IPADDRTABLE)(unsafe.Pointer(&b[0]))

	for i:=0;i<int(pIpAddrTable.DwNumEntries);i++ {
		table := pIpAddrTable.Table[i]
		ip := table.DwAddr
		mask := table.DwMask
		ips = append(ips,IpAddress{
			Ipv4: fmt.Sprintf("%d.%d.%d.%d", byte(ip),byte(ip>>8), byte(ip>>16), byte(ip>>24)),
			Gateway: fmt.Sprintf("%d.%d.%d.%d", byte(mask),byte(mask>>8), byte(mask>>16), byte(mask>>24)),
		})

	}
	return
}




type IpAdapterAddresses struct {
	Length                uint32
	IfIndex               uint32
	Next                  *IpAdapterAddresses
	AdapterName           *byte
	FirstUnicastAddress   *windows.IpAdapterUnicastAddress
	FirstAnycastAddress   *windows.IpAdapterAnycastAddress
	FirstMulticastAddress *windows.IpAdapterMulticastAddress
	FirstDnsServerAddress   *windows.IpAdapterDnsServerAdapter
	DnsSuffix             *uint16
	Description           *uint16
	FriendlyName          *uint16
	PhysicalAddress       [syscall.MAX_ADAPTER_ADDRESS_LENGTH]byte
	PhysicalAddressLength uint32
	Flags                 uint32
	Mtu                   uint32
	IfType                uint32
	OperStatus            uint32
	Ipv6IfIndex           uint32
	ZoneIndices           [16]uint32
	FirstPrefix           *windows.IpAdapterPrefix
	TransmitLinkSpeed	  uint64
	ReceiveLinkSpeed	  uint64
	FirstWinsServerAddress *IpAdapterGatewayAddress
	FirstGatewayAddress   *IpAdapterGatewayAddress
}

type IpAdapterGatewayAddress struct {
	Length             	  uint32
	Reserved              uint32
	Next    			  *IpAdapterGatewayAddress
	Address 			  windows.SocketAddress
}


func apiAdaptersAddresses(family uint32, flags uint32, reserved uintptr, adapterAddresses *IpAdapterAddresses, sizePointer *uint32) (errcode error) {
	r0, _, _ := syscall.Syscall6(getAdaptersAddresses.Addr(), 5, uintptr(family), uintptr(flags), uintptr(reserved), uintptr(unsafe.Pointer(adapterAddresses)), uintptr(unsafe.Pointer(sizePointer)), 0)
	if r0 != 0 {
		errcode = syscall.Errno(r0)
	}
	return
}

func adapterAddressesEx()([]*IpAdapterAddresses, error) {
	const GAA_FLAG_INCLUDE_GATEWAYS = 0x00000080
	var b []byte
	l := uint32(15000) // recommended initial size
	for {
		b = make([]byte, l)
		err := apiAdaptersAddresses(syscall.AF_UNSPEC, GAA_FLAG_INCLUDE_GATEWAYS, 0, (*IpAdapterAddresses)(unsafe.Pointer(&b[0])), &l)
		if err == nil {
			if l == 0 {
				return nil, nil
			}
			break
		}
		if err.(syscall.Errno) != syscall.ERROR_BUFFER_OVERFLOW {
			return nil, os.NewSyscallError("getadaptersaddresses", err)
		}
		if l <= uint32(len(b)) {
			return nil, os.NewSyscallError("getadaptersaddresses", err)
		}
	}
	var aas []*IpAdapterAddresses
	for aa := (*IpAdapterAddresses)(unsafe.Pointer(&b[0])); aa != nil; aa = aa.Next {
		aas = append(aas, aa)
	}
	return aas, nil
}