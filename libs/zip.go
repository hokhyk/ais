package libs

import (
	"archive/zip"
	"io"
	"log"
	"os"
	"strings"
)

//Zip 物件參數
type Zip struct{}

//New 建構式
func (z Zip) New() *Zip {
	return &z
}

//Compress 壓縮檔案
func (z *Zip) Compress(files []string, dst string) int {
	newZip, err := os.Create(dst)
	if err != nil { //無法新增Zip壓縮檔
		log.Println(err)
		return -1
	}
	defer newZip.Close()
	zipWriter := zip.NewWriter(newZip)
	defer zipWriter.Close()
	for _, file := range files {
		zipFile, err := os.Open(file)
		if err != nil { //伺服器中無對應之Zip壓縮檔
			log.Println(err)
			return -2
		}
		defer zipFile.Close()
		info, err := zipFile.Stat()
		if err != nil { //無法解析Zip壓縮檔
			log.Println(err)
			return -3
		}
		header, err := zip.FileInfoHeader(info)
		if err != nil { //無法取得Zip壓縮檔檔頭資訊
			log.Println(err)
			return -4
		}
		paths := strings.Split(file, "/")
		filename := paths[len(paths)-1]
		header.Name = filename
		header.Method = zip.Deflate
		writer, err := zipWriter.CreateHeader(header)
		if err != nil { //無法建立Zip壓縮檔檔頭
			log.Println(err)
			return -5
		}
		_, err = io.Copy(writer, zipFile)
		if err != nil { //無法寫入檔到到Zip壓縮檔
			log.Println(err)
			return -6
		}
	}
	return 1
}
