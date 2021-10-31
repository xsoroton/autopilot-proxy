package util

import "os"

func GetEnv(key string, defaultValue string) (setting string) {
	setting = os.Getenv(key)
	if setting == "" {
		setting = defaultValue
	}
	return
}
