package remotes

// StoreInS3 handles storing ReShifter archive in S3 compatible storage.
import (
	"fmt"

	log "github.com/Sirupsen/logrus"
	"github.com/mhausenblas/reshifter/pkg/types"
	"github.com/mhausenblas/reshifter/pkg/util"
	minio "github.com/minio/minio-go"
)

// StoreInS3 stores backup with backupid (in directory target)
// in bucket in an S3 compatible storage, using s3endpoint.
func StoreInS3(s3endpoint, bucket, target, backupid string) error {
	target += ".zip"
	accesskey, secret, err := util.S3CredFromEnv()
	if err != nil {
		return fmt.Errorf("No S3 credentials found: %s", err)
	}
	object := backupid + ".zip"
	mc, err := minio.New(s3endpoint, accesskey, secret, true)
	if err != nil {
		return fmt.Errorf(fmt.Sprintf("%s ", err))
	}
	// Note: we don't care about the error, that is,
	// if the bucket already exists, we ignore that fact.
	// Also, the region doesn't matter, it's a global resource.
	_ = mc.MakeBucket(bucket, "us-east-1")
	// if err != nil {
	// 	exists, berr := mc.BucketExists(bucket)
	// 	if berr == nil && exists {
	// 		return fmt.Errorf(fmt.Sprintf("Bucket %s already exists", bucket))
	// 	}
	// 	return fmt.Errorf(fmt.Sprintf("%s", err))
	// }
	nbytes, err := mc.FPutObject(bucket, object, target, types.ContentTypeZip)
	if err != nil {
		return fmt.Errorf(fmt.Sprintf("%s", err))
	}
	log.WithFields(log.Fields{"func": "remotes.StoreInS3"}).Debug(fmt.Sprintf("Successfully stored %s/%s (%d Bytes) in S3 compatible remote storage %s", bucket, object, nbytes, s3endpoint))
	return nil
}
