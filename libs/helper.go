package libs

//Helper 物件參數
type Helper struct{}

//New 建構式
func (h Helper) New() *Helper {
	return &h
}

//InArray 於陣列中有對應之值
func (h *Helper) InArray(arr []string, str string) bool {
	for _, v := range arr {
		if v == str {
			return true
		}
	}
	return false
}
