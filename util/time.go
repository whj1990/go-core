package util

import "time"

func ParseDateTime(timeStr string) (time.Time, error) {
	return time.ParseInLocation("2006-01-02 15:04:05", timeStr, time.Local)
}

func ParseDateTimeMMddHHmmss(timeStr string) (time.Time, error) {
	return time.ParseInLocation("20060102150405", timeStr, time.Local)
}

func ParseDateTimeMMddHHmmss2(timeStr string) (time.Time, error) {
	return time.ParseInLocation("2006010215:04:05", timeStr, time.Local)
}

func ParseDate(timeStr string) (time.Time, error) {
	return time.ParseInLocation("20060102", timeStr, time.Local)
}

func ParseDateSlash(timeStr string) (time.Time, error) {
	return time.ParseInLocation("2006/01/02", timeStr, time.Local)
}

func ParseDateBar(timeStr string) (time.Time, error) {
	return time.ParseInLocation("2006-01-02", timeStr, time.Local)
}

func FormatDateTime(time time.Time) string {
	return time.Format("2006-01-02 15:04:05")
}

func FormatTimeyyyyMMddHHmmss(time time.Time) string {
	return time.Format("20060102150405")
}

func FormatTimeyyyyMMdd(time time.Time) string {
	return time.Format("20060102")
}

func FormatTimeyyyyMMddSlash(time time.Time) string {
	return time.Format("2006/01/02")
}

func FormatTimeyyyyMMddBar(time time.Time) string {
	return time.Format("2006-01-02")
}
