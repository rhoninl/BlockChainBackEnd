package Model

import (
	"log"
	"main/Utils"
)

func GetJointVenture(companyId int64) ([]Utils.CompanyList, error) {
	template := `Select CompanyId, TargetCompanyId From ShippingTraceability.Relation Where CompanyId = ? Or TargetCompanyId = ?`
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
	var fromId int64
	rows.Scan(&fromId)
	template = `Insert Into Relation Set CompanyId = ?,TargetCompanyId = ?`
	result, err := Utils.DB().Exec(template, fromId, reply.CompanyId)
	if err != nil {
		return false
	}
	num, err := result.RowsAffected()
	return num == 1
}

func CheckCompanyFriend(company1, company2 int64) bool {
	template := `Select CompanyId From Relation Where CompanyId = ? And TargetCompanyId = ? Or TargetCompanyId = ? And CompanyId = ?`
	rows, err := Utils.DB().Query(template, company1, company2, company2, company1)
	if err != nil {
		log.Println("[CheckCompanyFriend] make a mistake", err)
		return false
	}
	return !rows.Next()
}
