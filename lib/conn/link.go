package conn

import "time"

type Secret struct {
	Password string
	Conn     *Conn
}

func NewSecret(passwd string, conn *Conn) *Secret {
	return &Secret{
		Password: passwd,
		Conn:     conn,
	}
}

type Link struct {
	ConnType   string //连接类型
	Host       string //目标
	Crypt      bool   //加密
	Compress   bool   //压缩
	LocalProxy bool   //本地代理
	RemoteAddr string //源头的IP信息
	Option     Options
}

type Option func(*Options)

type Options struct {
	Timeout time.Duration
}

var defaultTimeOut = time.Second * 5

func NewLink(connType string, host string, crypt bool, compress bool, remoteAddr string, localProxy bool, opts ...Option) *Link {
	options := newOptions(opts...)

	return &Link{
		RemoteAddr: remoteAddr,
		ConnType:   connType,
		Host:       host,
		Crypt:      crypt,
		Compress:   compress,
		LocalProxy: localProxy,
		Option:     options,
	}
}

func newOptions(opts ...Option) Options {
	opt := Options{
		Timeout: defaultTimeOut,
	}
	for _, o := range opts {
		o(&opt)
	}
	return opt
}

func LinkTimeout(t time.Duration) Option {
	return func(opt *Options) {
		opt.Timeout = t
	}
}
