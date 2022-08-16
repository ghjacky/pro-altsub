package duty

import "altsub/models"

func AddIssueHandling(ih *models.MIssueHandling) error {
	return ih.Add()
}

func DeleteIssueHandling(ih *models.MIssueHandling) error {
	return ih.HardDelete()
}

func UpdateIssueHandling(ih *models.MIssueHandling) error {
	return ih.Update()
}

func GetIssueHandling(ih *models.MIssueHandling) error {
	return ih.Get()
}

func FetchIssueHandlings(ihs *models.MIssueHandlings) error {
	return ihs.Fetch()
}
