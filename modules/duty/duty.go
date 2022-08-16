package duty

import "altsub/models"

func Add(d *models.MDuty) error {
	//
	// TODO：N月份排班班次已出，请注意各自排班安排
	//
	return d.Add()
}

func Delete(d *models.MDuty) error {
	return d.HardDelete()
}

func Fetch(ds *models.MDuties) error {
	return ds.Fetch("Groups.Ats", "Groups.Users")
}

func AddGroup(d *models.MDuty, g *models.MDutyGroup) error {
	return d.AddGroup(g)
}

func DeleteGroup(d *models.MDuty, g *models.MDutyGroup) error {
	return d.DeleteGroup(g)
}
