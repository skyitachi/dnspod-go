package dnspod

import (
	"fmt"
)

type User struct {
	RealName string `json:"real_name"`
	UserType string `json:"user_type"`
	Tel string `json:"telephone"`
	IM string `json:"im"`
	Nick string `json:"nick"`
	ID string `json:"id"`
	Email string `json:"email"`
	Status string `json:"status"`
	EmailVerified string `json:"email_verified"`
	TelVerified string `json:"telephone_verified"`
	WeixinBinded string `json:"weixin_binded"`
	AgentPending bool `json:"agent_pending"`
	Balance int `json:"balance"`
	SmsBalance int `json:"smsbalance"`
	UserGrade string `json:"user_grade"`
}

func (u User) String() string {
	bs, _ := json.Marshal(u)
	return string(bs)
}

type UserInfo struct {
	User User `json:"user"`
}

type userWrapper struct {
	Status Status `json:"status"`
	Info UserInfo `json:"info"`
}

func userAction(action string) string {
	if action == "" {
		return "User.Info"
	}
	return fmt.Sprintf("User.%s", action)
}

func (s *DomainsService) GetUserInfo() (User, *Response, error) {
	path := userAction("Detail")
	wrapper := userWrapper{}

	payload := newPayLoad(s.client.CommonParams)

	res, err := s.client.post(path, payload, &wrapper)
	if err != nil {
		return User{}, res, err
	}
	if wrapper.Status.Code != "1" {
		return User{}, res, fmt.Errorf("unexpected error: %s", wrapper.Status.Message)
	}
	return wrapper.Info.User, res, err
}
