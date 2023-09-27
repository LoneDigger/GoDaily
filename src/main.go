package main

import (
	"embed"
	"flag"

	"github.com/gin-gonic/gin"
	"me.daily/src/service"
)

//go:embed public/*
var fs embed.FS

var host, user, password, dbname, authUser, authPw string

func init() {
	flag.StringVar(&host, "host", "", "database host")
	flag.StringVar(&user, "user", "", "database user")
	flag.StringVar(&password, "password", "", "database password")
	flag.StringVar(&dbname, "dbname", "", "database dbname")
	flag.StringVar(&authUser, "authUser", "", "auth user")
	flag.StringVar(&authPw, "authPw", "", "auth password")
}

func main() {
	flag.Parse()

	gin.SetMode(gin.ReleaseMode)
	service.NewService(host, user, password, dbname, authUser, authPw, fs).Start()
}
