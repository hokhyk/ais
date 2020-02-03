package service

//TimeFormat 時間格式
var TimeFormat = "2006-01-02 15:04:05"

//Curl 物件介面
type Curl interface {
	Post(url string, params []byte, header map[string]string) []byte
}

//Redis 物件介面
type Redis interface {
	Get(key string) string
	Set(key, value string) bool
}

//Zip 物件介面
type Zip interface {
	Compress(files []string, dst string) int
}
