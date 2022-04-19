package main

import (
	"github.com/sirupsen/logrus"
	_ "xiaohuazhu/internal/config"
	"xiaohuazhu/internal/server/http"
)

func main() {
	logrus.Infof("程序启动中...")
	http.New()
}
