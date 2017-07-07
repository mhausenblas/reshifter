package cmd

import (
	"fmt"
	"io/ioutil"
	"path/filepath"
	"regexp"
	"time"

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
	Run: func(cmd *cobra.Command, args []string) {
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
		}
		fmt.Printf("Successfully created backup: %s/%s.zip\n", target, bid)
		if remote != "" {
			fmt.Printf("Pushed to remote %s in bucket %s\n\n", remote, bucket)
		}
	},
}

var listBackupCmd = &cobra.Command{
	Use:   "list",
	Short: "Lists backups of a Kubernetes cluster",
	Run: func(cmd *cobra.Command, args []string) {
		files, err := ioutil.ReadDir(types.DefaultWorkDir)
		if err != nil {
			log.Error(err)
			return
		}
		for _, file := range files {
			re := regexp.MustCompile("\\d{10}.zip")
			fn := file.Name()
			bid := fn[0 : len(fn)-len(filepath.Ext(fn))]
			if re.Match([]byte(fn)) {
				fmt.Printf("%s\n", bid)
			}
		}
	},
}

func init() {
	RootCmd.AddCommand(backupCmd)
	backupCmd.AddCommand(createBackupCmd)
	backupCmd.AddCommand(listBackupCmd)
	createBackupCmd.Flags().StringP("endpoint", "e", "http://127.0.0.1:2379", "The URL of the etcd to use")
	createBackupCmd.Flags().StringP("target", "t", "/tmp", "Optionally, the target directory for the resulting ZIP file of the backup")
	createBackupCmd.Flags().StringP("remote", "r", "", "Optionally, the S3-compatible storage endpoint")
	createBackupCmd.Flags().StringP("bucket", "b", "", "Optionally, the target bucket in the S3-compatible storage endpoint")
	_ = createBackupCmd.MarkFlagRequired("endpoint")
}
