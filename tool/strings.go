package tool

import "strings"

func TrimToNil(str string) *string {
	if len(str) == 0 {
		return nil
	}
	res := strings.TrimSpace(str)
	return &res
}
