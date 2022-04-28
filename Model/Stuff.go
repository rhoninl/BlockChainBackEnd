package Model

import (
	"log"
	"main/Utils"
)

func GetStuff(companyId int64) ([]Utils.Stuff, error) {
	template := `Select StaffId, StaffName, StaffJob From ShippingTraceability.Staff Where CompanyId = ?`
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
	template := `Insert Into ShippingTraceability.Staff Set StaffName = ?,StaffJob = ?,CompanyId = ?`
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
		log.Println("[CheckStuffCompany] make a mistake", err)
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
