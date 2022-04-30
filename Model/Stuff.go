package Model

import (
	"fmt"
	"log"
	"main/Utils"
)

func GetStuff(companyId int64) ([]Utils.Stuff, error) {
	template := `Select StaffId, StaffName, StaffJob From Staff Where CompanyId = ?`
	rows, err := Utils.DB().Query(template, companyId)
	if err != nil {
		log.Println("[GetStuff]", err)
		return nil, err
	}
	defer rows.Close()
	var stuffs []Utils.Stuff
	var stuff Utils.Stuff
	for rows.Next() {
		rows.Scan(&stuff.StuffId, &stuff.StuffName, &stuff.StuffJob)
		stuffs = append(stuffs, stuff)
	}
	return stuffs, nil
}

func InsertStuff(stuff Utils.Stuff, companyId int64) (int64, error) {
	template := `Insert Into Staff Set StaffName = ?,StaffJob = ?,CompanyId = ?`
	rows, err := Utils.DB().Exec(template, stuff.StuffName, stuff.StuffJob, companyId)
	if err != nil {
		return 0, err
	}
	line, _ := rows.LastInsertId()
	return line, nil
}

func CheckStuffCompany(stuffId, companyId int64) bool {
	template := `Select CompanyId From Staff Where StaffId = ?`
	rows, err := Utils.DB().Query(template, stuffId)
	if err != nil || !rows.Next() {
		return false
	}
	defer rows.Close()
	var sCompanyId int64
	rows.Scan(&sCompanyId)
	return sCompanyId == companyId
}

func DeleteStuff(stuffId int64) error {
	template := `Update Staff Set isDelete = t Where StaffId = ? limit 1`
	_, err := Utils.DB().Exec(template, stuffId)
	return err
}

func GetStuffInfo(stuffId int64) (Utils.StuffInfo, int64, error) {
	var info Utils.StuffInfo
	template := `Select JoinDate, Sex, Phone, Email, Fax, BirthDay, AddressId From StuffInfo Where StuffId = ?`
	rows, err := Utils.DB().Query(template, stuffId)
	if err != nil {
		log.Println("[GetStuffInfo] make a mistake ", err)
		return info, 0, err
	}
	defer rows.Close()
	if !rows.Next() {
		log.Println("[GetStuffInfo] Query data not exists")
		return info, 0, fmt.Errorf("not Exists")
	}
	var addressId int64
	rows.Scan(&info.JoinDate, &info.Sex, &info.Phone, &info.Email, &info.Fax, &info.BirthDay, &addressId)
	template = `Select Country, City, Address From Address Where AddressId = ?`
	rows, err = Utils.DB().Query(template, addressId)
	if err != nil {
		return info, 0, err
	}
	defer rows.Close()
	rows.Next()
	rows.Scan(&info.AddressInfo.Country, &info.AddressInfo.City, &info.AddressInfo.Address)
	return info, addressId, nil
}

func UpdateStuffInfo(info Utils.StuffInfo) bool {
	template := `Update StuffInfo Set Sex = ?,Phone = ?,Fax = ?,BirthDay = ?,Email = ? Where StuffId = ? limit 1`
	result, err := Utils.DB().Exec(template, info.Sex, info.Phone, info.Fax, info.BirthDay, info.Email, info.StuffId)
	if err != nil {
		log.Println("[UpdateStuffInfo] make a mistake ", err)
		return false
	}
	num, err := result.RowsAffected()
	return num == 1
}

func UpdateStuffAddressInfo(info Utils.AddressInfo, addressId int64, stuffId int64) bool {
	if addressId != 1 {
		template := `Update Address Set Country = ?,City = ?,Address = ? Where AddressId = ? limit 1`
		result, err := Utils.DB().Exec(template, info.Country, info.City, info.Address, addressId)
		if err != nil {
			log.Println("[UpdateStuffAddressInfo] make a mistake ", err)
			return false
		}
		num, _ := result.RowsAffected()
		return num == 1
	}
	template := `Insert Into Address Set Country = ? ,City = ? ,Address = ?`
	result, err := Utils.DB().Exec(template, info.Country, info.City, info.Address)
	if err != nil {
		log.Println("[UpdateAddressInfo] make a mistake ", err)
		return false
	}
	addressId, _ = result.LastInsertId()
	template = `Update StuffInfo Set AddressId = ? Where StuffId = ? limit 1`
	result, err = Utils.DB().Exec(template, addressId, stuffId)
	if err != nil {
		log.Println("[UpdateAddressInfo] make a mistake ", err)
		return false
	}
	num, _ := result.RowsAffected()
	return num == 1
}
