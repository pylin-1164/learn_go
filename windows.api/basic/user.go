package basic

/*
#include <windows.h>
#include <Lmcons.h>
#include <stdio.h>
#include <winbase.h>

LPTSTR  QueryUserName() {
	const int MAX_BUFFER_LEN = UNLEN + 1;
	TCHAR szBuffer[MAX_BUFFER_LEN];
	DWORD dwNameLen;
	LPTSTR userName;
	dwNameLen = MAX_BUFFER_LEN;
	if(GetUserNameA(szBuffer, &dwNameLen))
		userName = szBuffer;
	return userName;
}
*/
import "C"

/**
 * 查询电脑登录账号
 */
func GetUserName()string{
	name := C.GoString(C.QueryUserName())
	return name
}