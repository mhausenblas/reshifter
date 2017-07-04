package backup

import (
	"fmt"
	"os"

	log "github.com/Sirupsen/logrus"
)

// ExampleBackup shows how to back up a Kubernetes cluster
// by specifying the underlying etcd. It assumes that the
// etcd process is servering on 127.0.0.1:2379.
func ExampleBackup() {

	// define the port etcd is listening on:
	port := "2379"

	// define the URL etcd is available at:
	etcdurl := "http://127.0.0.1:" + port

	// using Minio play, a public S3-compatible sandbox
	// for the remote storage available at https://play.minio.io:9000
	// and the following credentials:
	_ = os.Setenv("ACCESS_KEY_ID", "Q3AM3UQ867SPQQA43P2F")
	_ = os.Setenv("SECRET_ACCESS_KEY", "zuf+tfteSlswRu7BJ86wekitnifILbZam1KYY3TG")

	// carry out the backup of etcd underlying the Kubernetes cluster
	// and handle errors as they occur:
	backupid, err := Backup(etcdurl, "/tmp", "play.minio.io:9000", "2017-07-some-bucket")
	if err != nil {
		log.Fatal(err)
		return
	}
	fmt.Printf("The backup completed successfully with ID %s", backupid)

	// Output:
	// backupid: "1498847078"
}
