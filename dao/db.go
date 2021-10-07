package dao

import (
	"github.com/mrlyc/cmdr/model"
	"github.com/mrlyc/cmdr/utils"
)

func GetClient() *model.Client {
	client, err := model.Open("sqlite3", "file:cmdr.db?cache=shared&_fk=1")
	utils.CheckError(err)
	return client
}
