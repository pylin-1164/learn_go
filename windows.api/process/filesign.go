package process

/*
#include <windows.h>
#include <wincrypt.h>
#include <tchar.h>
#include <stdio.h>
#include <string.h>
//#pragma comment(lib, "crypt32.lib")
#define ENCODING (X509_ASN_ENCODING | PKCS_7_ASN_ENCODING)

BOOL SignPath(){
	char szFileName[MAX_PATH];
    BOOL bIsSuccess;
    DWORD dwEncoding, dwContentType, dwFormatType;
    HCERTSTORE hStore = NULL;
    HCRYPTMSG hMsg = NULL;
    PVOID pvContext = NULL;
	strcpy(szFileName, "C:\\WINDOWS\\System32\\lsass.exe");
    bIsSuccess = CryptQueryObject(CERT_QUERY_OBJECT_FILE,
        szFileName,
        CERT_QUERY_CONTENT_FLAG_ALL,
        CERT_QUERY_FORMAT_FLAG_ALL,
        0,
        &dwEncoding,
        &dwContentType,
        &dwFormatType,
        &hStore,
        &hMsg,
        (const void**)&pvContext);
}
*/


func GetFileSign(){
}