package etcd

import (
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"
	"testing"

	"github.com/coreos/etcd/client"
)

var (
	tmpTestDir = "test/"
	storetests = []struct {
		path string
		val  string
	}{
		{"", ""},
		{"non-valid-key", ""},
		{"/", "root"},
		{"/" + tmpTestDir, "some"},
		{"/" + tmpTestDir + "/first-level", "another"},
		{"/" + tmpTestDir + "/this:also", "escaped"},
	}
)

func TestBackup(t *testing.T) {
	port := "2379"
	tetcd := "localhost:" + port
	err := etcd2up(port)
	if err != nil {
		t.Errorf("Can't launch local etcd at %s: %s", tetcd, err)
		return
	}
	// create some key-value pairs:
	c2, err := newClient2(tetcd, false)
	if err != nil {
		t.Errorf("Can't connect to local etcd2 at %s: %s", tetcd, err)
		return
	}

	kapi := client.NewKeysAPI(c2)
	err = setKV2(kapi, "/foo", "some")
	if err != nil {
		t.Errorf("Can't create key /foo: %s", err)
		return
	}
	err = setKV2(kapi, "/that", "")
	if err != nil {
		t.Errorf("Can't create key /that: %s", err)
		return
	}
	err = setKV2(kapi, "/that/here", "moar")
	if err != nil {
		t.Errorf("Can't create key /that/here: %s", err)
	}

	based, err := Backup(tetcd)
	if err != nil {
		t.Errorf("Error during backup: %s", err)
	}
	// TODO: check if content is as expected
	_, err = os.Stat(based + ".zip")
	if err != nil {
		t.Errorf("No archive found: %s", err)
	}
	// make sure to clean up:
	_ = os.Remove(based + ".zip")
	_ = etcddown()
}

func TestStore(t *testing.T) {
	for _, tt := range storetests {
		p, err := store(".", tt.path, tt.val)
		if err != nil {
			continue
		}
		got := readcontent(p)
		want := tt.val
		if got != want {
			t.Errorf("etcd.store(\".\", %q, %q) => %q, want %q", tt.path, tt.val, got, want)
		}
	}
	// make sure to clean up remaining directories:
	_ = os.RemoveAll(tmpTestDir)
}

func readcontent(path string) string {
	// make sure to clean up individual files
	defer func() {
		if path != "." {
			_ = os.Remove(path)
		}
	}()
	content, _ := ioutil.ReadFile(path)
	return string(content)
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
