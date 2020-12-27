package dateInput

import "time"

func ToTime(input string) (*time.Time, error) {
	dateTime, err := time.Parse("2006-01-02", input)
	return &dateTime, err
}
