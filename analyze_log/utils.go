package analyze_log

import (
	"encoding/json"
	"io"
	"log"
	"strconv"
)

func CloseAll(cr ...io.Closer) {
	for _, c := range cr {
		if c != nil {
			if err := c.Close(); err != nil {
				log.Println(err)
			}
		}
	}
}

type StringType string

func (s *StringType) UnmarshalJSON(data []byte) error {
	if data[0] == '"' {
		// 如果是字符串，直接解析
		return json.Unmarshal(data, (*string)(s))
	}

	var num int64
	if err := json.Unmarshal(data, &num); err != nil {
		return err
	}
	*s = StringType(strconv.FormatInt(num, 10))
	return nil
}
