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

// statsCmd represents the stats command
var statsCmd = &cobra.Command{
	Use:   "stats",
	Short: "Collects stats about Kubernetes-related keys from an etcd endpoint",
	RunE: func(cmd *cobra.Command, args []string) (error) {
		ep := cmd.Flag("endpoint").Value.String()
		fmt.Printf("Collecting stats from etcd endpoint %s\n", ep)
		return docollectstats(ep)
	},
}

func init() {
	RootCmd.AddCommand(statsCmd)
	statsCmd.Flags().StringP("endpoint", "e", "http://127.0.0.1:2379", "The URL of the etcd to collect stats from")
}

func docollectstats(endpoint string) (error) {
	if endpoint == "" || strings.Index(endpoint, "http") != 0 {
		merr := "The endpoint is malformed"
		log.Error(merr)
		return errors.New(merr)
	}
	_, _, _, err := discovery.ProbeEtcd(endpoint)
	if err != nil {
		log.Errorf(fmt.Sprintf("%s", err))
		return err
	}
	vlk, vls, derr := discovery.CountKeysFor(endpoint, types.LegacyKubernetesPrefix, types.LegacyKubernetesPrefixLast)
	if derr != nil {
		log.Info("Didn't find legacy keys, trying modern keys now …")
	}
	vk, vs, err := discovery.CountKeysFor(endpoint, types.KubernetesPrefix, types.KubernetesPrefixLast)
	if err != nil && derr != nil {
		serr := fmt.Sprintf("Having problems calculating stats: %s", err)
		log.Error(serr)
		return errors.New(serr)
	}
	fmt.Printf("Vanilla Kubernetes [keys:%d, size:%d]\n", vlk+vk, vls+vs)
	osk, oss, _ := discovery.CountKeysFor(endpoint, types.OpenShiftPrefix, types.OpenShiftPrefixLast)
	if osk > 0 {
		fmt.Printf("OpenShift [keys:%d, size:%d]\n\n", osk, oss)
	}

	return nil
}
