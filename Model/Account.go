package Model

import (
	"log"
	"main/Utils"
)

func Login(account Utils.Account) (Utils.Account, bool, error) {
	var info Utils.Account
	template := `Select PassWord,CompanyId From ShippingTraceability.Account Where Account = ? limit 1`
	rows, err := Utils.DB().Query(template, account.Account)
	if err != nil {
		log.Println("[Login]make Mistake", err)
		return info, false, err
	}
	defer rows.Close()
	if !rows.Next() {
		return info, false, nil
	}
	rows.Scan(&info.Password, &info.CompanyId)
	return info, true, nil
}

func Info(companyId string) (Utils.CompanyInfo, bool, error) {
	var info Utils.CompanyInfo
	template := `Select Phone,AddressId, Email From ShippingTraceability.CompanyInfo Where CompanyId = ? limit 1`
	rows, err := Utils.DB().Query(template, companyId)
	if err != nil {
		log.Println("[Info]数据库发生异常", err)
		return info, false, err
	}
	defer rows.Close()
	if !rows.Next() {
		return info, false, nil
	}
	addressId := ""
	rows.Scan(&info.Phone, &addressId, &info.Email)
	// 获取地址信息
	addressInfo, err := QueryAddress(addressId)
	if err != nil {
		return info, false, err
	}
	info.Country = addressInfo.Country
	info.City = addressInfo.City
	info.Address = addressInfo.Address
	info.EnglishAddress = addressInfo.EnglishAddress
	//获取 基础信息
	basicInfo, err := CompanyBasicInfo(companyId)
	if err != nil {
		return info, false, err
	}
	info.CompanyName = basicInfo.CompanyName
	info.CompanyType = basicInfo.CompanyType
	info.CompanyId = companyId
	return info, true, nil

}

func QueryAddress(addressId string) (Utils.AddressInfo, error) {
	var addressInfo Utils.AddressInfo
	template := `Select Country, City, Address, EnglishAddress From ShippingTraceability.Address Where AddressId = ? Limit 1`
	rows, err := Utils.DB().Query(template, addressId)
	if err != nil {
		log.Println("[QueryAddress]数据库发生异常", err)
		return addressInfo, err
	}
	defer rows.Close()
	rows.Next()
	rows.Scan(&addressInfo.Country, &addressInfo.City, &addressInfo.Address, &addressInfo.EnglishAddress)
	return addressInfo, nil
}

func CompanyBasicInfo(companyId string) (Utils.CompanyBasicInfo, error) {
	var companyInfo Utils.CompanyBasicInfo
	template := `Select CompanyName, CompanyType From ShippingTraceability.Company Where CompanyId = ? Limit 1`
	rows, err := Utils.DB().Query(template, companyId)
	if err != nil {
		log.Println("[GetCompanyById]数据库发生异常", err)
		return companyInfo, err
	}
	defer rows.Close()
	rows.Next()
	rows.Scan(&companyInfo.CompanyName, &companyInfo.CompanyType)
	return companyInfo, nil
}
