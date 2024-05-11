package api

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"strconv"
	"tunpx/pkg/crypt"
	"tunpx/pkg/models"
	myUitls "tunpx/pkg/myUtils"
	"tunpx/web/dto"
	"tunpx/web/service"
)

type Group struct {
	kit *service.Base
}

func NewGroup(dr *service.Base) *Group {
	return &Group{kit: dr}
}

func (b *Group) GetAll(ctx *gin.Context) {
	req := dto.GroupGetReq{}
	var err error
	if err = ctx.BindJSON(req); err != nil {
		RepErrorWithMsg(ctx, dto.ErrParam, err.Error())
		return
	}

	db := b.kit.Orm(models.Group{})
	if len(req.Filters.Ids) > 0 {
		db.Where("id in ?", req.Filters.Ids)
	}
	if req.Filters.Uid != 0 {
		db.Where("uid = ?", req.Filters.Uid)
	}
	if req.Filters.Name != "" {
		db.Where("name = ?", req.Filters.Name)
	}
	res := make([]models.Group, 0)
	db.Find(&res)

	ctx.JSON(http.StatusOK, res)
}

func (b *Group) Create(ctx *gin.Context) {
	group := models.Group{}
	var err error
	if err = ctx.BindJSON(&group); err != nil {
		RepErrorWithMsg(ctx, dto.ErrParam, err.Error())
		return
	}
	group.Id = crypt.SnowID(1)
	user := GetUser(ctx)
	if user == nil {
		RepErrorWithMsg(ctx, dto.ErrInternalErr, err.Error())
		return
	}
	group.Uid = user.Uid
	group.Id = 0
	if err = myUitls.NewPorts(nil).Load(group.UsagePorts); err != nil {
		RepErrorWithMsg(ctx, dto.ErrParam, err.Error()+" example string:8080-9000,9901-9901")
		return
	}

	err = b.kit.Orm(models.Group{}).Create(&group).Error
	if err != nil {
		RepError(ctx, dto.ErrDbError)
		return
	}

	ctx.JSON(http.StatusOK, group)
}
func (b *Group) Delete(ctx *gin.Context) {
	id, ok := ctx.GetQuery("id")
	if !ok {
		RepErrorWithMsg(ctx, dto.ErrParam, "id should not empty")
		return
	}
	gid, err := strconv.Atoi(id)
	if err != nil {
		RepErrorWithMsg(ctx, dto.ErrParam, err.Error())
		return
	}
	user := GetUser(ctx)
	if user == nil {
		RepErrorWithMsg(ctx, dto.ErrInternalErr, err.Error())
		return
	}

	err = b.kit.Orm(models.Group{}).Where("id = ? and uid = ?", gid, user.Uid).Delete(models.Group{}).Error
	if err != nil {
		RepError(ctx, dto.ErrDbError)
		return
	}
	ctx.JSON(http.StatusOK, id)
}
func (b *Group) Update(ctx *gin.Context) {
	group := models.Group{}
	var err error
	if err = ctx.BindJSON(&group); err != nil {
		RepErrorWithMsg(ctx, dto.ErrParam, err.Error())
		return
	}
	if group.Id == 0 {
		RepErrorWithMsg(ctx, dto.ErrParam, "id should not be null")
		return
	}
	group.Id = crypt.SnowID(1)
	user := GetUser(ctx)
	if user == nil {
		RepErrorWithMsg(ctx, dto.ErrInternalErr, err.Error())
		return
	}
	group.Uid = user.Uid
	group.Id = 0
	if err = myUitls.NewPorts(nil).Load(group.UsagePorts); err != nil {
		RepErrorWithMsg(ctx, dto.ErrParam, err.Error()+" example string:8080-9000,9901-9901")
		return
	}

	err = b.kit.Orm(models.Group{}).Omit("id", "create_at", "uid").Where("id = ?", group.Id).Updates(&group).Error
	if err != nil {
		RepError(ctx, dto.ErrDbError)
		return
	}

	ctx.JSON(http.StatusOK, group)
}
