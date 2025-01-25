package login

//
//import (
//	"bytes"
//	"encoding/json"
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
//func TestLoginHandler(t *testing.T) {
//	ctrl := gomock.NewController(t)
//	defer ctrl.Finish()
//
//	s := mocks.NewMockStorage(ctrl)
//
//	type want struct {
//		code        int
//		response    string
//		contentType string
//	}
//
//	tests := []struct {
//		name string
//		body common.UserCredentials
//		want want
//	}{
//		{
//			name: "success log in",
//			body: common.UserCredentials{
//				Login:    "my_login",
//				Password: "my_password",
//			},
//			want: want{
//				code:        200,
//				response:    "successful authorization",
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
//		{
//			name: "wrong password",
//			body: common.UserCredentials{
//				Login:    "my_login",
//				Password: "my_wrong_password",
//			},
//			want: want{
//				code:        401,
//				response:    "wrong password\n",
//				contentType: "text/plain; charset=utf-8",
//			},
//		},
//		{
//			name: "invalid json",
//			body: common.UserCredentials{
//				Login:    "my_login",
//				Password: "my_password",
//			},
//			want: want{
//				code:        400,
//				response:    "json: cannot unmarshal string into Go value of type common.UserCredentials\n",
//				contentType: "text/plain; charset=utf-8",
//			},
//		},
//	}
//
//	testHashedPassword := "$2a$10$L5kuKZBPemzy9/aUCCfKlulon./C1Hzl4MdrI1stNW.3hDpf0cOq."
//
//	for _, test := range tests {
//		t.Run(test.name, func(t *testing.T) {
//			body, _ := json.Marshal(test.body)
//
//			if test.name == "invalid json" {
//				body, _ = json.Marshal(test.name)
//			}
//
//			r := httptest.NewRequest(http.MethodPost, "/api/user/login", bytes.NewReader(body))
//			w := httptest.NewRecorder()
//
//			s.EXPECT().GetUserPasswordByLogin(r.Context(), test.body.Login).Return(testHashedPassword, nil).AnyTimes()
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
