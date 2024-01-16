package models

type Group struct {
	Id         int64 `gorm:"column:id;type:int;auto_increment;not null;primaryKey;" json:"Id"`
	Name       int64 `gorm:"column:name;type:int;auto_increment;not null;primaryKey;" json:"Name"`
	Ports      int64 `gorm:"column:ports;type:int;auto_increment;not null;primaryKey;" json:"Ports"`
	GroupType  int64 `gorm:"column:group_type;type:int;auto_increment;not null;primaryKey;" json:"GroupType"`
	CreateTime int64 `gorm:"column:create_time;type:int;auto_increment;not null;primaryKey;" json:"CreateTime"`
	Remark     int64 `gorm:"column:remark;type:int;auto_increment;not null;primaryKey;" json:"Remark"`
}

func (s *Group) TableName() string {
	return "group"
}
