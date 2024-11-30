package cmd

import (
	"github.com/spf13/cobra"
	"log"
	"runtimectl/config"
	"runtimectl/dao"
	client "runtimectl/pkg/k8s"
)

func newInitCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "init",
		Short: "database init",
		RunE:  initAction,
	}

	cmd.Flags().StringVar(&path, "path", "config.json", "Path to the config file")
	return cmd
}

func initAction(cmd *cobra.Command, args []string) error {
	config.Init()
	dao.Init()
	k8sClient := client.Init(k8sConfig)
	if err := k8sClient.SyncToDB(); err != nil {
		log.Println(err)
	}
	return nil
}
