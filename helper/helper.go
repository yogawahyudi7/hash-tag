package helper

import (
	"fmt"
	"strings"
)

func QueryIN(ids []int) (string, []any) {
	ar := []string{}
	er := []any{}
	for i, v := range ids {

		ar = append(ar, fmt.Sprintf("$%v", i+1))
		er = append(er, fmt.Sprintf("%v", v))
	}

	iQuery := strings.Join(ar, ",")

	return iQuery, er
}
