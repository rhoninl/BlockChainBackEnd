package Utils

type Account struct {
	Account   string `json:"account"`
	Password  string `json:"password"`
	CompanyId int64  `json:"companyId"`
}
type Info struct {
	CompanyInfo
	AddressInfo
	CompanyBasicInfo
}

type CompanyInfo struct {
	CompanyId int64  `json:"-"`
	AddressId int64  `json:"addressId,omitempty"`
	Phone     string `json:"phone"`
	Email     string `json:"email"`
}

type AddressInfo struct {
	Country string `json:"country"`
	City    string `json:"city"`
	Address string `json:"address"`
}

type Staff struct {
	StaffId   int64  `json:"staffId"`
	StaffName string `json:"staffName"`
	StaffJob  string `json:"staffJob"`
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
	Category    string  `json:"category"`
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
	MessageType int64  `json:"messageType"`
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

type StaffInfo struct {
	StaffId  int64  `json:"staffId"`
	Sex      string `json:"sex"`
	Phone    string `json:"phone"`
	Email    string `json:"email"`
	Fax      string `json:"fax"`
	JoinDate string `json:"joinDate"`
	BirthDay string `json:"birthDay"`
	AddressInfo
}

type MessageStruct struct {
	Message
	FromId int64 `json:"fromId"`
	ToId   int64 `json:"toId"`
}

type AskPrice struct {
	TargetCompanyId int64 `json:"toCompanyId"`
	OrderId         int64 `json:"orderId"`
}

type Bargain struct {
	Price  int64 `json:"price"`
	Status int64 `json:"status"`
	CompanyBasicInfo
}

type ReplyBargain struct {
	IsPass  bool  `json:"isPass"`
	Bargain int64 `json:"bargain"`
	OrderId int64 `json:"orderId"`
}
