package process

import (
	"fmt"
	"github.com/gonutz/w32"
)

type FileVersionInfo struct {
	Company 		string
	Version			string
	Translations	[]string
}

func QueryFileInfo(path string)*FileVersionInfo{
	defer func() {
		recover()
	}()
	size := w32.GetFileVersionInfoSize(path)
	if size <= 0 {
		//fmt.Println("GetFileVersionInfoSize failed")
	}

	info := make([]byte, size)
	ok := w32.GetFileVersionInfo(path, info)

	if !ok {
		panic("GetFileVersionInfo failed")
	}

	fixed, ok := w32.VerQueryValueRoot(info)
	if !ok {
		//fmt.Println("VerQueryValueRoot failed")
	}
	versionBuff := fixed.FileVersion()
	version := fmt.Sprintf(
		"%d.%d.%d.%d",
		versionBuff&0xFFFF000000000000>>48,
		versionBuff&0x0000FFFF00000000>>32,
		versionBuff&0x00000000FFFF0000>>16,
		versionBuff&0x000000000000FFFF>>0,
	)

	translations, ok := w32.VerQueryValueTranslations(info)
	if !ok {
		//fmt.Println("VerQueryValueTranslations failed")
	}
	if len(translations) == 0 {
		//fmt.Println("no translation found")
	}

	t := translations[0]
	// w32.CompanyName simply translates to "CompanyName"
	company, ok := w32.VerQueryValueString(info, t, w32.CompanyName)
	if !ok {
		//fmt.Println("cannot get company name")
	}
	return &FileVersionInfo{
		Company:      company,
		Version:      version,
		Translations: translations,
	}
}