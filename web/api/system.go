package api

import (
	"github.com/gin-gonic/gin"
	"math/rand"
	"net/http"
	"time"
	"tunpx/pkg/config"
	"tunpx/pkg/crypt"
	"tunpx/pkg/models"
	"tunpx/web/dto"
	"tunpx/web/service"
)

/**********        USER          *********/

type System struct {
	kit *service.Base
}

func NewSystem(kit *service.Base) *System {
	return &System{kit: kit}
}

func (s *System) GetConfig(ctx *gin.Context) {
	ctx.JSON(http.StatusOK, "")

}
func (s *System) Update(ctx *gin.Context) {
	ctx.JSON(http.StatusOK, "")
}

// Init to write config to file,and remove temp config file
func (s *System) Init(ctx *gin.Context) {

	var conf dto.ConfigReq
	var err error
	if err = ctx.BindJSON(&conf); err != nil {
		RepErrorWithMsg(ctx, dto.ErrParam, err.Error())
		ctx.Abort()
		return
	}
	if err = conf.Validity(); err != nil {
		RepErrorWithMsg(ctx, dto.ErrParam, err.Error())
		ctx.Abort()
		return
	}
	var fConf config.Config
	fConf.BasePath = conf.BasePath
	fConf.WebPort = conf.WebPort
	fConf.InitTime = time.Now().Unix()
	fConf.BridgePort = conf.BridgePort
	fConf.Remark = conf.OrgName
	fConf.DbType = "sqlite"
	fConf.AppKeys = crypt.RandStr().AddAll().GenerateList(16, 20)
	fConf.UsagePorts = conf.UsagePorts
	err = config.New(fConf.BasePath + "/conf.toml").Update(fConf)
	if err != nil {
		RepErrorWithMsg(ctx, dto.ErrCreateConfigFile, err.Error())
		ctx.Abort()
		return
	}
	auth := models.AuthUser{
		Username:     conf.Username,
		Password:     conf.Password,
		Emile:        conf.Username,
		OTAKeys:      crypt.RandStr().Generate(40),
		Uid:          crypt.SnowID(rand.Int63n(1024)),
		Level:        0,
		LastLoginIp:  ctx.Request.RemoteAddr,
		ExpirationAt: 4102415999,
		MaxConn:      conf.MaxConn,
		Valid:        true,
	}
	err = service.NewAuthUser(s.kit).Create(&auth)
	if err != nil {
		RepErrorWithMsg(ctx, dto.ErrCreateUser, err.Error())
		ctx.Abort()
		return
	}
	group := models.Group{
		Valid:        true,
		Uid:          auth.Uid,
		Name:         "default",
		UsagePorts:   config.NewPorts(conf.UsagePorts).String(),
		GroupType:    "default",
		MaxClientNum: int32(conf.MaxConn),
		Remark:       "default group",
	}
	err = service.NewGroup(s.kit).Create(&group)
	if err != nil {
		RepErrorWithMsg(ctx, dto.ErrCreateGroup, err.Error())
		ctx.Abort()
		return
	}

	ctx.JSON(http.StatusOK, "OK")
	ctx.Next()

}

// StaticInit  to load system init html and other static files
func (s *System) StaticInit(ctx *gin.Context) {
	ctx.JSON(http.StatusOK, "")
}

func (s *System) StaticSuccess(ctx *gin.Context) {
	ctx.JSON(http.StatusOK, "")
}
