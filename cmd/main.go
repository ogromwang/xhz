package main

import (
	_ "xiaohuazhu/internal/config"
	"xiaohuazhu/internal/server/http"
)

func main() {

	http.New()
}