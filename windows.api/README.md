
#项目打包命令

```shell script
 go build -ldflags="-linkmode internal" -o winEnvListen.exe
```

#项目开发编译过程问题记录

>   CGO代码使用IED工具运行报错

```shell script
# 添加启动参数：
-ldflags="-linkmode internal"
```

> CGO代码返回数据类型转换为go string
```go
// From https://golang.org/cmd/cgo/
// A few special functions convert between Go and C types by making copies of the data. 
// Go string to C string
// The C string is allocated in the C heap using malloc.
// It is the caller's responsibility to arrange for it to be
// freed, such as by calling C.free (be sure to include stdlib.h
// if C.free is needed).

func C.CString(string) *C.char

// Go []byte slice to C array
// The C array is allocated in the C heap using malloc.
// It is the caller's responsibility to arrange for it to be
// freed, such as by calling C.free (be sure to include stdlib.h
// if C.free is needed).
func C.CBytes([]byte) unsafe.Pointer

// C string to Go string
func C.GoString(*C.char) string

// C data with explicit length to Go string
func C.GoStringN(*C.char, C.int) string

// C data with explicit length to Go []byte
func C.GoBytes(unsafe.Pointer, C.int) []byte
```