package main

import (
	"encoding/json"
	"fmt"
	"os"
	"path"
	"path/filepath"
)

type configuration struct {
	Host string
}

func readConfig() configuration {
	// 打开文件
	ex, _ := os.Executable()
	exPath := filepath.Dir(ex)
	confPath := path.Join(exPath, "cromwell-om.conf")
	file, _ := os.Open(confPath)
	conf := configuration{}

	// 关闭文件
	defer file.Close()

	//NewDecoder创建一个从file读取并解码json对象的*Decoder，解码器有自己的缓冲，并可能超前读取部分json数据。
	decoder := json.NewDecoder(file)

	//Decode从输入流读取下一个json编码值并保存在v指向的值里
	err := decoder.Decode(&conf)
	if err != nil {
		fmt.Println("Error:", err)
	}
	return conf
}
