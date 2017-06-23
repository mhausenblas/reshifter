package util

import (
	"context"
	"fmt"
	"os/exec"
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
// The endpoint is both host and port, for example, localhost:2379.
func NewClient2(endpoint string, secure bool) (client.Client, error) {
	if secure {
		return nil, fmt.Errorf("Secure etcd2 connection not yet supported")
	}
	// create plain HTTP-based client:
	cfg := client.Config{
		Endpoints:               []string{"http://" + endpoint},
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

// Etcd2up launches an etcd2 server on port.
func Etcd2up(port string) error {
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

// Etcd3up launches an etcd3 server on port.
func Etcd3up(port string) error {
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

// Etcddown tears down an etcd server.
func Etcddown() error {
	cmd := exec.Command("docker", "kill", "test-etcd")
	err := cmd.Run()
	if err != nil {
		return err
	}
	time.Sleep(time.Second * 1)
	return nil
}
