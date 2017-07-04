package restore

import (
	"fmt"
	"os"

	log "github.com/Sirupsen/logrus"
)

// Shows how to restore a Kubernetes cluster
// by specifying the underlying etcd. It assumes that the
// etcd process is servering on 127.0.0.1:2379.
// Takes the backup from Minio play, a public S3-compatible storage sandbox.
func ExampleRestore() {
	// define the URL etcd is available at:
	etcdurl := "http://127.0.0.1:2379"

	// using Minio play, a public S3-compatible sandbox
	// for the remote storage available at https://play.minio.io:9000
	// and the following credentials which need to be exposed
	// as environment variables to the ReShifter process:
	_ = os.Setenv("ACCESS_KEY_ID", "Q3AM3UQ867SPQQA43P2F")
	_ = os.Setenv("SECRET_ACCESS_KEY", "zuf+tfteSlswRu7BJ86wekitnifILbZam1KYY3TG")

	// carry out the restore into etcd underlying the Kubernetes cluster
	// and handle errors as they occur:
	keysrestored, err := Restore(etcdurl, "1498847078", "/tmp", "play.minio.io:9000", "2017-07-some-bucket")
	if err != nil {
		log.Fatal(err)
		return
	}
	fmt.Printf("The restore completed successfully. Restored %d keys", keysrestored)

	// Output:
	// keysrestored: 1042
}
