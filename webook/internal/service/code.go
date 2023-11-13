package service

import (
	"context"
	"fmt"
	"geektime-basic/webook/internal/repository"
	"geektime-basic/webook/internal/service/sms"
	"math/rand"
)

const codeTplId = "1877556"

var (
	ErrCodeVerifyTooManyTimes = repository.ErrCodeVerifyTooManyTimes
	ErrCodeSendTooMany        = repository.ErrCodeSendTooMany
)

type CodeService struct {
	repo   *repository.CodeRepository
	smsSvc sms.Service
}

func (svc *CodeService) Send(ctx context.Context,
	biz string,
	phone string) error {
	// 生成验证码
	code := svc.generateCode()
	fmt.Println("code", code)

	// 放进去redis
	err := svc.repo.Store(ctx, biz, phone, code)
	if err != nil {
		// 有问题
		return err
	}

	// 前面成功了，再发送出去
	err = svc.smsSvc.Send(ctx, codeTplId, []string{code}, phone)
	//if err != nil {
	// 这个地方怎么办？
	// 这意味着，Redis 有这个验证码，但是不好意思，
	// 我能不能删掉这个验证码？
	// 你这个 err 可能是超时的 err，你都不知道，发出了没
	// 在这里重试
	// 要重试的话，初始化的时候，传入一个自己就会重试的 smsSvc
	//}
	return err
}

func (svc *CodeService) generateCode() string {
	// 六位数， num 在 0， 999999 之间，包含 0 和 999999
	num := rand.Intn(100000)
	return fmt.Sprintf("%6d", num)
}
