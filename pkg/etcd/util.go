package etcd

import (
	"context"
	"fmt"
	"os/exec"
	"time"

	"github.com/coreos/etcd/client"
)

// newClient2 create an etcd2 client, optionally using SSL/TLS if secure is true.
// The endpoint is both host and port, for example, localhost:2379.
func newClient2(endpoint string, secure bool) (client.Client, error) {

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

// setKV2 sets the key with val in an etcd2 cluster and
// iff val is empty, creates a directory key.
func setKV2(kapi client.KeysAPI, key, val string) error {
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

func etcd2up(port string) error {
	// var out bytes.Buffer
	cmd := exec.Command("docker", "run", "--rm", "-d",
		"-p", port+":"+port, "--name", "test-etcd", "quay.io/coreos/etcd:v2.3.8",
		"--advertise-client-urls", "http://0.0.0.0:"+port,
		"--listen-client-urls", "http://0.0.0.0:"+port)
	// cmd.Stdout = &out
	fmt.Printf("%s\n", cmd.Args)
	err := cmd.Run()
	if err != nil {
		return err
	}
	// fmt.Printf("%s\n", out.String())
	// time.Sleep(time.Second * 2)
	return nil
}

func etcd3up(port string) error {
	// var out bytes.Buffer
	cmd := exec.Command("docker", "run", "--rm", "-d",
		"-p", port+":"+port, "--name", "test-etcd",
		"quay.io/coreos/etcd:v3.1.0", "/usr/local/bin/etcd",
		"--advertise-client-urls", "http://0.0.0.0:"+port,
		"--listen-client-urls", "http://0.0.0.0:"+port)
	// cmd.Stdout = &out
	fmt.Printf("%s\n", cmd.Args)
	err := cmd.Run()
	if err != nil {
		return err
	}
	// fmt.Printf("%s\n", out.String())
	// time.Sleep(time.Second * 2)
	return nil
}

func etcddown() error {
	cmd := exec.Command("docker", "kill", "test-etcd")
	err := cmd.Run()
	if err != nil {
		return err
	}
	return nil
}
