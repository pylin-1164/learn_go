package foreground

import (
	"log"
	"syscall"
)

func LockStatusOpen() bool{
	const successCallMessage = "The operation completed successfully."
	// 加载类库
	user32 := syscall.NewLazyDLL("user32.dll")
	// 创建新的调用进程
	getForegroundWindow := user32.NewProc("GetForegroundWindow")
	// 调用相应的函数
	activeWindowId, _, err := getForegroundWindow.Call()

	if err != nil && err.Error() != successCallMessage {
		return false
		log.Println(err)
	}
	log.Println("activeWindowId:", activeWindowId)
	if activeWindowId == 0{
		return true
	}
	return false

}
