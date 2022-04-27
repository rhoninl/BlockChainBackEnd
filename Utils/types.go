package Utils

type Account struct {
	Account   string `json:"account"`
	Password  string `json:"password"`
	CompanyId int64  `json:"companyId"`
}

type CompanyInfo struct {
	//CompanyId string `json:"companyId"`
	Phone string `json:"phone"`
	Email string `json:"email"`
	AddressInfo
	CompanyBasicInfo
}

type AddressInfo struct {
	Country string `json:"country"`
	City    string `json:"city"`
	Address string `json:"address"`
}

type Stuff struct {
	StuffId   int64  `json:"stuffId"`
	StuffName string `json:"stuffName"`
	StuffJob  string `json:"stuffJob"`
}

type CompanyBasicInfo struct {
	CompanyId   int64  `json:"companyId"`
	CompanyName string `json:"companyName"`
	CompanyType string `json:"companyType"`
}

type Order struct {
	OrderId              int64  `json:"orderId"`
	ClientCompanyName    string `json:"clientCompanyName"`
	StartDate            string `json:"startDate"`
	LandTransCompanyName string `json:"landTransCompanyName"`
	SeaTransCompanyName  string `json:"seaTransCompanyName"`
	Status               string `json:"status"`
}

type CompanyList struct {
	CompanyId   int64  `json:"companyId"`
	CompanyName string `json:"companyName"`
	CompanyType string `json:"companyType"`
}

type RegisterInfo struct {
	AuthCode
	Account
	CompanyBasicInfo
}

type OrderInfo struct {
	OrderId     int64  `json:"orderId"`
	Phone       string `json:"phone"`
	Email       string `json:"email"`
	Fax         string `json:"fax"`
	UnStackable bool   `json:"unStackable"`
	Perishable  bool   `json:"perishable"`
	Dangerous   bool   `json:"dangerous"`
	Clearance   bool   `json:"clearance"`
	Incoterms   string `json:"incoterms"`
	Other       string `json:"other"`

	SendAddress    AddressInfo `json:"sendAddress"`
	ReceiveAddress AddressInfo `json:"receiveAddress"`

	ClientCompanyId int64 `json:"clientCompanyId"`

	Cargos []Cargo `json:"cargos"`

	DeliveryDate  string `json:"deliveryDate"`
	HopeReachDate string `json:"hopeReachDate"`
}

type Cargo struct {
	CargoId     int64   `json:"cargoId"`
	CargoName   string  `json:"cargoName"`
	CargoModel  string  `json:"cargoModel"`
	CargoSize   string  `json:"cargoSize"`
	CargoNum    int64   `json:"cargoNum"`
	CategoryId  int64   `json:"categoryId"`
	CargoWeight float64 `json:"cargoWeight"`
}

type AuthCode struct {
	ToEmail string `json:"email"`
	Code    string `json:"authCode"`
}

type ForgetPasswordForm struct {
	Account string `json:"account"`
	Email   string `json:"email"`
}

type Message struct {
	Context string `json:"context"`
}

type MessageList struct {
	MessageId   int64  `json:"messageId"`
	CompanyName string `json:"companyName"`
	MessageType string `json:"messageType"`
	SendTime    string `json:"sendTime"`
	IsRead      int    `json:"isRead"`
	IsReply     int    `json:"isReply"`
}

type GetAuth struct {
	Email string `json:"email"`
	Tag   string `json:"tag"`
}

type ReplyFriend struct {
	CompanyId int64 `json:"companyId"`
	ToId      int64 `json:"toId"`
	MessageId int64 `json:"messageId"`
	Ok        bool  `json:"ok"`
}

type MessageInfo struct {
	MessageId   int64 `json:"messageId"`
	FromId      int64 `json:"fromId"`
	ToId        int64 `json:"toId"`
	MessageType int64 `json:"messageType"`
	IsReply     int64 `json:"isReply"`
}
