package cmd

import (
	"github.com/spf13/cobra"
	"log"
	"runtimectl/config"
	"runtimectl/dao"
	client "runtimectl/pkg/k8s"
)

func newPatchCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "patch",
		Short: "patch the devbox pod",
		RunE:  patchAction,
	}

	cmd.Flags().StringVar(&k8sConfig, "k8s-config", "/root/.kube/config", "the k8s config file path")

	return cmd
}

func patchAction(cmd *cobra.Command, args []string) error {
	config.Init()
	dao.Init()
	k8sClient := client.Init(k8sConfig)
	if err := k8sClient.Patch(); err != nil {
		log.Println("Error patching devbox pod: ", err)
	}
	return nil
}
