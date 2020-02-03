package dto

//Users 使用者相關資料結構
type Users struct {
	ID               int    `json:"id,string"`
	FirstName        string `json:"firstname"`
	LastName         string `json:"lastname"`
	Account          string `json:"login"`
	EMail            string `json:"email"`
	Role             int    `json:"role,string"`
	Manager          int    `json:"manager,string"`
	Country          string `json:"country"`
	OrganizationID   int    `json:"organization_id,string"`
	OrganizationName string `json:"organization_name"`
	Contract         string `json:"contract"`
	Position         int    `json:"position,string"`
	Identifier       string `json:"identifier"`
}
