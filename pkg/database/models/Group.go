package models

type Group struct {
	Id           int64  `gorm:"column:id;type:int;auto_increment;not null;primaryKey;" json:"id"`
	Name         string `gorm:"column:name;type:varchar(40);not null;" json:"name"`
	UserId       int64  `gorm:"column:user_id;type:int;not null;" json:"user_id"`
	UsagePorts   string `gorm:"column:ports;type:text;not null;" json:"ports"`
	GroupType    int64  `gorm:"column:group_type;type:int;not null;" json:"group_type"`
	CreateAt     int64  `gorm:"column:create_time;type:int;not null;" json:"create_at"`
	ModifyAt     int64  `gorm:"column:create_time;type:int;not null;" json:"modify_at"`
	MaxClientNum int32  `gorm:"column:create_time;type:int;not null;default:0" json:"max_client_num"`
	CurClientNum int32  `gorm:"column:create_time;type:int;not null;default:0" json:"cur_client_num"`  //可能不需要写入数据库
	OnClientNum  int32  `gorm:"column:on_client_num;type:int;not null;default:0" json:"on_client_num"` //可能不需要写入数据库
	Valid        bool   `gorm:"column:valid;not null;default:true" json:"valid"`
	Remark       string `gorm:"column:remark;type:text;not null;" json:"remark"`
}

func (s *Group) TableName() string {

	return "client_group"
}
