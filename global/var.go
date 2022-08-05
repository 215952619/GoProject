package global

import (
	"fmt"
	"github.com/sirupsen/logrus"
	"sync"
	"time"
)

var (
	WebServeAddr            = ":9507"
	AppJwtSecret            = "gowebproject"
	AppIssuer               = "18435186204@163.com"
	Mode                    string
	Test                    bool
	Level                   string
	Action                  string
	Logger                  *logrus.Logger
	LogPath                 = ".\\log\\project.log"
	DbPath                  = ".\\data.db"
	DefaultValidityPeriod   = time.Minute * 30
	TokenHeaderKey          = "Authorization"
	TokenQueryKey           = "_t"
	AuthedKey               = "AuthedUser"
	DefaultEventTimeout     = time.Second * 10
	AesSecret               = "defaultAesCryptoSecret"
	AesLength               = 16
	DefaultCacheExpiredTime = time.Minute * 30
	GiteeClientId           = "be838d5776c3dfe6081c0fb75b5923cf464e877e9c12587098693a03cc3436ce"
	GiteeClientSecret       = "a731ffd81de778d95611a7cd80cd3cd18b2762e3e836a549cc935d739522b5b4"
	GithubClientId          = "c0dd74661681f232b222"
	GithubClientSecret      = "84394c77b869f7e52d11f69a8530ad8e8f977dfd"
	SsoRedirectUrl          = "http://10.0.7.112:3000/admin/login"
	GiteePersonalToken      = "f67ed18693a0e8fa0187e00f7fe19249"
	MaxTopArticle           = 5
	Local                   = "/"
	ChanPool                = &sync.Pool{
		New: func() interface{} {
			return make(chan bool)
		},
	}
	InitChan = ChanPool.Get().(chan bool)
)

func init() {
	go func() {
		select {
		case <-InitChan:
			if Mode == "dev" {
				Local = fmt.Sprintf("http://10.0.7.112%s/", WebServeAddr)
			}
		}
		defer func() {
			InitChan <- true
		}()
	}()
}
