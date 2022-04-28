package Model

import (
	"fmt"
	"golang.org/x/crypto/bcrypt"
	"log"
	"main/Utils"
	"reflect"
	"strconv"
	"time"
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

func Info(companyId int64) (Utils.CompanyInfo, bool, error) {
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
	var addressId int64
	rows.Scan(&info.Phone, &addressId, &info.Email)
	// 获取地址信息
	addressInfo, err := QueryAddress(addressId)
	if err != nil {
		return info, false, err
	}
	info.Country = addressInfo.Country
	info.City = addressInfo.City
	info.Address = addressInfo.Address
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

func QueryAddress(addressId int64) (Utils.AddressInfo, error) {
	var addressInfo Utils.AddressInfo
	template := `Select Country, City, Address From ShippingTraceability.Address Where AddressId = ? Limit 1`
	rows, err := Utils.DB().Query(template, addressId)
	if err != nil {
		log.Println("[QueryAddress]数据库发生异常", err)
		return addressInfo, err
	}
	defer rows.Close()
	rows.Next()
	rows.Scan(&addressInfo.Country, &addressInfo.City, &addressInfo.Address)
	return addressInfo, nil
}

func CompanyBasicInfo(companyId int64) (Utils.CompanyBasicInfo, error) {
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

func RegisterInfo(info Utils.RegisterInfo) (bool, error) {
	template := `Insert Into Company Set CompanyName = ?,CompanyType = ?`
	rows, err := Utils.DB().Exec(template, "未命名", "未选择")
	if err != nil {
		return false, err
	}
	companyId, err := rows.LastInsertId()
	if err != nil {
		return false, err
	}
	template = `Insert Into CompanyInfo Set CompanyId = ?,Email = ?`
	Utils.DB().Exec(template, companyId, info.ToEmail)
	cryPassword, _ := bcrypt.GenerateFromPassword([]byte(info.Account.Password), 10)
	template = `Insert Into Account Set Account = ?,PassWord = ? ,CompanyId = ?`
	Utils.DB().Exec(template, info.Account.Account, string(cryPassword), companyId)
	return true, nil
}

func CheckAccountUnique(account string) bool {
	template := `Select CompanyId From ShippingTraceability.Account Where Account = ? limit 1`
	rows, err := Utils.DB().Query(template, account)
	if err != nil {
		log.Println("[CheckAccountUnique]数据库异常", err)
		return false
	}
	defer rows.Close()
	return !rows.Next()
}

func CheckEmailUnique(email string) bool {
	template := `Select CompanyId From ShippingTraceability.CompanyInfo Where Email = ? limit 1`
	rows, err := Utils.DB().Query(template, email)
	if err != nil {
		log.Println("[CheckAccountUnique]数据库异常", err)
		return false
	}
	defer rows.Close()
	return !rows.Next()
}

func TryUpdateCompany(info Utils.CompanyBasicInfo) bool {
	template := `Select CompanyName, CompanyType From Company Where CompanyId = ? Limit 1`
	rows, err := Utils.DB().Query(template, info.CompanyId)
	if err != nil {
		log.Println("[TryUpdateCompany]数据库异常", err)
		return false
	}
	defer rows.Close()
	if !rows.Next() {
		return false
	}
	var oldInfo Utils.CompanyBasicInfo
	rows.Scan(&oldInfo.CompanyName, &oldInfo.CompanyType)
	oldInfo.CompanyId = info.CompanyId
	fmt.Println(info, oldInfo)
	if reflect.DeepEqual(oldInfo, info) {
		return false
	}
	template = `Update Company Set CompanyName = ?,CompanyType = ? WHere CompanyId = ?`
	result, err := Utils.DB().Exec(template, info.CompanyName, info.CompanyType, info.CompanyId)
	if err != nil {
		log.Println("[TryUpdateCompany]数据库异常", err)
		return false
	}
	Utils.RDB().Del(string(info.CompanyId))
	num, _ := result.RowsAffected()
	return num == 1
}

func TryUpdateCompanyInfo(info Utils.CompanyInfo) bool {
	template := `Select Phone, AddressId, Email From CompanyInfo Where CompanyId = ?`
	rows, err := Utils.DB().Query(template, info.CompanyId)
	if err != nil {
		log.Println("[TryUpdateCompany]数据库异常", err)
		return false
	}
	defer rows.Close()
	if !rows.Next() {
		return false
	}
	var phone, email string
	var addressId int64
	rows.Scan(&phone, &addressId, &email)
	fmt.Println(phone, addressId, email, info)
	Utils.RDB().Set(string(info.CompanyId)+"addressId", addressId, time.Minute)
	if info.Email == email && info.Phone == phone {
		return false
	}
	template = `Update CompanyInfo Set Phone = ?,Email = ? Where CompanyId = ?`
	result, err := Utils.DB().Exec(template, info.Phone, info.Email, info.CompanyId)
	if err != nil {
		log.Println("[TryUpdateCompany]数据库异常", err)
		return false
	}
	num, _ := result.RowsAffected()
	return num == 1
}

func TryUpdateAddress(info Utils.AddressInfo, id int64) bool {
	addressId, _ := Utils.RDB().Get(string(id) + "#addressId").Result()
	if addressId == "1" {
		template := `Insert Into Address Set Country=?,City=?,Address=?`
		result, err := Utils.DB().Exec(template, info.Country, info.City, info.Address)
		if err != nil {
			return false
		}
		aId, _ := result.LastInsertId()
		template = `Update CompanyInfo Set AddressId = ? Where CompanyId = ?`
		rows, err := Utils.DB().Exec(template, aId, id)
		num, _ := rows.RowsAffected()
		return num == 1
	}
	template := `Select Country, City, Address From Address Where AddressId = ?`
	rows, err := Utils.DB().Query(template, addressId)
	if err != nil {
		log.Println("[TryUpdateAddress]数据库异常", err)
		return false
	}
	defer rows.Close()
	if !rows.Next() {
		return false
	}
	var oldInfo Utils.AddressInfo
	rows.Scan(&oldInfo.Country, &oldInfo.City, &oldInfo.Address)
	fmt.Println(oldInfo, info)
	if reflect.DeepEqual(oldInfo, info) {
		fmt.Println(false)
		return false
	}
	template = `Update Address Set Country = ?,City = ?,Address = ? Where AddressId = ?`
	result, err := Utils.DB().Exec(template, info.Country, info.City, info.Address, addressId)
	if err != nil {
		return false
	}
	num, _ := result.RowsAffected()
	return num == 1
}

func CheckEmail(account interface{}, email string) bool {
	if _, ok := account.(string); ok {
		template := `Select CompanyId From Account Where Account = ?`
		rows, err := Utils.DB().Query(template, account)
		if err != nil {
			log.Println("[CheckEmail]出错了")
			return false
		}
		defer rows.Close()
		if !rows.Next() {
			return false
		}
		rows.Scan(&account)
	}
	template := `Select Email From CompanyInfo Where CompanyId = ?`
	rows, err := Utils.DB().Query(template, account)
	if err != nil {
		log.Println("[CheckEmail]出错了")
		return false
	}
	if !rows.Next() {
		return false
	}
	var oldEmail string
	rows.Scan(&oldEmail)
	return oldEmail == email
}

func ChangePassword(account interface{}, password string) bool {
	var template string
	if _, ok := account.(int64); ok {
		template = `Update Account Set PassWord = ? Where CompanyId = ?`
	} else {
		template = `Update Account Set PassWord = ? Where Account = ?`
	}
	cryPassword, err := bcrypt.GenerateFromPassword([]byte(password), 10)
	if err != nil {
		return false
	}
	result, err := Utils.DB().Exec(template, cryPassword, account)
	if err != nil {
		log.Println("[ChangePassword]出错了")
		return false
	}
	num, _ := result.RowsAffected()
	return num == 1
}

//GetCompanyBasicInfo 通过Id获取企业的类型以及名称 (name,type)
func GetCompanyBasicInfo(companyId int64) (string, string) {
	if companyId == 0 {
		return "", ""
	}
	aCompanyId := strconv.FormatInt(companyId, 10)
	companyName, err := Utils.RDB().Get(aCompanyId + "#Companyname").Result()
	companyType, err := Utils.RDB().Get(aCompanyId + "#Companytype").Result()
	if err != nil { //Redis中没有找到则进行查找
		basicInfo, _ := CompanyBasicInfo(companyId)
		Utils.RDB().Set(aCompanyId+"#Comanytype", basicInfo.CompanyType, time.Minute*5)
		Utils.RDB().Set(aCompanyId+"#Companyname", basicInfo.CompanyName, time.Minute*5)
		return basicInfo.CompanyName, basicInfo.CompanyType
	}
	return companyName, companyType
}
