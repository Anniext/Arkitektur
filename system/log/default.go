package log

import (
	"path"
)

func InitSystemLogger(tmp, mode string) (err error) {
	logPath := path.Join(tmp, "log/server.log")

	slogConfig := NewSlogOption(
		WithFilenameOption(logPath),
		WithMaxSizeOption(10),
		WithMaxBackupsOption(10),
		WithMaxAgeOption(10),
		WithCompressOption(true),
		WithProdLevelOption("info"),
		WithModeOption(mode),
	)

	_, err = NewSlogCore(slogConfig)
	if err != nil {
		return err
	}
	return nil
}
