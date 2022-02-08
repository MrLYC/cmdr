package model

type Command struct {
	ID        int    `storm:"id,increment"`
	Name      string `storm:"index" json:"name"`
	Version   string `storm:"index" json:"version"`
	Activated bool   `storm:"index" json:"activated"`
	Location  string `json:"location"`
}
