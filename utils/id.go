package utils

import "github.com/segmentio/ksuid"

func GenerateKsuid() string {
	id := ksuid.New()
	return id.String()
}
