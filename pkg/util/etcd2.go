package util

import (
	"context"
	"crypto/tls"
	"crypto/x509"
	"fmt"
	"io/ioutil"
	"net/http"
	"time"

	"go.etcd.io/etcd/client"
)

// NewClient2 creates an etcd2 client, optionally using SSL/TLS if secure is true.
// The endpoint is an URL such as http://localhost:2379.
func NewClient2(endpoint string, secure bool) (client.Client, error) {
	if secure { // manually create a SSL/TLS client
		tr, err := etcd2transport()
		if err != nil {
			return nil, err
		}
		cs, err := newSecureEtcd2Client(tr, endpoint)
		if err != nil {
			return nil, err
		}
		return cs, nil
	}
	// create plain HTTP-based client:
	c, err := client.New(client.Config{
		Endpoints:               []string{endpoint},
		Transport:               client.DefaultTransport,
		HeaderTimeoutPerRequest: time.Second,
	})
	if err != nil {
		return nil, fmt.Errorf("Can't connect to etcd2: %s", err)
	}
	return c, nil
}

func newSecureEtcd2Client(tr *http.Transport, endpoint string) (client.Client, error) {
	cli, err := client.New(client.Config{
		Endpoints: []string{endpoint},
		Transport: tr,
	})
	if err != nil {
		return nil, err
	}
	return cli, nil
}

func etcd2transport() (*http.Transport, error) {
	clientcert, clientkey, err := ClientCertAndKeyFromEnv()
	if err != nil {
		return nil, err
	}
	cert, err := tls.LoadX509KeyPair(clientcert, clientkey)
	if err != nil {
		return nil, err
	}
	cafile, err := CACertFromEnv()
	if err != nil {
		return nil, err
	}
	cacert, err := ioutil.ReadFile(cafile)
	if err != nil {
		return nil, err
	}
	cacertpool := x509.NewCertPool()
	cacertpool.AppendCertsFromPEM(cacert)
	tr := &http.Transport{
		TLSClientConfig: &tls.Config{
			Certificates:       []tls.Certificate{cert},
			RootCAs:            cacertpool,
			InsecureSkipVerify: true,
		},
	}
	return tr, nil
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
