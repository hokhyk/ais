package controller

import (
	"fmt"
	"io"
	"net/http"
	"os"

	"github.com/teed7334-restore/ais/service"
)

//Download 下載檔案資料結構
type Download struct{}

var download = service.Download{}.New()

//New 建構式
func (d Download) New() *Download {
	return &d
}

//GetFile 下載檔案
func (d *Download) GetFile(w http.ResponseWriter, r *http.Request) {
	proof, ok := r.URL.Query()["proof"]
	if !ok {
		content := RO.BuildJSON(0, "檔名不得為空值")
		fmt.Fprintf(w, content)
		return
	}

	token := r.FormValue("token")

	fileName := fmt.Sprintf("./resources/proof/%s", proof[0])
	dtoRO := download.GetFile(token, fileName)

	if dtoRO.Status != 1 {
		PrintRO(w, dtoRO, "")
		return
	}

	resp, err := os.OpenFile(fileName, os.O_RDWR, 0666)
	if err != nil {
		content := RO.BuildJSON(0, "檔案無法開啟或是無此檔案")
		fmt.Fprintf(w, content)
		return
	}

	w.Header().Set("Content-Disposition", "attachment; filename="+proof[0])
	w.Header().Set("Content-Type", r.Header.Get("Content-Type"))
	io.Copy(w, resp)
}
