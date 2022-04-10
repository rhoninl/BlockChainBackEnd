package Utils

type Account struct {
	Account   string `json:"account"`
	Password  string `json:"password"`
	CompanyId string `json:"companyId"`
}

type CompanyInfo struct {
	CompanyId string `json:"companyId"`
	Phone     string `json:"phone"`
	Email     string `json:"email"`
	AddressInfo
	CompanyBasicInfo
}

type AddressInfo struct {
	Country        string `json:"country"`
	City           string `json:"city"`
	Address        string `json:"address"`
	EnglishAddress string `json:"englishAddress"`
}

type Stuff struct {
	StuffId   string `json:"stuffId"`
	StuffName string `json:"stuffName"`
	StuffJob  string `json:"stuffJob"`
}

type CompanyBasicInfo struct {
	CompanyId   string `json:"companyId"`
	CompanyName string `json:"companyName"`
	CompanyType string `json:"companyType"`
}
