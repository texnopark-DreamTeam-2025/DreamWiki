package repository

import (
	"strconv"
	"strings"
	"time"
)

func decodeCursor(cursor *string) (timeFrom time.Time, idFrom int64) {
	minTime := time.Unix(0, 0)

	if cursor == nil {
		return minTime, 0
	}

	splittedCursor := strings.Split(*cursor, "\n")
	if len(splittedCursor) != 2 {
		return minTime, 0
	}

	timeFrom, err := time.Parse(time.RFC3339, splittedCursor[0])
	if err != nil {
		return minTime, 0
	}

	idFrom, err = strconv.ParseInt(splittedCursor[1], 10, 64)
	if err != nil {
		return minTime, 0
	}

	return timeFrom, idFrom
}

func encodeCursor(timeFrom time.Time, idFrom int64) string {
	return timeFrom.Format(time.RFC3339) + "\n" + strconv.FormatInt(idFrom, 10)
}
