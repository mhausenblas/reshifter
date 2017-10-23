package cmd

import (
	"fmt"
	"strings"
	"errors"

	log "github.com/Sirupsen/logrus"
	"github.com/mhausenblas/reshifter/pkg/discovery"
	"github.com/mhausenblas/reshifter/pkg/types"
	"github.com/spf13/cobra"
)

// exploreCmd represents the explore command
var exploreCmd = &cobra.Command{
	Use:   "explore",
	Short: "Probes an etcd endpoint",
	Long:  `Probes an etc endpoint at path /version to figure which version of etcd it is and in which mode (secure or insecure) it is used as well as if a Kubernetes distro can be detected`,
	RunE: func(cmd *cobra.Command, args []string) (error) {
		ep := cmd.Flag("endpoint").Value.String()
		fmt.Printf("Exploring etcd endpoint %s\n", ep)
		return doexplore(ep)
	},
}

func init() {
	RootCmd.AddCommand(exploreCmd)
	exploreCmd.Flags().StringP("endpoint", "e", "http://127.0.0.1:2379", "The URL of the etcd to probe")
}

func doexplore(endpoint string) (error) {
	if endpoint == "" || strings.Index(endpoint, "http") != 0 {
		merr := "The endpoint is malformed"
		log.Error(merr)
		return errors.New(merr)
	}
	version, apiversion, issecure, err := discovery.ProbeEtcd(endpoint)
	if err != nil {
		log.Errorf(fmt.Sprintf("%s", err))
		return err
	}
	distrotype, err := discovery.ProbeKubernetesDistro(endpoint)
	if err != nil {
		log.Errorf(fmt.Sprintf("Can't determine Kubernetes distro: %s", err))
		return err
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
	fmt.Printf("etcd version: %s\nAPI version: %s\nSecure: %s\nDistro: %s\n\n", version, apiversion, secure, distro)

	return nil
}
