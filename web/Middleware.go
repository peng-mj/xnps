package web

import (
	"github.com/gin-gonic/gin"
	"github.com/juju/ratelimit"
	"net/http"
	"strings"
	"time"
	"tunpx/pkg/cache"
	"tunpx/pkg/jwtTool"
	"tunpx/web/api"
	"tunpx/web/dto"
	"tunpx/web/service"
)

type MiddleBase struct {
	token *jwtTool.Token
	salt  *cache.Cache
	kid   *service.Base
}
type AuthUser struct {
	MiddleBase
	User     string
	UserId   int32
	AuthCode int32
	IsAdmin  bool
}

func NewMiddle(db *service.Base) *MiddleBase {
	middle := MiddleBase{token: jwtTool.NewToken(), salt: cache.New(100), kid: db}
	return &middle
}

func (m *MiddleBase) AuthMiddle(ctx *gin.Context) {
	tokenString := ctx.Request.Header.Get("Authorization")
	if tokenString == "" {
		ctx.JSON(http.StatusUnauthorized, dto.Response{Code: http.StatusUnauthorized, ErMsg: "token empty"})
		ctx.Abort()
		return
	}
	tokenString = strings.Replace(tokenString, "Bearer ", "", 1)
	uid, err := m.token.Parse(tokenString)
	if err != nil {
		ctx.JSON(http.StatusUnauthorized, dto.Response{Code: http.StatusUnauthorized, ErMsg: err.Error()})
		ctx.Abort()
		return
	}
	ctx.Set("uid", uid)
	ctx.Next()
}

func (m *MiddleBase) GetUser(ctx *gin.Context) {
	if uid, ok := ctx.Get("uid"); ok {
		if user, err := service.NewAuthUser(m.kid).GetUserByUid(uid.(int64)); err != nil {
			api.RepError(ctx, http.StatusForbidden)
			ctx.Abort()
			return
		} else {
			ctx.Set("user", dto.User{
				Uid:       user.Uid,
				Id:        user.Id,
				AuthLevel: user.Level,
				ExpireAt:  user.ExpirationAt,
				OtpKey:    user.OTAKeys,
				Valid:     user.ExpirationAt >= time.Now().Unix(),
			})
			ctx.Next()
		}
	}
}

// Login 登录信息
func (m *MiddleBase) Login(ctx *gin.Context) {
	var login dto.LoginReq
	var err error
	if err = ctx.ShouldBindJSON(&login); err == nil {
		user, code := service.NewAuthUser(m.kid).CheckPasswd(&login)
		if code != 200 {
			api.RepError(ctx, code.Int())
		}
		token := m.token.Generate(user.Uid, time.Hour*12)
		api.Response(ctx, token)
		ctx.Abort()
		return
	}
	api.RepError(ctx, dto.ErrParam)
}

// func (m *MiddleBase) RateLimitMiddle(ctx *gin.Context) {
//
//		ctx.Next()
//	}

func (m *MiddleBase) RateLimitMiddle(fillInterval time.Duration, cap, quantum int64) gin.HandlerFunc {
	bucket := ratelimit.NewBucketWithQuantum(fillInterval, cap, quantum)
	return func(ctx *gin.Context) {
		// Maybe it's better to limit the frequency by IP address
		// ctx.Request.RemoteAddr   -> limit ip
		if bucket.TakeAvailable(1) < 1 {
			api.RepError(ctx, dto.ErrRateLimit)
			ctx.Abort()
			return
		}
		ctx.Next()
	}
}
