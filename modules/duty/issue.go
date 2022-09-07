package duty

import "altsub/models"

func AddIssueHandling(ih *models.MIssueHandling) error {
	return ih.Add()
}

func DeleteIssueHandling(ih *models.MIssueHandling) error {
	return ih.HardDelete()
}

func CloseIssueHandling(ih *models.MIssueHandling) error {
	return ih.SoftDelete()
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

func FetchIssueHandlingEvents(ih *models.MIssueHandling) error {
	return ih.FetchEvents()
}