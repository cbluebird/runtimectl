package cmd

import (
	"fmt"
	"github.com/spf13/cobra"
	"runtimectl/config"
	"runtimectl/dao"
	"runtimectl/model"
	"runtimectl/util"
)

func newSyncCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "sync",
		Short: "sync the database",
		RunE:  syncAction,
	}
}

func syncAction(cmd *cobra.Command, args []string) error {
	config.Init()
	dao.Init()
	config := util.ParseJson()
	if len(config.Runtime.Framework) != 0 {
		if err := sync("FRAMEWORK", config.Runtime.Framework); err != nil {
			fmt.Println("Error syncing framework:", err)
			return err
		}
	}
	if len(config.Runtime.Language) != 0 {
		err := sync("LANGUAGE", config.Runtime.Language)
		if err != nil {
			return err
		}
	}
	if len(config.Runtime.OS) != 0 {
		err := sync("OS", config.Runtime.OS)
		if err != nil {
			return err
		}
	}
	if len(config.Runtime.Custom) != 0 {
		err := sync("CUSTOM", config.Runtime.Custom)
		if err != nil {
			return err
		}
	}
	return nil
}

func sync(kind string, runtime []model.RuntimeVersion) error {
	for _, o := range runtime {
		if err := dao.CreateOrUpdateTemplateRepository(o.Name, kind); err != nil {
			fmt.Println("Error creating or updating template repository:", err)
			return err
		}
		t := dao.GetTemplateRepository(o.Name)
		for _, version := range o.Version {
			if err := dao.CreateOrUpdateTemplate(version.Name, t.UID, version.Image, version.Config); err != nil {
				fmt.Println("Error creating or updating template:", err)
				return err
			}
		}
	}
	return nil
}
