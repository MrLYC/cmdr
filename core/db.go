package core

import (
	"fmt"
	"path"

	"github.com/mrlyc/cmdr/define"
	"github.com/mrlyc/cmdr/model"
	"github.com/mrlyc/cmdr/utils"
)

func GetClient() *model.Client {
	cfg := define.Configuration
	cmdrDir := cfg.GetString("cmdr.root")
	name := cfg.GetString("database.name")
	client, err := model.Open("sqlite3", fmt.Sprintf("file:%s?cache=shared&_fk=1", path.Join(cmdrDir, name)))
	utils.CheckError(err)
	return client
}
