// Copyright © 2017 NAME HERE <EMAIL ADDRESS>
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package cmd

import (
	"fmt"
	"time"

	log "github.com/Sirupsen/logrus"
	"github.com/mhausenblas/reshifter/pkg/backup"
	"github.com/spf13/cobra"
)

// backupCmd represents the backup command
var backupCmd = &cobra.Command{
	Use:   "backup",
	Short: "Creates a backup of a Kubernetes cluster",
	Long:  `Creates a backup of a Kubernetes cluster by travesing the underlying etcd and storing the content in a ZIP file, either local or in an S3-compatible remote storage`,
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

func init() {
	RootCmd.AddCommand(backupCmd)
	backupCmd.Flags().StringP("endpoint", "e", "http://127.0.0.1:2379", "The URL of the etcd to use")
	backupCmd.Flags().StringP("target", "t", "/tmp", "Optionally, the target directory for the resulting ZIP file of the backup")
	backupCmd.Flags().StringP("remote", "r", "", "Optionally, the S3-compatible storage endpoint")
	backupCmd.Flags().StringP("bucket", "b", "", "Optionally, the target bucket in the S3-compatible storage endpoint")
}
