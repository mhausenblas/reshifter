package cmd

import (
	"fmt"
	"os"
	"time"
	"errors"

	log "github.com/Sirupsen/logrus"
	"github.com/mhausenblas/reshifter/pkg/backup"
	"github.com/mhausenblas/reshifter/pkg/types"
	"github.com/spf13/cobra"
)

var backupCmd = &cobra.Command{
	Use:   "backup",
	Short: "Manages backups of a Kubernetes cluster",
}

var createBackupCmd = &cobra.Command{
	Use:   "create",
	Short: "Creates a backup of a Kubernetes cluster",
	Long:  `Backups are created by travesing the underlying etcd and storing the content in a ZIP file in the local filesystem and optionally in an S3-compatible remote storage`,
	RunE: func(cmd *cobra.Command, args []string) (error) {
		ep := cmd.Flag("endpoint").Value.String()
		target := cmd.Flag("target").Value.String()
		remote := cmd.Flag("remote").Value.String()
		bucket := cmd.Flag("bucket").Value.String()
		if remote != "" && bucket == "" {
			bucket = "reshifter-" + time.Now().UTC().Format("2006-01-02")
			fmt.Printf("You didn't tell me which bucket to use, using %s as a fallback\n", bucket)
		}
		bid, err := backup.Backup(ep, target, remote, bucket)
		if err != nil {
			log.Error(err)
			return err
		}
		if os.Getenv("RS_BACKUP_STRATEGY") != types.ReapFunctionRender {
			fmt.Printf("Successfully created backup: %s/%s.zip\n", target, bid)
			if remote != "" {
				fmt.Printf("Pushed to remote %s in bucket %s\n\n", remote, bucket)
			}
		}

		return nil
	},
}

var listBackupCmd = &cobra.Command{
	Use:   "list",
	Short: "Lists backups of a Kubernetes cluster",
	RunE: func(cmd *cobra.Command, args []string) (error) {
		remote := cmd.Flag("remote").Value.String()
		bucket := cmd.Flag("bucket").Value.String()
		backupIDs, err := backup.List(remote, bucket)
		if err != nil {
			log.Error(err)
			return err
		}
		for _, bid := range backupIDs {
			fmt.Printf("%s\n", bid)
		}

		return nil
	},
}

func init() {
	RootCmd.AddCommand(backupCmd)
	backupCmd.AddCommand(createBackupCmd)
	backupCmd.AddCommand(listBackupCmd)
	backupCmd.PersistentFlags().StringP("remote", "r", "", "Optionally, the S3-compatible storage endpoint")
	backupCmd.PersistentFlags().StringP("bucket", "b", "", "Optionally, the target bucket in the S3-compatible storage endpoint")
	createBackupCmd.Flags().StringP("endpoint", "e", "http://127.0.0.1:2379", "The URL of the etcd to use")
	createBackupCmd.Flags().StringP("target", "t", types.DefaultWorkDir, "Optionally, the target directory for the resulting ZIP file of the backup")
	_ = createBackupCmd.MarkFlagRequired("endpoint")
}
