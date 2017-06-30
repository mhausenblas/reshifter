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
	"strings"

	log "github.com/Sirupsen/logrus"
	"github.com/mhausenblas/reshifter/pkg/discovery"
	"github.com/mhausenblas/reshifter/pkg/types"
	"github.com/spf13/cobra"
)

// exploreCmd represents the explore command
var exploreCmd = &cobra.Command{
	Use:   "explore",
	Short: "Probes an etcd endpoint",
	Long:  `Probes an etc endpoint at path /version to figure which version of etcd it is and in which mode (secure or insecure) it is used`,
	Run: func(cmd *cobra.Command, args []string) {
		ep := cmd.Flag("endpoint").Value.String()
		fmt.Printf("Exploring etcd endpoint %s\n", ep)
		doexplore(ep)
	},
}

func init() {
	RootCmd.AddCommand(exploreCmd)
	exploreCmd.Flags().StringP("endpoint", "e", "http://127.0.0.1:2379", "The URL of the etcd to probe.")
}

func doexplore(endpoint string) {
	if endpoint == "" || strings.Index(endpoint, "http") != 0 {
		merr := "The endpoint is malformed"
		log.Error(merr)
		return
	}
	version, issecure, err := discovery.ProbeEtcd(endpoint)
	if err != nil {
		log.Errorf(fmt.Sprintf("%s", err))
		return
	}
	distrotype, err := discovery.ProbeKubernetesDistro(endpoint)
	if err != nil {
		log.Errorf(fmt.Sprintf("Can't determine Kubernetes distro: %s", err))
		return
	}
	secure := "insecure etcd, no SSL/TLS configured"
	if issecure {
		secure = "secure etcd, SSL/TLS configure"
	}
	var distro string
	switch distrotype {
	case types.Vanilla:
		distro = "Vanilla Kubernetes"
	case types.OpenShift:
		distro = "OpenShift"
	default:
		distro = "no Kubernetes distro found"
	}
	fmt.Printf("Version: %s\nSecure: %s\nDistro: %s\n\n", version, secure, distro)
}
