package service

import (
	"xiaohuazhu/internal/service/account"
	"xiaohuazhu/internal/service/record"
)

type Service struct {
	Account *account.Service
	Record  *record.Service
}

func New() *Service {
	return &Service{
		Account: account.NewService(),
		Record: record.NewService(),
	}
}
