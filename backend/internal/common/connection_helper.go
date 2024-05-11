package common

import (
	"encoding/base64"
	"fmt"
	"strconv"
)

func EncodeCursor(id int) string {
	idStr := fmt.Sprintf("%d", id)
	return base64.StdEncoding.EncodeToString([]byte(idStr))
}

func DecodeCursor(after *string) (*int, error) {
	skip := 0
	if after != nil && *after != "" {
		cursorData, err := base64.StdEncoding.DecodeString(*after)
		if err != nil {
			return nil, err
		}
		skip, err = strconv.Atoi(string(cursorData))
		if err != nil {
			return nil, err
		}
	}
	return &skip, nil
}
