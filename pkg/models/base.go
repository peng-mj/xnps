package models

func GetAllTable() []interface{} {
	return []interface{}{AuthUser{}, Group{}, Client{}, Tunnel{}, Blacklist{}, Config{}, Log{}}
}
