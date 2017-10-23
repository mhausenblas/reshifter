package util

import (
	"fmt"
	"time"

	"github.com/coreos/etcd/clientv3"
	"github.com/coreos/etcd/pkg/transport"
	"google.golang.org/grpc"
	"math"
)

// NewClient3 creates an etcd3 client, optionally using SSL/TLS if secure is true.
// The endpoint is an URL such as http://localhost:2379.
func NewClient3(endpoint string, secure bool) (*clientv3.Client, error) {
	if secure { // manually create a SSL/TLS client
		cs, err := newSecureEtcd3Client(endpoint)
		if err != nil {
			return nil, err
		}
		return cs, nil
	}

	dialOptions := createGRPCOptions()

	// create plain HTTP-based client:
	c, err := clientv3.New(clientv3.Config{
		Endpoints:   []string{endpoint},
		DialTimeout: time.Second,
		DialOptions: dialOptions,
	})
	if err != nil {
		return nil, fmt.Errorf("Can't connect to etcd3: %s", err)
	}
	return c, nil
}

func newSecureEtcd3Client(endpoint string) (*clientv3.Client, error) {
	clientcert, clientkey, err := ClientCertAndKeyFromEnv()
	if err != nil {
		return nil, err
	}
	cafile, err := CACertFromEnv()
	if err != nil {
		return nil, err
	}
	tlsInfo := transport.TLSInfo{
		CertFile: clientcert,
		KeyFile:  clientkey,
		CAFile:   cafile,
	}
	tlsConfig, err := tlsInfo.ClientConfig()
	if err != nil {
		return nil, err
	}

	dialOptions := createGRPCOptions()

	cli, err := clientv3.New(clientv3.Config{
		DialOptions: dialOptions,
		Endpoints:   []string{endpoint},
		TLS:         tlsConfig,
	})
	if err != nil {
		return nil, err
	}
	return cli, nil
}

// This sets the sizes for grpc transfers to MaxUint32. This is the largest
// payload supported from the goloang client.
func createGRPCOptions() []grpc.DialOption {
	var dialOptions []grpc.DialOption
	dialOptions = append(dialOptions, grpc.WithDefaultCallOptions(grpc.MaxCallRecvMsgSize(math.MaxUint32)))
	dialOptions = append(dialOptions, grpc.WithDefaultCallOptions(grpc.MaxCallSendMsgSize(math.MaxUint32)))
	return dialOptions
}