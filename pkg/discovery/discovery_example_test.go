package discovery

import (
	"fmt"

	log "github.com/Sirupsen/logrus"
)

// Shows how to probe etcd to determine which version is running
// and in which mode (secure or insecure) it is used. It assumes that the
// etcd process is servering on 127.0.0.1:2379.
func ExampleProbeEtcd() {
	// define the URL etcd is available at:
	etcdurl := "http://127.0.0.1:2379"

	// carry out the probe and handle errors as they occur:
	version, secure, err := ProbeEtcd(etcdurl)
	if err != nil {
		log.Fatal(err)
		return
	}
	fmt.Printf("Discovered etcd in version %s, running in secure mode: %t", version, secure)

	// Output:
	// version: 3.1.0, secure: false
}
