package logger

import "github.com/sirupsen/logrus"

func init() {
	// Text output formatter, can also be set to JsonFormatter
	textFormatter := &Formatter{
		TimestampFormat: "2006-01-02 15:04:05.000",
		LogFormat:       "%time% %file%:%line% %lvl%:%msg%",
	}
	logrus.SetFormatter(textFormatter)

	logrus.SetLevel(logrus.InfoLevel)

	// Print the caller location
	logrus.SetReportCaller(true)
}
