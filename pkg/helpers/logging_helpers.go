package helpers

import "testProject/pkg/logging"

func LogAndReturnError(logger *logging.Logger, message string, err error) error {
	logger.Errorf("%s: %v", message, err)
	return err
}
