package cmd

import (
	"fmt"
	"time"

	log "github.com/Sirupsen/logrus"
	"github.com/mhausenblas/reshifter/pkg/restore"
	"github.com/mhausenblas/reshifter/pkg/util"
	"github.com/spf13/cobra"
)

// restoreCmd represents the restore command
var restoreCmd = &cobra.Command{
	Use:   "restore",
	Short: "Performs a restore of a Kubernetes cluster",
	Long:  `Performs a restore of a Kubernetes cluster from the content of a ZIP file, either local or from an S3-compatible remote storage`,
	Run: func(cmd *cobra.Command, args []string) {
		ep := cmd.Flag("endpoint").Value.String()
		bid := cmd.Flag("backupid").Value.String()
		target := cmd.Flag("target").Value.String()
		remote := cmd.Flag("remote").Value.String()
		bucket := cmd.Flag("bucket").Value.String()

		if !util.IsBackupID(bid) {
			abortreason := fmt.Sprintf("Aborting restore: %s is not a valid backup ID", bid)
			log.Error(abortreason)
			return
		}
		if remote != "" && bucket == "" {
			bucket = "reshifter-" + time.Now().UTC().Format("2006-01-02")
			fmt.Printf("You didn't tell me which bucket to use, using %s as a fallback\n", bucket)
		}
		krestored, etime, err := restore.Restore(ep, bid, target, remote, bucket)
		if err != nil {
			log.Error(err)
		}
		if remote != "" {
			fmt.Printf("Using source: remote %s, bucket %s\n", remote, bucket)
		}
		fmt.Printf("Successfully restored %d key(s) from %s.zip in %f sec\n\n", krestored, bid, etime.Seconds())
	},
}

func init() {
	RootCmd.AddCommand(restoreCmd)

	restoreCmd.Flags().StringP("endpoint", "e", "http://127.0.0.1:2379", "The URL of the etcd to use")
	restoreCmd.Flags().StringP("backupid", "i", "", "The ID of the backup to use for the restore operation")
	restoreCmd.Flags().StringP("target", "t", "/tmp", "Optionally, the target directory for the resulting ZIP file of the backup")
	restoreCmd.Flags().StringP("remote", "r", "", "Optionally, the S3-compatible storage endpoint")
	restoreCmd.Flags().StringP("bucket", "b", "", "Optionally, the target bucket in the S3-compatible storage endpoint")
}
