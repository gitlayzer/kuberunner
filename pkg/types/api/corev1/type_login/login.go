package type_login

import (
	"errors"
	"github.com/gitlayzer/kuberunner/pkg/config"
)

var Login login

type login struct{}

func (l *login) Login(username, password string) (err error) {
	if username == config.GetUsername() && password == config.GetPassword() {
		return nil
	} else {
		return errors.New("username or password is wrong")
	}
}
