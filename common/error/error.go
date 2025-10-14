package error

import (
	"github.com/sirupsen/logrus"
)

func WrapError(err error) error {
	logrus.Errorf("error: %v", err)
	return err
}
