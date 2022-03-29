package auth

import (
	"github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"
	"testing"
	"xiaohuazhu/internal/model"

	_ "xiaohuazhu/internal/config"
)

func TestGenerateToken(t *testing.T) {
	mo := &model.AccountDTO{
		Id:             12,
		Username:       "test",
		Password:       "123123",
		ProfilePicture: "image/test.jpg",
	}
	token, err := GenerateToken(mo)
	assert.Nil(t, err, err.Error())
	assert.NotEmpty(t, token)

	parseToken, err := ParseToken(token)
	assert.Nil(t, err, err.Error())
	assert.NotEmpty(t, parseToken)

	logrus.Debugf("[%+v]", parseToken)
}
