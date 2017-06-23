package util

import (
	"fmt"
	"os"
	"regexp"
)

// IsBackupID tests if a string is a valid backup ID.
// A valid backup ID is a 10 digit integer, representing
// a Unix timestamp.
func IsBackupID(id string) bool {
	re := regexp.MustCompile("\\d{10}")
	return re.Match([]byte(id))
}

// ClientCertAndKeyFromEnv loads the client cert and key filepaths
// from the respective environment variables RS_ETCD_CLIENT_CERT
// and RS_ETCD_CLIENT_KEY.
func ClientCertAndKeyFromEnv() (string, string, error) {
	certfile := os.Getenv("RS_ETCD_CLIENT_CERT")
	if certfile == "" {
		return "", "", fmt.Errorf("Can't find client cert file: RS_ETCD_CLIENT_CERT environment variable is not set.")
	}
	clientkey := os.Getenv("RS_ETCD_CLIENT_KEY")
	if clientkey == "" {
		return "", "", fmt.Errorf("Can't find client key file: RS_ETCD_CLIENT_KEY environment variable is not set.")
	}
	return certfile, clientkey, nil
}

// CACertFromEnv loads the CA cert filepath
// from the respective environment variable RS_ETCD_CA_CERT.
func CACertFromEnv() (string, error) {
	cacertfile := os.Getenv("RS_ETCD_CA_CERT")
	if cacertfile == "" {
		return "", fmt.Errorf("Can't find CA cert file: RS_ETCD_CA_CERT environment variable is not set.")
	}
	return cacertfile, nil
}
