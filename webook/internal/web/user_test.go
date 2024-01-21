package web

import (
	"bytes"
	"context"
	"errors"
	"geektime-basic/webook/internal/domain"
	"geektime-basic/webook/internal/service"
	svcmocks "geektime-basic/webook/internal/service/mocks"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/mock/gomock"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestUserHandler_SignUp(t *testing.T) {
	testCases := []struct {
		name string

		mock func(ctrl *gomock.Controller) service.UserService

		reqBody  string
		wantCode int
		wantBody string
	}{
		{
			name: "注册成功",
			mock: func(ctrl *gomock.Controller) service.UserService {
				userSvc := svcmocks.NewMockUserService(ctrl)
				userSvc.EXPECT().SingUp(gomock.Any(), domain.User{
					Email:    "123@qq.com",
					Password: "hello#work123",
				}).Return(nil)
				return userSvc
			},
			reqBody: `
{
		"email": "123@qq.com",
		"password": "hello#work123",
		"confirmPassword": "hello#work123"
}`,
			wantCode: http.StatusOK,
			wantBody: "注册成功",
		},
		{
			name: "参数不对, bind失败",
			mock: func(ctrl *gomock.Controller) service.UserService {
				userSvc := svcmocks.NewMockUserService(ctrl)
				return userSvc
			},
			reqBody: `
{
		"email": "123@qq.com",
		"password": "hello#work123",
		"confirmPassword": "hello#work123"
}`,
			wantCode: http.StatusBadRequest,
		},
		{
			name: "邮箱格式不对",
			mock: func(ctrl *gomock.Controller) service.UserService {
				userSvc := svcmocks.NewMockUserService(ctrl)
				return userSvc
			},
			reqBody: `
{
		"email": "123@q",
		"password": "hello#work123",
		"confirmPassword": "hello#work123"
}`,
			wantCode: http.StatusOK,
			wantBody: "邮箱不正确",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			server := gin.Default()
			// 用不上codeSvc
			h := NewUserHandler(tc.mock(ctrl), nil)
			h.RegisterRoutes(server)

			req, err := http.NewRequest(http.MethodPost,
				"/users/signup",
				bytes.NewBuffer([]byte(tc.reqBody)))
			require.NoError(t, err)
			//数据是json格式
			req.Header.Set("Content-Type", "application/json")

			t.Log(req)

			resp := httptest.NewRecorder()
			t.Log(resp)

			// 这就是HTTP请求进去GIN框架的入口
			// 当你这么调用的时候，GIN就会这么处理请求
			// 响应写回到 resp 里
			server.ServeHTTP(resp, req)

			assert.Equal(t, tc.wantCode, resp.Code)
			assert.Equal(t, tc.wantBody, resp.Body.String())
		})
	}
}

func TestMock(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	usersvc := svcmocks.NewMockUserService(ctrl)
	usersvc.EXPECT().SingUp(gomock.Any(), gomock.Any()).
		Return(errors.New("mock error"))

	err := usersvc.SingUp(context.Background(), domain.User{
		Email: "123@qq.com",
		//Password: "123456",
	})
	t.Log(err)
}
