package Utils

type Account struct {
	Account   string `json:"account"`
	Password  string `json:"password"`
	CompanyId string `json:"companyId"`
}

type Stuff struct {
	StuffId   string `json:"stuffId"`
	StuffName string `json:"stuffName"`
	StuffJob  string `json:"stuffJob"`
}
