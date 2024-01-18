package tencent

import (
	"context"
	"fmt"
	"github.com/go-playground/assert/v2"
	"github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/common"
	"github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/common/profile"
	sms "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/sms/v20210111"
	"os"
	"testing"
)

func TestSender(t *testing.T) {
	secretId, ok := os.LookupEnv("SMS_SECRET_ID")
	if !ok {
		t.Fatal()
	}
	fmt.Println("secretId:", secretId)

	secretKey, ok := os.LookupEnv("SMS_SECRET_KEY")
	fmt.Println("secretKey:", secretKey)

	c, err := sms.NewClient(common.NewCredential(secretId, secretKey),
		"ap-nanjing", profile.NewClientProfile())
	if err != nil {
		t.Fatal(err)
	}

	s := NewService(c, "1400868098", "深圳麦好多科技有限公司")

	testCases := []struct {
		name    string
		tplId   string
		params  []string
		numbers []string
		wantErr error
	}{
		{
			name:    "发送验证码",
			tplId:   "1993013",
			params:  []string{"258396"},
			numbers: []string{"18682325051"},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			er := s.Send(context.Background(), tc.tplId, tc.params, tc.numbers...)
			assert.Equal(t, tc.wantErr, er)
		})
	}
}
