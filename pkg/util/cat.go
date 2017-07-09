package util

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"time"
)

// Functions that provide container-assisted testing (CAT)

// LaunchEtcd2 launches etcd in v2 on port, either in secure
// or in insecure mode, depending on the scheme used in tetcd.
func LaunchEtcd2(tetcd, port string) (bool, error) {
	secure := false
	switch {
	case strings.Index(tetcd, "https") == 0:
		err := etcd2SecureUp(port)
		secure = true
		_ = os.Setenv("RS_ETCD_CLIENT_CERT", filepath.Join(Certsdir(), "client.pem"))
		_ = os.Setenv("RS_ETCD_CLIENT_KEY", filepath.Join(Certsdir(), "client-key.pem"))
		_ = os.Setenv("RS_ETCD_CA_CERT", filepath.Join(Certsdir(), "ca.pem"))
		if err != nil {
			return secure, fmt.Errorf("Can't launch secure etcd2 at %s: %s", tetcd, err)
		}
	case strings.Index(tetcd, "http") == 0:
		err := etcd2Up(port)
		if err != nil {
			return secure, fmt.Errorf("Can't launch insecure etcd2 at %s: %s", tetcd, err)
		}
	default:
		return secure, fmt.Errorf("That's not a valid etcd2 endpoint: %s", tetcd)
	}
	return secure, nil
}

// LaunchEtcd3 launches etcd in v3 on port, either in secure
// or in insecure mode, depending on the scheme used in tetcd.
func LaunchEtcd3(tetcd, port string) (bool, error) {
	secure := false
	switch {
	case strings.Index(tetcd, "https") == 0:
		err := etcd3SecureUp(port)
		secure = true
		_ = os.Setenv("RS_ETCD_CLIENT_CERT", filepath.Join(Certsdir(), "client.pem"))
		_ = os.Setenv("RS_ETCD_CLIENT_KEY", filepath.Join(Certsdir(), "client-key.pem"))
		_ = os.Setenv("RS_ETCD_CA_CERT", filepath.Join(Certsdir(), "ca.pem"))
		if err != nil {
			return secure, fmt.Errorf("Can't launch secure etcd2 at %s: %s", tetcd, err)
		}
	case strings.Index(tetcd, "http") == 0:
		err := etcd3Up(port)
		if err != nil {
			return secure, fmt.Errorf("Can't launch insecure etcd3 at %s: %s", tetcd, err)
		}
	default:
		return secure, fmt.Errorf("That's not a valid etcd3 endpoint: %s", tetcd)

	}
	return secure, nil
}

// etcd2Up launches an etcd2 server on port.
func etcd2Up(port string) error {
	cmd := exec.Command("docker", "run", "--rm", "-d",
		"-p", port+":"+port, "--name", "test-etcd", "--dns", "8.8.8.8",
		"quay.io/coreos/etcd:v2.3.8",
		"--advertise-client-urls", "http://0.0.0.0:"+port,
		"--listen-client-urls", "http://0.0.0.0:"+port)
	fmt.Printf("%s\n", cmd.Args)
	err := cmd.Run()
	if err != nil {
		return err
	}
	time.Sleep(time.Second * 2)
	return nil
}

// etcd2SecureUp launches a secure etcd2 server on port.
func etcd2SecureUp(port string) error {
	cmd := exec.Command("docker", "run", "--rm", "-d",
		"-v", Certsdir()+"/:/etc/ssl/certs", "-p", port+":"+port,
		"--name", "test-etcd", "--dns", "8.8.8.8",
		"quay.io/coreos/etcd:v2.3.8",
		"--ca-file", "/etc/ssl/certs/ca.pem",
		"--cert-file", "/etc/ssl/certs/server.pem",
		"--key-file", "/etc/ssl/certs/server-key.pem",
		"--advertise-client-urls", "https://0.0.0.0:"+port,
		"--listen-client-urls", "https://0.0.0.0:"+port)
	fmt.Printf("%s\n", cmd.Args)
	err := cmd.Run()
	if err != nil {
		return err
	}
	time.Sleep(time.Second * 2)
	return nil
}

// etcd3Up launches an etcd3 server on port.
func etcd3Up(port string) error {
	cmd := exec.Command("docker", "run", "--rm", "-d",
		"-p", port+":"+port, "--name", "test-etcd", "--dns", "8.8.8.8",
		"quay.io/coreos/etcd:v3.1.0", "/usr/local/bin/etcd",
		"--advertise-client-urls", "http://0.0.0.0:"+port,
		"--listen-client-urls", "http://0.0.0.0:"+port)
	fmt.Printf("%s\n", cmd.Args)
	err := cmd.Run()
	if err != nil {
		return err
	}
	time.Sleep(time.Second * 2)
	return nil
}

// etcd3SecureUp launches a secure etcd3 server on port.
func etcd3SecureUp(port string) error {
	cmd := exec.Command("docker", "run", "--rm", "-d",
		"-v", Certsdir()+"/:/etc/ssl/certs", "-p", port+":"+port,
		"--name", "test-etcd", "--dns", "8.8.8.8",
		"quay.io/coreos/etcd:v3.1.0", "/usr/local/bin/etcd",
		"--ca-file", "/etc/ssl/certs/ca.pem",
		"--cert-file", "/etc/ssl/certs/server.pem",
		"--key-file", "/etc/ssl/certs/server-key.pem",
		"--advertise-client-urls", "https://0.0.0.0:"+port,
		"--listen-client-urls", "https://0.0.0.0:"+port)
	fmt.Printf("%s\n", cmd.Args)
	err := cmd.Run()
	if err != nil {
		return err
	}
	time.Sleep(time.Second * 2)
	return nil
}

// EtcdDown tears down an etcd server.
func EtcdDown() error {
	cmd := exec.Command("docker", "kill", "test-etcd")
	err := cmd.Run()
	if err != nil {
		return err
	}
	time.Sleep(time.Second * 1)
	return nil
}

// Certsdir returns the absolute path to the directory
// where the pre-generated certs and keys are.
func Certsdir() string {
	base, _ := os.Getwd()
	certsrel := filepath.Join(base, "../../testbed/certs")
	certs, _ := filepath.Abs(certsrel)
	return certs
}
