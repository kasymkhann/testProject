package helpers

import "testProject/pkg/logging"

// LogAndReturnError записывает ошибку в журнал и возвращает ее в качестве результата.
// - logger: Экземпляр логгера, который будет использоваться для записи ошибки.
// - msg: Сообщение или описание ошибки.
// - err: Сама ошибка, которая будет залогирована.
func LogAndReturnError(logger *logging.Logger, message string, err error) error {
	logger.Errorf("%s: %v", message, err)
	return err
}
