package models

type Client struct {
	Id                 int64   `gorm:"column:id;type:integer;auto_increment;not null;primaryKey;" json:"id"`
	GroupId            int64   `gorm:"column:group_id;type:integer;not null;default:0" json:"group_id"`              //所属分组ID
	AccessId           string  `gorm:"column:access_id;type:varchar(64);not null" json:"access_id"`                  //认证ID
	AccessKey          string  `gorm:"column:access_key;type:varchar(128);not null" json:"access_key"`               //认证密钥
	RemoteAddr         string  `gorm:"column:remote_addr;type:varchar(20);not null;default: " json:"remote_addr"`    //客户端地址
	Name               string  `gorm:"column:name;type:varchar(64);not null;default: " json:"name"`                  //名称
	Remark             string  `gorm:"column:remark;type:text;not null;default: " json:"remark"`                     //备注
	Valid              bool    `gorm:"column:valid;not null;default:true" json:"valid"`                              //是否启用
	Connected          bool    `gorm:"column:connected;not null;default:true" json:"connected"`                      //是否已经连接
	Crypt              bool    `gorm:"column:crypt;not null;default:false" json:"crypt"`                             //是否加密
	Compress           bool    `gorm:"column:compress;not null;default:false:" json:"compress"`                      //是否压缩
	RateLimit          uint32  `gorm:"column:rate_limit;type:integer;default:0;not null" json:"rate_limit"`          //网速限制
	FlowExport         float64 `gorm:"column:flow_export;type:real;not null;default:0" json:"flow_export"`           //流出流量的 KB
	FlowInput          float64 `gorm:"column:flow_in;type:real;not null;default:0" json:"flow_input"`                //流如的流量 KB
	MaxConn            int32   `gorm:"column:max_conn;type:integer;not null;default:100" json:"max_conn"`            //最大连接数 ，这个应该放在隧道上边
	NowConn            int32   `gorm:"column:now_conn;type:integer;not null;default:0" json:"now_conn"`              //当前连接数，应该放在隧道上边
	AllowUseConfigFile bool    `gorm:"column:allow_file_config;not null;default:true" json:"allow_use_config_file"`  //不允许用户使用配置文件登录
	MaxTunnelNum       uint32  `gorm:"column:max_tunnel_num;type:integer;not null;default:20" json:"max_tunnel_num"` //最大限制20个，单个终端
	Version            string  `gorm:"column:version;type:text;not null;default:0" json:"version"`                   //客户端的软件版本
	ActiveAt           int64   `gorm:"column:active_time;type:integer;default:1672502400;not null" json:"active_at"` //活跃时间，上一次连接时间，断开连接时间
	BlackId            uint32  `gorm:"column:black_id;type:integer;default:0;not null" json:"blackId"`               //黑名单分组ID
}

func (*Client) TableName() string {
	return "client"
}
