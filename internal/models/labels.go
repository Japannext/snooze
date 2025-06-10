package models

import "fmt"

func Identity(key string) string {
	return fmt.Sprintf("identity__%s", key)
}

func Internal(key string) string {
	return fmt.Sprintf("internal__%s", key)
}
