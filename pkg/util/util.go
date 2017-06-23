package util

import (
	"context"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"regexp"
	"time"

	"github.com/coreos/etcd/client"
)

// IsBackupID tests if a string is a valid backup ID.
// A valid backup ID is a 10 digit integer, representing
// a Unix timestamp.
func IsBackupID(id string) bool {
	re := regexp.MustCompile("\\d{10}")
	return re.Match([]byte(id))
}

// NewClient2 creates an etcd2 client, optionally using SSL/TLS if secure is true.
// The endpoint is an URL such as http://localhost:2379.
func NewClient2(endpoint string, secure bool) (client.Client, error) {
	if secure {
		return nil, fmt.Errorf("Secure etcd2 connection not yet supported")
	}
	// create plain HTTP-based client:
	cfg := client.Config{
		Endpoints:               []string{endpoint},
		Transport:               client.DefaultTransport,
		HeaderTimeoutPerRequest: time.Second,
	}
	c, err := client.New(cfg)
	if err != nil {
		return nil, fmt.Errorf("Can't connect to etcd2: %s", err)
	}
	return c, nil
}

// SetKV2 sets the key with val in an etcd2 cluster and
// iff val is empty, creates a directory key.
func SetKV2(kapi client.KeysAPI, key, val string) error {
	if val == "" {
		_, err := kapi.Set(context.Background(), key, "", &client.SetOptions{Dir: true, PrevExist: client.PrevIgnore})
		if err != nil {
			return fmt.Errorf("Can't set directory key %s: %s", key, err)
		}
		return nil
	}
	_, err := kapi.Set(context.Background(), key, val, &client.SetOptions{Dir: false, PrevExist: client.PrevIgnore})
	if err != nil {
		return fmt.Errorf("Can't set key %s with value %s: %s", key, val, err)
	}
	return nil
}

// Etcd2Up launches an etcd2 server on port.
func Etcd2Up(port string) error {
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
	time.Sleep(time.Second * 1)
	return nil
}

// Etcd2SecureUp launches a secure etcd2 server on port.
func Etcd2SecureUp(port string) error {
	cmd := exec.Command("docker", "run", "--rm", "-d",
		"-v", Certsdir("")+"/:/etc/ssl/certs", "-p", port+":"+port,
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
	time.Sleep(time.Second * 1)
	return nil
}

// Etcd3Up launches an etcd3 server on port.
func Etcd3Up(port string) error {
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
	time.Sleep(time.Second * 1)
	return nil
}

// Etcd3SecureUp launches a secure etcd3 server on port.
func Etcd3SecureUp(port string) error {
	cmd := exec.Command("docker", "run", "--rm", "-d",
		"-v", Certsdir("")+"/:/etc/ssl/certs", "-p", port+":"+port,
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
	time.Sleep(time.Second * 1)
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
func Certsdir(base string) string {
	if base == "" {
		base, _ = os.Getwd()
	}
	certsrel := filepath.Join(base, "../../certs")
	certs, _ := filepath.Abs(certsrel)
	return certs
}
