package util

import (
	"regexp"
)

// IsBackupID tests if a string is a valid backup ID.
// A valid backup ID is a 10 digit integer, representing
// a Unix timestamp.
func IsBackupID(id string) bool {
	re := regexp.MustCompile("\\d{10}")
	return re.Match([]byte(id))
}
