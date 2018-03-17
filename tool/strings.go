package tool

import "strings"

func TrimToNil(str string) *string {
	res := strings.TrimSpace(str)
	if len(res) == 0 {
		return nil
	}
	return &res
}
