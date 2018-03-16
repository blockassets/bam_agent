package tool

import "strings"

func TrimToNil(str string) *string {
	if len(str) == 0 {
		return nil
	}
	res := strings.TrimSpace(str)
	if len(res) == 0 {
		return nil
	}
	return &res
}
