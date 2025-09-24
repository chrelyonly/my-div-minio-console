package chrelyonly

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"time"
)

func init() {
	//允许exe直接运行
	//cobra.MousetrapHelpText = ""
	//kernel32, _ := syscall.LoadLibrary(`kernel32.dll`)
	//sct, _ := syscall.GetProcAddress(kernel32, `SetConsoleTitleW`)
	//syscall.Syscall(sct, 1, uintptr(unsafe.Pointer(syscall.StringToUTF16Ptr("minio文件服务器 by_chrelyonly"))), 0, 0)
	//syscall.FreeLibrary(kernel32)
}

type Person struct {
	ApiUrl string `json:"apiUrl"`
}

func Optimize(args []string) []string {
	//打印当前时间
	currentTime := time.Now()
	fmt.Println("当前时间: " + currentTime.Format("2006-01-02 15:04:05"))
	var configPath = ""
	//判断是否传入配置文件路径
	if len(args) > 1 {
		configPath = args[1] + "/config.json"
	} else {
		//获取当前程序目录下的config.json文件
		configPath = "./config.json"
	}
	//判断文件是否存在
	_, err := os.Stat(configPath)
	fmt.Println("读取到配置文件,路径: " + configPath)
	//配置信息
	var apiUrl string
	//新的参数
	var newArgs []string
	//文件对象
	var config *os.File
	if err != nil {
		fmt.Println("当前配置文件不存在,将进行初始化配置")
		//创建文件
		config, err = os.Create(configPath)
		if err != nil {
			fmt.Println("配置文件创建失败,请检查是否有权限")
			os.Exit(1)
		}
		fmt.Println("请输入API地址(默认: http://127.0.0.1:30000))")
		_, err := fmt.Scanln(&apiUrl)
		if err != nil {
			apiUrl = "http://127.0.0.1:30000"
			fmt.Println("当前API地址: http://127.0.0.1:30000")
		}
		//将配置信息保存
		person := Person{
			ApiUrl: apiUrl,
		}
		//将配置信息写入文件
		encoder := json.NewEncoder(config)
		err = encoder.Encode(person)
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
	} else {
		//	如果存在文件则,从文件中读取配置
		fmt.Println("当前配置文件存在,将从配置文件中读取配置")
		//读取配置文件,json格式
		config, err = os.OpenFile(configPath, os.O_RDWR, 0666)
		if err != nil {
			fmt.Println("配置文件读取失败,请检查配置文件是否存在")
			os.Exit(1)
		}
		// 从文件中读取 JSON 数据
		decoder := json.NewDecoder(config)
		var person Person
		err = decoder.Decode(&person)
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
		apiUrl = person.ApiUrl
	}

	//关闭文件
	err = config.Close()
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	//apiUrl
	err = os.Setenv("CONSOLE_MINIO_SERVER", apiUrl)
	if err != nil {
		fmt.Println("apiUrl设置失败")
		os.Exit(1)
	}
	appName := filepath.Base(args[0])
	args = append(newArgs, appName, "server")
	return args
}
