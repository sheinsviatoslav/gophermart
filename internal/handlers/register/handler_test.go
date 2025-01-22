package register

//
//import (
//	"bytes"
//	"encoding/json"
//	"errors"
//	"github.com/golang/mock/gomock"
//	"github.com/sheinsviatoslav/gophermart/internal/common"
//	"github.com/sheinsviatoslav/gophermart/internal/mocks"
//	"github.com/stretchr/testify/assert"
//	"github.com/stretchr/testify/require"
//	"io"
//	"net/http"
//	"net/http/httptest"
//	"testing"
//)
//
//func TestRegisterHandler(t *testing.T) {
//	ctrl := gomock.NewController(t)
//	defer ctrl.Finish()
//
//	s := mocks.NewMockStorage(ctrl)
//
//	type want struct {
//		code           int
//		response       string
//		contentType    string
//		loginExists    bool
//		loginExistsErr error
//	}
//
//	tests := []struct {
//		name string
//		body common.UserCredentials
//		want want
//	}{
//		{
//			name: "success register",
//			body: common.UserCredentials{
//				Login:    "my_login",
//				Password: "my_password",
//			},
//			want: want{
//				code:        200,
//				response:    "user successfully registered",
//				contentType: "text/plain",
//			},
//		},
//		{
//			name: "empty login",
//			body: common.UserCredentials{
//				Login:    "",
//				Password: "my_password",
//			},
//			want: want{
//				code:        400,
//				response:    "login is required\n",
//				contentType: "text/plain; charset=utf-8",
//			},
//		},
//		{
//			name: "short login",
//			body: common.UserCredentials{
//				Login:    "my",
//				Password: "my_password",
//			},
//			want: want{
//				code:        400,
//				response:    "login is too short\n",
//				contentType: "text/plain; charset=utf-8",
//			},
//		},
//		{
//			name: "short password",
//			body: common.UserCredentials{
//				Login:    "my_login",
//				Password: "my_pass",
//			},
//			want: want{
//				code:        400,
//				response:    "password is too short\n",
//				contentType: "text/plain; charset=utf-8",
//			},
//		},
//		{
//			name: "empty password",
//			body: common.UserCredentials{
//				Login:    "my_login",
//				Password: "",
//			},
//			want: want{
//				code:        400,
//				response:    "password is required\n",
//				contentType: "text/plain; charset=utf-8",
//			},
//		},
//		//{
//		//	name: "login already exists",
//		//	body: common.UserCredentials{
//		//		Login:    "my_login",
//		//		Password: "my_password",
//		//	},
//		//	want: want{
//		//		code:        409,
//		//		response:    "user with this login already exists",
//		//		contentType: "text/plain; charset=utf-8",
//		//		loginExists: true,
//		//	},
//		//},
//		//{
//		//	name: "login already exists error",
//		//	body: common.UserCredentials{
//		//		Login:    "my_login",
//		//		Password: "my_password",
//		//	},
//		//	want: want{
//		//		code:        400,
//		//		response:    "json: cannot unmarshal string into Go value of type common.UserCredentials\n",
//		//		contentType: "text/plain; charset=utf-8",
//		//	},
//		//},
//		{
//			name: "invalid json",
//			body: common.UserCredentials{
//				Login:    "my_login",
//				Password: "my_password",
//			},
//			want: want{
//				code:           400,
//				response:       "json: cannot unmarshal string into Go value of type common.UserCredentials\n",
//				contentType:    "text/plain; charset=utf-8",
//				loginExistsErr: errors.New("error with checking login"),
//			},
//		},
//	}
//
//	for _, test := range tests {
//		t.Run(test.name, func(t *testing.T) {
//			body, _ := json.Marshal(test.body)
//
//			if test.name == "invalid json" {
//				body, _ = json.Marshal(test.name)
//			}
//
//			r := httptest.NewRequest(http.MethodPost, "/api/user/register", bytes.NewReader(body))
//			w := httptest.NewRecorder()
//
//			s.EXPECT().CheckLoginExists(r.Context(), test.body.Login).Return(
//				test.want.loginExists,
//				test.want.loginExistsErr,
//			).AnyTimes()
//			s.EXPECT().CreateUser(r.Context(), common.UserCredentials{
//				Login:    test.body.Login,
//				Password: test.body.Password,
//			}).Return(test.want.loginExistsErr).AnyTimes()
//			NewHandler(s).Handle(w, r)
//
//			res := w.Result()
//			assert.Equal(t, test.want.code, res.StatusCode)
//			defer res.Body.Close()
//			resBody, err := io.ReadAll(res.Body)
//
//			require.NoError(t, err)
//
//			assert.Equal(t, test.want.response, string(resBody))
//			assert.Equal(t, test.want.contentType, res.Header.Get("Content-Type"))
//		})
//	}
//}
