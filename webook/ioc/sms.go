package ioc

import (
	"geektime-basic/webook/internal/service/sms"
	"geektime-basic/webook/internal/service/sms/memory"
)

func InitSMSService() sms.Service {
	return memory.NewService()
}
