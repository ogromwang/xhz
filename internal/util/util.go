package util

import (
	"fmt"
	"github.com/gofrs/uuid"
	"github.com/nfnt/resize"
	"image"
	_ "image/gif"
	"image/jpeg"
	_ "image/png"
	"os"
	"path"
)

func IntContains(array []int64, val int64) (index int) {
	index = -1
	for i := 0; i < len(array); i++ {
		if array[i] == val {
			index = i
			return
		}
	}
	return
}

// FormatFileSize 字节的单位转换 保留两位小数
func FormatFileSize(fileSize int64) (size string) {
	if fileSize < 1024 {
		//return strconv.FormatInt(fileSize, 10) + "B"
		size = fmt.Sprintf("%.2fB", float64(fileSize)/float64(1))
	} else if fileSize < (1024 * 1024) {
		size = fmt.Sprintf("%.2fKB", float64(fileSize)/float64(1024))
	} else if fileSize < (1024 * 1024 * 1024) {
		size = fmt.Sprintf("%.2fMB", float64(fileSize)/float64(1024*1024))
	} else if fileSize < (1024 * 1024 * 1024 * 1024) {
		size = fmt.Sprintf("%.2fGB", float64(fileSize)/float64(1024*1024*1024))
	} else if fileSize < (1024 * 1024 * 1024 * 1024 * 1024) {
		size = fmt.Sprintf("%.2fTB", float64(fileSize)/float64(1024*1024*1024*1024))
	} else { //if fileSize < (1024 * 1024 * 1024 * 1024 * 1024 * 1024)
		size = fmt.Sprintf("%.2fEB", float64(fileSize)/float64(1024*1024*1024*1024*1024))
	}
	return
}

func ImgFileResize(file, out *os.File, width uint) (err error) {
	img, _, err := image.Decode(file)
	// img, err := jpeg.Decode(file)
	if err != nil {
		return
	}

	// resize to width 1000 using Lanczos resampling
	// and preserve aspect ratio
	m := resize.Resize(width, 0, img, resize.Lanczos3)

	// write new image to file
	err = jpeg.Encode(out, m, nil)

	return
}

func NewImgTempPath(ext string) (temp *os.File, err error) {
	uu, _ := uuid.NewV4()
	if temp, err = os.CreateTemp("", uu.String()+"*"+ext); err != nil {
		return
	}
	return
}

func GetFileExt(fileName string) string {
	// 获取文件名带后缀
	filenameWithSuffix := path.Base(fileName)
	// 获取文件后缀
	return path.Ext(filenameWithSuffix)
}
