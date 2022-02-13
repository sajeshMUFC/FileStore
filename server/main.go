package main

//
// HIGHLY CONFIDENTIAL
//
// Copyright (c) 2021 Wabtec Corporation. All rights reserved.
//
// NOTICE: This file contains material that is confidential and proprietary to
// Wabtec and/or other developers and may be protected by patents, copyright, and/or
// trade secrets. No license is granted under any intellectual property rights
// of Wabtec except as may be provided in a duly executed agreement with Wabtec.
// Any unauthorized reproduction or distribution of material from this file is
// expressly prohibited. This source code is and remains the property of Wabtec.
//

import (
	"fmt"
	"http-filestore/httpd/handler"
	filestore "http-filestore/platform/file-ops"
	"log"
	"os"
	"strconv"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

type App struct {
	addr string
	port int
}

func main() {
	log.Println("starting file store http server")

	a, err := NewApp()
	a.port = 8000
	if err != nil {
		log.Fatal(err)
	}
	defer a.Close()

	router := gin.Default()
	config := cors.DefaultConfig()
	// Replace * wth trusted domains
	config.AllowOrigins = []string{"*"}
	router.Use(cors.New(config))
	fstoreservice := filestore.NewFileStoreService()
	// max file upload limit
	router.MaxMultipartMemory = 8 << 20 // 8 MiB
	v1 := router.Group("/v1/file")
	{
		v1.POST("", handler.FileUploadHandler(fstoreservice))
		v1.GET("", handler.FileListHandler(fstoreservice))
		v1.DELETE("/:filename", handler.FileDeleteHandler(fstoreservice))
		v1.PUT("", handler.FileUpdateHandler(fstoreservice))
		v1.GET("/:search", handler.FileWordCountHandler(fstoreservice))
		v1.GET("/freqword", handler.FileFreqWordCountHandler(fstoreservice))
	}

	router.Run(fmt.Sprintf("%v:%v", a.addr, a.port))
}

func NewApp() (*App, error) {
	p, err := strconv.Atoi(os.Getenv("SERVE_PORT"))
	if err != nil {
		p = 8000
	}

	addr := os.Getenv("SERVE_ADDR")
	l := App{
		addr: addr,
		port: p,
	}
	return &l, nil
}

func (a *App) Close() {
	//close all connection while exit
}
