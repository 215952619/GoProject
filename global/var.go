package global

import (
	"github.com/sirupsen/logrus"
	"time"
)

var (
	WebServeAddr          = ":9507"
	AppSecret             = "gowebproject"
	AppIssuer             = "18435186204@163.com"
	Mode                  string
	Test                  bool
	Level                 string
	Action                string
	Logger                *logrus.Logger
	LogPath               string = ".\\log\\project.log"
	DbPath                string = ".\\data.db"
	DefaultValidityPeriod        = time.Minute * 30
	TokenHeaderKey               = "Authorization"
	TokenQueryKey                = "_t"
	AuthedKey                    = "AuthedUser"
	DefaultEventTimeout          = time.Second * 10
	AesSecret                    = "defaultAesCryptoSecret"
	DefaultAesSecret             = []byte(AesSecret)[:16]
)
