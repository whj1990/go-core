package handler

import (
	"github.com/pkg/errors"
	"go.uber.org/zap"
)

func HandleError(err error) error {
	if err != nil {
		zap.L().Error(err.Error())
	}
	return err
}

func HandleNewError(message string) error {
	err := errors.New(message)
	zap.L().Error(err.Error())
	return err
}
