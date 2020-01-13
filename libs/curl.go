package libs

import (
	"io/ioutil"
	"log"
	"net/http"
	"strings"
)

//Curl 物件參數
type Curl struct {
}

//New 建構式
func (c Curl) New() *Curl {
	return &c
}

//Post 進行HTTP POST
func (c *Curl) Post(url string, params []byte, header map[string]string) []byte {
	data := &strings.Reader{}
	json := string(params)
	if json != "" {
		data = strings.NewReader(json)
	}
	req, err := http.NewRequest("POST", url, data)
	if err != nil {
		log.Println(err)
	}
	for k, v := range header {
		req.Header.Add(k, v)
	}
	clt := http.Client{}
	resp, err := clt.Do(req)
	if err != nil {
		log.Println(err)
	}
	defer req.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Println(err)
	}
	return body
}
