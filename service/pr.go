package service

import (
	"fmt"
	"io"
	"log"
	"mime/multipart"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/teed7334-restore/ais/dto"
	"github.com/teed7334-restore/ais/libs"
)

//PR 請購單資料結構
type PR struct{}

const (
	uploadPath = "./resources/uploads"
	proofPath  = "./resources/proof"
)

// New 建構式
func (pr PR) New() *PR {
	return &pr
}

//Add 新增請購單
func (pr *PR) Add(params *dto.PR, files *multipart.Form) int {
	status, uploads := pr.doUpload(files)
	if status != 1 {
		return status
	}
	ziper := libs.Zip{}.New()
	status, proof := pr.doCompress(uploads, ziper)
	if status != 1 {
		return status - 3
	}
	params.List.Proof = proof
	pr.doRemoveTempFiles(uploads)
	return 1
}

//doUpload 處理上傳檔案
func (pr *PR) doUpload(files *multipart.Form) (int, []string) {
	var dst []string
	for _, fs := range files.File {
		for i := range fs {
			file, err := fs[i].Open()
			defer file.Close()
			if err != nil { //未上傳佐証資料
				log.Println(err)
				return -1, dst
			}
			ext := strings.Split(fs[i].Filename, ".")[1]
			now := time.Now()
			filename := fmt.Sprintf("%d%d.%s", now.UnixNano(), i, ext)
			fp := filepath.Join(uploadPath, filename)
			filename = fmt.Sprintf("%s/%s", uploadPath, filename)
			dst = append(dst, filename)
			out, err := os.Create(fp)
			defer out.Close()
			if err != nil { //無法新增檔案
				log.Println(err)
				return -2, dst
			}
			_, err = io.Copy(out, file)
			if err != nil { //上傳資料夾權限不符
				log.Println(err)
				return -3, dst
			}
		}
	}
	return 1, dst
}

//doCompress 壓縮檔案
func (pr *PR) doCompress(files []string, ziper Zip) (int, string) {
	now := time.Now()
	filename := fmt.Sprintf("%d.zip", now.UnixNano())
	dst := fmt.Sprintf("%s/%s", proofPath, filename)
	status := ziper.Compress(files, dst)
	if status != 1 {
		return status, dst
	}
	return 1, dst
}

//doRemoveTempFiles 刪除暫存檔案
func (pr *PR) doRemoveTempFiles(files []string) {
	for _, file := range files {
		os.Remove(file)
	}
}
