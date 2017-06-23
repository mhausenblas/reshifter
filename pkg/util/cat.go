package util

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"time"
)

// Functions that provide container-assisted testing (CAT)

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
