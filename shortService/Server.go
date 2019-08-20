/*************************************************************************
	> File Name: Server.go
	> Author: Wu Yinghao
	> Mail: wyh817@gmail.com
	> Created Time: 日  6/14 16:00:54 2015
 ************************************************************************/

package main

import (
	"flag"
	"fmt"
	"net/http"
	"os"
	"github.com/Maginawork/shortService/src/shortlib"
)

func main() {
	var configFile string
	flag.StringVar(&configFile, "conf", "config.ini", "configure file full path")
	flag.Parse()

	//读取配置文件
	fmt.Printf("[INFO] Read configure file...\n")
	configure, err := shortlib.NewConfigure(configFile)
	if err != nil {
		fmt.Printf("[ERROR] Parse Configure File Error: %v\n", err)
		return
	}

	//启动Redis客户端
	fmt.Printf("[INFO] Start Redis Client...\n")
	redis_cli, err := shortlib.NewRedisAdaptor(configure)
	if err != nil {
		fmt.Printf("[ERROR] Redis init fail..\n")
		return
	}
	//是否初始化Redis计数器，如果为ture就初始化计数器
	if configure.GetRedisStatus() {
		err = redis_cli.InitCountService()
		if err != nil {
			fmt.Printf("[ERROR] Init Redis key count fail...\n")
		}
	}

	//不使用redis的情况下，启动短链接计数器
	count_channl := make(chan shortlib.CountChannl, 1000)
	go CountThread(count_channl)

	countfunction := shortlib.CreateCounter(configure.GetCounterType(), count_channl, redis_cli)
	//启动LRU缓存
	fmt.Printf("[INFO] Start LRU Cache System...\n")
	lru, err := shortlib.NewLRU(redis_cli)
	if err != nil {
		fmt.Printf("[ERROR]LRU init fail...\n")
	}
	//初始化两个短连接服务
	fmt.Printf("[INFO] Start Service...\n")
	baseprocessor := &shortlib.BaseProcessor{redis_cli, configure, configure.GetHostInfo(), lru, countfunction}

	original := &OriginalProcessor{baseprocessor, count_channl}
	short := &ShortProcessor{baseprocessor}

	//启动http handler
	router := &shortlib.Router{configure, map[int]shortlib.Processor{
		0: short,
		1: original,
	}}

	//启动服务

	port, _ := configure.GetPort()
	addr := fmt.Sprintf(":%d", port)
	fmt.Printf("[INFO]Service Starting addr :%v,port :%v\n", addr, port)
	err = http.ListenAndServe(addr, router)
	if err != nil {
		//logger.Error("Server start fail: %v", err)
		os.Exit(1)
	}

}
