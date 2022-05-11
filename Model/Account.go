package Model

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"golang.org/x/crypto/bcrypt"
	"log"
	"main/Utils"
	"reflect"
	"strconv"
)

func Login(account Utils.Account) (Utils.Account, bool, error) {
	var info Utils.Account
	template := `Select PassWord,CompanyId From Account Where Account = ? limit 1`
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

func Info(companyId int64) (info Utils.Info, err error) {
	info.CompanyInfo, err = getCompanyT(companyId)
	if err != nil {
		log.Println("[Info]->[getCompanyT] make a mistake ", err)
		return
	}
	info.AddressInfo, err = getAddressT(info.CompanyInfo.AddressId)
	if err != nil {
		log.Println("[Info]->[getAddressT] make a mistake ", err)
		return
	}
	info.CompanyBasicInfo, err = getCompanyBasicInfoT(companyId)
	if err != nil {
		log.Println("[Info]->[getAddressT] make a mistake ", err)
		return
	}
	return info, nil
}

func getCompanyT(companyId int64) (Utils.CompanyInfo, error) {
	var info Utils.CompanyInfo
	result, err := Utils.RDB().HMGet("company_"+string(companyId), "phone", "addressId", "email").Result()
	if err != nil || result[0] == nil || result[1] == nil || result[2] == nil {
		template := `Select Phone,AddressId, Email From CompanyInfo Where CompanyId = ? limit 1`
		rows, err := Utils.DB().Query(template, companyId)
		if err != nil {
			log.Println("[Info]数据库发生异常", err)
			return info, err
		}
		defer rows.Close()
		if !rows.Next() {
			return info, nil
		}
		rows.Scan(&info.Phone, &info.AddressId, &info.Email)

		tmp, _ := json.Marshal(info)
		var data map[string]interface{}
		json.Unmarshal(tmp, &data)
		Utils.RDB().HMSet("company_"+string(companyId), data)

	} else {
		info.Phone = result[0].(string)
		info.AddressId, err = strconv.ParseInt(result[1].(string), 10, 64)
		info.Email = result[2].(string)
	}
	return info, nil
}

func getAddressT(addressId int64) (info Utils.AddressInfo, err error) {
	result, err := Utils.RDB().HMGet("address_"+string(addressId), "country", "city", "address").Result()
	if err != nil || result[0] == nil || result[1] == nil || result[2] == nil {
		template := `Select Country, City, Address From Address Where AddressId = ? Limit 1`
		var rows *sql.Rows
		rows, err = Utils.DB().Query(template, addressId)
		if err != nil {
			log.Println("[QueryAddress]数据库发生异常", err)
			return
		}
		defer rows.Close()
		rows.Next()
		rows.Scan(&info.Country, &info.City, &info.Address)

		var data map[string]interface{}
		tmp, _ := json.Marshal(info)
		json.Unmarshal(tmp, &data)
		Utils.RDB().HMSet("address_"+string(addressId), data)

	} else {
		info.Country = result[0].(string)
		info.City = result[1].(string)
		info.Address = result[2].(string)
	}
	return
}

func getCompanyBasicInfoT(companyId int64) (info Utils.CompanyBasicInfo, err error) {
	result, err := Utils.RDB().HMGet("company_"+string(companyId), "companyName", "companyType").Result()
	if err != nil || result[0] == nil || result[1] == nil {
		template := `Select CompanyName, CompanyType From Company Where CompanyId = ? Limit 1`
		var rows *sql.Rows
		rows, err = Utils.DB().Query(template, companyId)
		if err != nil {
			log.Println("[GetCompanyById]数据库发生异常", err)
			return
		}
		defer rows.Close()
		rows.Next()
		rows.Scan(&info.CompanyName, &info.CompanyType)

		var data map[string]interface{}
		tmp, _ := json.Marshal(info)
		json.Unmarshal(tmp, &data)
		Utils.RDB().HMSet("company_"+string(companyId), data)
	} else {
		info.CompanyName = result[0].(string)
		info.CompanyType = result[1].(string)
	}
	info.CompanyId = companyId
	return
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
	template := `Select CompanyId From Account Where Account = ? limit 1`
	rows, err := Utils.DB().Query(template, account)
	if err != nil {
		log.Println("[CheckAccountUnique]数据库异常", err)
		return false
	}
	defer rows.Close()
	return !rows.Next()
}

func CheckEmailUnique(email string) bool {
	template := `Select CompanyId From CompanyInfo Where Email = ? limit 1`
	rows, err := Utils.DB().Query(template, email)
	if err != nil {
		log.Println("[CheckAccountUnique]数据库异常", err)
		return false
	}
	defer rows.Close()
	return !rows.Next()
}

func TryUpdateCompany(info Utils.CompanyBasicInfo) bool {

	oldInfo, err := getCompanyBasicInfoT(info.CompanyId)
	if err != nil {
		return false
	}
	oldInfo.CompanyId = info.CompanyId
	if reflect.DeepEqual(oldInfo, info) {
		return false
	}
	template := `Update Company Set CompanyName = ?,CompanyType = ? WHere CompanyId = ?`
	result, err := Utils.DB().Exec(template, info.CompanyName, info.CompanyType, info.CompanyId)
	if err != nil {
		log.Println("[TryUpdateCompany]数据库异常", err)
		return false
	}
	Utils.RDB().HDel("company_"+string(info.CompanyId), "companyName", "companyType")
	num, _ := result.RowsAffected()
	return num == 1
}

func TryUpdateCompanyInfo(info Utils.CompanyInfo) bool {
	oldInfo, err := getCompanyT(info.CompanyId)
	if err != nil {
		return false
	}
	info.AddressId = oldInfo.AddressId
	oldInfo.CompanyId = info.CompanyId
	if reflect.DeepEqual(oldInfo, info) {
		return false
	}
	template := `Update CompanyInfo Set Phone = ?,Email = ? Where CompanyId = ?`
	result, err := Utils.DB().Exec(template, info.Phone, info.Email, info.CompanyId)
	if err != nil {
		log.Println("[TryUpdateCompany]数据库异常", err)
		return false
	}
	Utils.RDB().HDel("company_"+string(info.CompanyId), "phone", "email")
	num, _ := result.RowsAffected()
	return num == 1
}

func TryUpdateAddress(info Utils.AddressInfo, id int64) bool {
	emm, err := getCompanyT(id)
	if err != nil {
		return false
	}
	oldInfo, err := getAddressT(emm.AddressId)
	if err != nil {
		return false
	}
	if reflect.DeepEqual(oldInfo, info) {
		fmt.Println(false)
		return false
	}
	if emm.AddressId == 1 {
		template := `Insert Into Address Set Country=?,City=?,Address=?`
		result, err := Utils.DB().Exec(template, info.Country, info.City, info.Address)
		if err != nil {
			return false
		}
		aId, _ := result.LastInsertId()
		template = `Update CompanyInfo Set AddressId = ? Where CompanyId = ?`
		rows, err := Utils.DB().Exec(template, aId, id)
		num, _ := rows.RowsAffected()
		Utils.RDB().HSet("company_"+string(id), "addressId", aId)
		return num == 1
	}
	template := `Update Address Set Country = ?,City = ?,Address = ? Where AddressId = ?`
	result, err := Utils.DB().Exec(template, info.Country, info.City, info.Address, emm.AddressId)
	if err != nil {
		return false
	}
	Utils.RDB().HDel("address_"+string(id), "country", "city", "address")
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
	info, _ := getCompanyBasicInfoT(companyId)
	return info.CompanyName, info.CompanyType
}
