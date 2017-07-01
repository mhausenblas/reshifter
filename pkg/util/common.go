package util

import (
	"errors"
	"fmt"
	"net"
	"os"
	"regexp"

	"github.com/mhausenblas/reshifter/pkg/types"
)

// LookupDistro returns a textual description for a Kube distro by type.
func LookupDistro(distrotype types.KubernetesDistro) string {
	var distro string
	switch distrotype {
	case types.Vanilla:
		distro = "Vanilla Kubernetes"
	case types.OpenShift:
		distro = "OpenShift"
	default:
		distro = "not a Kubernetes distro"
	}
	return distro
}

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

// S3CredFromEnv loads S3 access key and secret
// from the respective environment variable ACCESS_KEY_ID
// and SECRET_ACCESS_KEY.
func S3CredFromEnv() (string, string, error) {
	accesskeyID := os.Getenv("ACCESS_KEY_ID")
	if accesskeyID == "" {
		return "", "", fmt.Errorf("Can't find S3 access key: ACCESS_KEY_ID environment variable is not set.")
	}
	secretaccessKey := os.Getenv("SECRET_ACCESS_KEY")
	if secretaccessKey == "" {
		return "", "", fmt.Errorf("Can't find S3 access key: SECRET_ACCESS_KEY environment variable is not set.")
	}
	return accesskeyID, secretaccessKey, nil
}

// ExternalIP retrieves the public IP of the host
// ReShifter is running on, adapted from:
// https://stackoverflow.com/questions/23558425/how-do-i-get-the-local-ip-address-in-go
func ExternalIP() (string, error) {
	ifaces, err := net.Interfaces()
	if err != nil {
		return "", err
	}
	for _, iface := range ifaces {
		if iface.Flags&net.FlagUp == 0 {
			continue // interface down
		}
		if iface.Flags&net.FlagLoopback != 0 {
			continue // loopback interface
		}
		addrs, err := iface.Addrs()
		if err != nil {
			return "", err
		}
		for _, addr := range addrs {
			var ip net.IP
			switch v := addr.(type) {
			case *net.IPNet:
				ip = v.IP
			case *net.IPAddr:
				ip = v.IP
			}
			if ip == nil || ip.IsLoopback() {
				continue
			}
			ip = ip.To4()
			if ip == nil {
				continue // not an IPv4 address
			}
			return ip.String(), nil
		}
	}
	return "", errors.New("are you connected to the network?")
}
