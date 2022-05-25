package time

import (
	"time"
)

// UtcToLocal utc time to local time
func UtcToLocal(utcTime string, format string, zone string) (localTime string, err error) {
	location, _ := time.LoadLocation(zone)
	parseTime, err := time.ParseInLocation(format, utcTime, location)
	if err != nil {
		return
	}
	return parseTime.Local().Format("2006-01-02 15:04:05"), nil
}

// EstToLocal est time to local time
func EstToLocal(estTime string, format string, zone string) (localTime string, err error) {
	location, _ := time.LoadLocation(zone)
	parseTime, err := time.ParseInLocation(format, estTime, location)
	if err != nil {
		return
	}
	return parseTime.Local().Format("2006-01-02 15:04:05"), nil
}

// GetCurrentDate get current date
func GetCurrentDate(format string) string {
	return time.Now().Format(format)
}

// GetPlusDate get date by ...
func GetPlusDate(format string, duration int) string {
	t := time.Now()
	b := t.AddDate(0, 0, duration)
	return b.Format(format)
}

// GetPlusYear get after or before year
func GetPlusYear(duration int) string {
	t := time.Now()
	b := t.AddDate(duration, 0, 0)
	return b.Format("2006")
}

// Str2Time time format change
func Str2Time(timeStr string, oldFormat string, newFormat string) (string, error) {
	t, err := time.Parse(oldFormat, timeStr)
	if err != nil {
		return "", err
	}
	return t.Format(newFormat), nil
}
