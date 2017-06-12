package gvar

type RuntimeConfig struct {
	Comm       BaseConf `json:"Comm"`
	ModuleList []Module `json:"moduleList"`
}

type BaseConf struct {
	Interval int    `json:"interval"`
	Listen   string `json:"listen"`
	Scheme   string `json:"scheme"`
	Username string `json:"username"`
	Password string `json:"password"`
	Uri      string `json:"uri"`
	Server   string `json:"server"`
	Server1  string `json:"server1"`
	Server2  string `json:"server2"`
}
type Module struct {
	Currverion    string `json:"currverion"`
	Enable        bool   `json:"enable"`
	Filetype      string `json:"filetype"`
	RunLevel      int    `json:"runlevel"`
	Md5sum        string `json:"md5sum"`
	Name          string `json:"name"`
	Updatenable   bool   `json:"updatenable"`
	Updateversion string `json:"updateversion"`
	URL           string `json:"url"`
}
