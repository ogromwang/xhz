package oss

import (
	"testing"
)

func TestMinioClient(t *testing.T) {
	MinioClient()
}

func TestPushObject(t *testing.T) {
	client := MinioClient()
	client.PushObject("/Users/ogromwang/Downloads/test.png", "icon")
}