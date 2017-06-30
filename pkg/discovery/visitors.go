package discovery

import (
	"fmt"

	log "github.com/Sirupsen/logrus"
	"github.com/coreos/etcd/client"
	"github.com/coreos/etcd/clientv3"
	"github.com/mhausenblas/reshifter/pkg/types"
	"golang.org/x/net/context"
)

// Visit2 recursively visits an etcd2 server from the root and applies
// the reap function on leaf nodes (keys that don't have sub-keys),
// otherwise descents the tree.
func Visit2(kapi client.KeysAPI, path string, fn types.Reap) error {
	log.WithFields(log.Fields{"func": "discovery.Visit2"}).Debug(fmt.Sprintf("Processing %s", path))
	copts := client.GetOptions{
		Recursive: false,
		Quorum:    false,
	}
	res, err := kapi.Get(context.Background(), path, &copts)
	if err != nil {
		return err
	}
	if res.Node.Dir { // there are children
		log.WithFields(log.Fields{"func": "discovery.Visit2"}).Debug(fmt.Sprintf("%s has %d children", path, len(res.Node.Nodes)))
		for _, node := range res.Node.Nodes {
			log.WithFields(log.Fields{"func": "discovery.Visit2"}).Debug(fmt.Sprintf("Next visiting child %s", node.Key))
			_ = Visit2(kapi, node.Key, fn)
		}
		return nil
	}
	// otherwise we're on a leaf node:
	return fn(res.Node.Key, string(res.Node.Value))
}

// Visit3 visits all paths of an etcd3 server and applies the reap function
// on the keys.
func Visit3(c3 *clientv3.Client, path string, distro types.KubernetesDistro, fn types.Reap) error {
	log.WithFields(log.Fields{"func": "discovery.Visit3"}).Debug(fmt.Sprintf("Processing %s", path))
	endkey := ""
	if distro == types.Vanilla {
		endkey = types.KubernetesPrefixLast
	}
	if distro == types.OpenShift {
		endkey = types.OpenShiftPrefixLast
	}
	res, err := c3.Get(context.Background(), path+"/*", clientv3.WithRange(endkey))
	// res, err := c3.Get(context.Background(), "/kubernetes.io/namespaces/kube-system")
	if err != nil {
		return err
	}
	log.WithFields(log.Fields{"func": "discovery.Visit3"}).Debug(fmt.Sprintf("Got %v", res))
	for _, ev := range res.Kvs {
		log.WithFields(log.Fields{"func": "discovery.Visit3"}).Debug(fmt.Sprintf("key: %s, value: %s", ev.Key, ev.Value))
		err = fn(string(ev.Key), string(ev.Value))
		if err != nil {
			return err
		}
	}
	return nil
}
