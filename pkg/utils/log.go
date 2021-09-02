package utils

import log "github.com/sirupsen/logrus"

func WrapErrLog(err error, l *log.Entry, message, errMessage string) error {
	if err != nil {
		l.WithError(err).Error(errMessage)
	} else {
		l.Debug(message)
	}

	return err
}

//TODO find a better name
func PrintError(err error, l *log.Entry, errMessage string) error {
	if err != nil {
		l.WithError(err).Error(errMessage)
	}

	return err
}

func WrapWarningLog(err error, l *log.Entry, message, errMessage string) error {
	if err != nil {
		l.WithError(err).Warning(errMessage)
	} else {
		l.Debug(message)
	}

	return err
}
