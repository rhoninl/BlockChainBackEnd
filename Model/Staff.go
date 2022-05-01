package Model

import (
	"fmt"
	"log"
	"main/Utils"
	"time"
)

func GetStaff(companyId int64) ([]Utils.Staff, error) {
	template := `Select StaffId, StaffName, StaffJob From Staff Where CompanyId = ?`
	rows, err := Utils.DB().Query(template, companyId)
	if err != nil {
		log.Println("[GetStaff]", err)
		return nil, err
	}
	defer rows.Close()
	var staffs []Utils.Staff
	var staff Utils.Staff
	for rows.Next() {
		rows.Scan(&staff.StaffId, &staff.StaffName, &staff.StaffJob)
		staffs = append(staffs, staff)
	}
	return staffs, nil
}

func InsertStaff(staff Utils.Staff, companyId int64) (int64, error) {
	template := `Insert Into Staff Set StaffName = ?,StaffJob = ?,CompanyId = ?`
	rows, err := Utils.DB().Exec(template, staff.StaffName, staff.StaffJob, companyId)
	if err != nil {
		return 0, err
	}
	line, _ := rows.LastInsertId()
	template = `Insert Into StaffInfo Set StaffId = ?,JoinDate = ?`
	Utils.DB().Exec(template, line, time.Now().Unix())
	return line, nil
}

func CheckStaffCompany(staffId, companyId int64) bool {
	template := `Select CompanyId From Staff Where StaffId = ?`
	rows, err := Utils.DB().Query(template, staffId)
	if err != nil || !rows.Next() {
		return false
	}
	defer rows.Close()
	var sCompanyId int64
	rows.Scan(&sCompanyId)
	return sCompanyId == companyId
}

func DeleteStaff(staffId int64) error {
	template := `Update Staff Set isDelete = 1 Where StaffId = ? limit 1`
	_, err := Utils.DB().Exec(template, staffId)
	return err
}

func GetStaffInfo(staffId int64) (Utils.StaffInfo, int64, error) {
	var info Utils.StaffInfo
	template := `Select JoinDate, Sex, Phone, Email, Fax, BirthDay, AddressId From StaffInfo Where StaffId = ?`
	rows, err := Utils.DB().Query(template, staffId)
	if err != nil {
		log.Println("[GetStaffInfo] make a mistake ", err)
		return info, 0, err
	}
	defer rows.Close()
	if !rows.Next() {
		log.Println("[GetStaffInfo] Query data not exists")
		return info, 0, fmt.Errorf("not Exists")
	}
	var addressId, birthday, joinDate int64
	rows.Scan(&joinDate, &info.Sex, &info.Phone, &info.Email, &info.Fax, &birthday, &addressId)
	info.JoinDate = time.Unix(joinDate, 0).Format("2006-01-02")
	info.BirthDay = time.Unix(birthday, 0).Format("2006-01-02")
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

func UpdateStaffInfo(info Utils.StaffInfo) bool {
	birthday, err := time.Parse("2006-01-02", info.BirthDay)
	if err != nil {
		log.Println("timestamp transform default")
		return false
	}
	template := `Update StaffInfo Set Sex = ?,Phone = ?,Fax = ?,BirthDay = ?,Email = ? Where StaffId = ? limit 1`
	result, err := Utils.DB().Exec(template, info.Sex, info.Phone, info.Fax, birthday.Unix(), info.Email, info.StaffId)
	if err != nil {
		log.Println("[UpdateStaffInfo] make a mistake ", err)
		return false
	}
	num, err := result.RowsAffected()
	return num == 1
}

func UpdateStaffAddressInfo(info Utils.AddressInfo, addressId int64, staffId int64) bool {
	if addressId != 1 {
		template := `Update Address Set Country = ?,City = ?,Address = ? Where AddressId = ? limit 1`
		result, err := Utils.DB().Exec(template, info.Country, info.City, info.Address, addressId)
		if err != nil {
			log.Println("[UpdateStaffAddressInfo] make a mistake ", err)
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
	template = `Update StaffInfo Set AddressId = ? Where StaffId = ? limit 1`
	result, err = Utils.DB().Exec(template, addressId, staffId)
	if err != nil {
		log.Println("[UpdateAddressInfo] make a mistake ", err)
		return false
	}
	num, _ := result.RowsAffected()
	return num == 1
}
