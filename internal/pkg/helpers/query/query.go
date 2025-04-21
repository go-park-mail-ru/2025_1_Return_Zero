package query

import (
	"net/url"
	"strconv"
)

func ReadInt(qs url.Values, key string, defaultValue int) (int, error) {
	str := qs.Get(key)
	if str == "" {
		return defaultValue, nil
	}

	i, err := strconv.Atoi(str)
	if err != nil {
		return 0, err
	}

	return i, nil
}
