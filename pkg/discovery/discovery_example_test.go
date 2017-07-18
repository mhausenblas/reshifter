// +build !example

package discovery

import (
	"fmt"

	log "github.com/Sirupsen/logrus"
	"github.com/mhausenblas/reshifter/pkg/types"
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
	fmt.Printf("Discovered etcd in version %s, running in secure mode: %t\n", version, secure)

	// Output:
	// Discovered etcd in version 3.1.0, running in secure mode: false
}

// Shows how to probe etcd to determine which Kubernetes distribution is
// present (if at all) by scanning the available keys. It assumes that the
// etcd process is servering on 127.0.0.1:2379 and the OpenShift Kubernetes
// distro is using it.
func ExampleProbeKubernetesDistro() {
	// define the URL etcd is available at:
	etcdurl := "http://127.0.0.1:2379"

	// a textual description of the Kubernetes distro
	var distro string

	// carry out the probe and handle errors as they occur:
	distrotype, err := ProbeKubernetesDistro(etcdurl)
	if err != nil {
		log.Fatal(err)
		return
	}
	switch distrotype {
	case types.Vanilla:
		distro = "Vanilla Kubernetes"
	case types.OpenShift:
		distro = "OpenShift"
	default:
		distro = "no Kubernetes distro found"
	}
	fmt.Printf("Discovered Kubernetes distro: %s\n", distro)

	// Output:
	// Discovered Kubernetes distro: OpenShift
}

// Shows how to glean stats from etcd by walking well-known keys
// of a given Kubernetes distro. It assumes that the etcd process
// is servering on 127.0.0.1:2379 and the vanilla Kubernetes
// distro is using it.
func ExampleCountKeysFor() {
	// define the URL etcd is available at:
	etcdurl := "http://127.0.0.1:2379"

	// carry out the stats walk and handle errors as they occur:
	numkeys, totalsize, err := CountKeysFor(etcdurl, types.LegacyKubernetesPrefix, types.LegacyKubernetesPrefixLast)
	if err != nil {
		log.Fatal(err)
		return
	}
	fmt.Printf("Found %d keys and overall payload size of %d bytes.\n", numkeys, totalsize)

	// Output:
	// Found 1024 keys and overall payload size of 350876 bytes.
}
