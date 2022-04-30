package Model

import (
	"log"
	"main/Utils"
)

func GetJointVenture(companyId int64) ([]Utils.CompanyList, error) {
	template := `Select CompanyId, TargetCompanyId From ( Select * From Relation Where isDelete = 0) as A Where CompanyId = ? Or TargetCompanyId = ?`
	rows, err := Utils.DB().Query(template, companyId, companyId)
	if err != nil {
		log.Println("[GetJointVenture]数据库异常", err)
		return nil, err
	}
	defer rows.Close()
	var companyList []Utils.CompanyList
	var company Utils.CompanyList
	var aCompanyId, bCompanyId int64
	for rows.Next() {
		rows.Scan(&aCompanyId, &bCompanyId)
		if companyId == aCompanyId {
			aCompanyId = bCompanyId
		}
		company.CompanyId = aCompanyId
		company.CompanyName, company.CompanyType = GetCompanyBasicInfo(aCompanyId)
		companyList = append(companyList, company)
	}
	return companyList, nil
}

func PassReply(reply Utils.ReplyFriend) bool {
	template := `Select FromId From MessageQueue Where MessageId = ? limit 1`
	rows, err := Utils.DB().Query(template, reply.MessageId)
	if err != nil || !rows.Next() {
		return false
	}
	defer rows.Close()
	var fromId int64
	rows.Scan(&fromId)
	aCompanyId, bCompanyId := fromId, reply.CompanyId
	if aCompanyId > bCompanyId {
		aCompanyId, bCompanyId = bCompanyId, aCompanyId
	}
	template = `Select CompanyId From Relation WHere CompanyId = ? And TargetCompanyId = ? And isDelete = 1 Limit 1`
	rows, err = Utils.DB().Query(template, fromId, reply.CompanyId)
	if err != nil {
		log.Println("[PassReply] make a mistake", err)
		return false
	}
	defer rows.Close()
	if !rows.Next() {
		template = `Insert Into Relation Set CompanyId = ?,TargetCompanyId = ?`
	} else {
		template = `Update Relation Set isDelete = 1 Where CompanyId = ? And TargetCompanyId = ?`
	}
	result, err := Utils.DB().Exec(template, aCompanyId, bCompanyId)
	if err != nil {
		return false
	}
	num, err := result.RowsAffected()
	return num == 1

}

func CheckCompanyFriend(company1, company2 int64) bool {
	if company1 > company2 {
		company1, company2 = company2, company1
	}
	template := `Select CompanyId From (Select * From Relation Where isDelete = 0) as A Where CompanyId = ? And TargetCompanyId = ?`
	rows, err := Utils.DB().Query(template, company1, company2)
	if err != nil {
		log.Println("[CheckCompanyFriend] make a mistake", err)
		return false
	}
	defer rows.Close()
	return rows.Next()
}

func DeleteCompanyFriend(company1, company2 int64) error {
	if company1 > company2 {
		company1, company2 = company2, company1
	}
	template := `Update Relation Set isDelete = 1 Where CompanyId = ? And TargetCompanyId = ? limit 1`
	_, err := Utils.DB().Exec(template, company1, company2)
	return err
}
