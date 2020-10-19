package setup

import (
	"encoding/binary"
	"fmt"
	"golang.org/x/sys/windows"
	"golang.org/x/sys/windows/registry"
	"runtime"
	"unsafe"
)

// SetupDiGetClassDevsEx function returns a handle to a device information set that contains requested device information elements for a local or a remote computer.
func SetupDiGetClassDevsEx(classGUID *windows.GUID, enumerator string, hwndParent uintptr, flags DIGCF, deviceInfoSet DevInfo, machineName string) (handle DevInfo, err error) {
	var enumeratorUTF16 *uint16
	if enumerator != "" {
		enumeratorUTF16, err = windows.UTF16PtrFromString(enumerator)
		if err != nil {
			return
		}
	}
	var machineNameUTF16 *uint16
	if machineName != "" {
		machineNameUTF16, err = windows.UTF16PtrFromString(machineName)
		if err != nil {
			return
		}
	}
	return setupDiGetClassDevsEx(classGUID, enumeratorUTF16, hwndParent, flags, deviceInfoSet, machineNameUTF16, 0)
}

// Close method deletes a device information set and frees all associated memory.
func (deviceInfoSet DevInfo) Close() error {
	return SetupDiDestroyDeviceInfoList(deviceInfoSet)
}

// EnumDeviceInfo method returns a DevInfoData structure that specifies a device information element in a device information set.
func (deviceInfoSet DevInfo) EnumDeviceInfo(memberIndex int) (*DevInfoData, error) {
	return SetupDiEnumDeviceInfo(deviceInfoSet, memberIndex)
}



//sys	setupDiEnumDeviceInfo(deviceInfoSet DevInfo, memberIndex uint32, deviceInfoData *DevInfoData) (err error) = setupapi.SetupDiEnumDeviceInfo

// SetupDiEnumDeviceInfo function returns a DevInfoData structure that specifies a device information element in a device information set.
func SetupDiEnumDeviceInfo(deviceInfoSet DevInfo, memberIndex int) (*DevInfoData, error) {
	data := &DevInfoData{}
	data.size = uint32(unsafe.Sizeof(*data))

	return data, setupDiEnumDeviceInfo(deviceInfoSet, uint32(memberIndex), data)
}


// DeviceRegistryProperty method retrieves a specified Plug and Play device property.
func (deviceInfoSet DevInfo) GetDeviceRegistryProperty(deviceInfoData *DevInfoData, property SPDRP) (interface{}, error) {
	return SetupDiGetDeviceRegistryProperty(deviceInfoSet, deviceInfoData, property)
}



//sys	setupDiGetDeviceRegistryProperty(deviceInfoSet DevInfo, deviceInfoData *DevInfoData, property SPDRP, propertyRegDataType *uint32, propertyBuffer *byte, propertyBufferSize uint32, requiredSize *uint32) (err error) = setupapi.SetupDiGetDeviceRegistryPropertyW

// SetupDiGetDeviceRegistryProperty function retrieves a specified Plug and Play device property.
func SetupDiGetDeviceRegistryProperty(deviceInfoSet DevInfo, deviceInfoData *DevInfoData, property SPDRP) (value interface{}, err error) {
	reqSize := uint32(256)
	for {
		var dataType uint32
		buf := make([]byte, reqSize)
		err = setupDiGetDeviceRegistryProperty(deviceInfoSet, deviceInfoData, property, &dataType, &buf[0], uint32(len(buf)), &reqSize)
		if err == windows.ERROR_INSUFFICIENT_BUFFER {
			continue
		}
		if err != nil {
			return
		}
		return getRegistryValue(buf[:reqSize], dataType)
	}
}


func getRegistryValue(buf []byte, dataType uint32) (interface{}, error) {
	switch dataType {
	case windows.REG_SZ:
		ret := windows.UTF16ToString(bufToUTF16(buf))
		runtime.KeepAlive(buf)
		return ret, nil
	case windows.REG_EXPAND_SZ:
		ret, err := registry.ExpandString(windows.UTF16ToString(bufToUTF16(buf)))
		runtime.KeepAlive(buf)
		return ret, err
	case windows.REG_BINARY:
		return buf, nil
	case windows.REG_DWORD_LITTLE_ENDIAN:
		return binary.LittleEndian.Uint32(buf), nil
	case windows.REG_DWORD_BIG_ENDIAN:
		return binary.BigEndian.Uint32(buf), nil
	case windows.REG_MULTI_SZ:
		bufW := bufToUTF16(buf)
		a := []string{}
		for i := 0; i < len(bufW); {
			j := i + wcslen(bufW[i:])
			if i < j {
				a = append(a, windows.UTF16ToString(bufW[i:j]))
			}
			i = j + 1
		}
		runtime.KeepAlive(buf)
		return a, nil
	case windows.REG_QWORD_LITTLE_ENDIAN:
		return binary.LittleEndian.Uint64(buf), nil
	default:
		return nil, fmt.Errorf("Unsupported registry value type: %v", dataType)
	}
}


// bufToUTF16 function reinterprets []byte buffer as []uint16
func bufToUTF16(buf []byte) []uint16 {
	sl := struct {
		addr *uint16
		len  int
		cap  int
	}{(*uint16)(unsafe.Pointer(&buf[0])), len(buf) / 2, cap(buf) / 2}
	return *(*[]uint16)(unsafe.Pointer(&sl))
}

// utf16ToBuf function reinterprets []uint16 as []byte
func utf16ToBuf(buf []uint16) []byte {
	sl := struct {
		addr *byte
		len  int
		cap  int
	}{(*byte)(unsafe.Pointer(&buf[0])), len(buf) * 2, cap(buf) * 2}
	return *(*[]byte)(unsafe.Pointer(&sl))
}


func wcslen(str []uint16) int {
	for i := 0; i < len(str); i++ {
		if str[i] == 0 {
			return i
		}
	}
	return len(str)
}