package cmd

import (
	"fmt"
	"github.com/spf13/cobra"
	"runtimectl/config"
	"runtimectl/dao"
	"runtimectl/model"
	"runtimectl/pkg/util"
	"time"
)

func newSyncCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "sync",
		Short: "sync the database",
		RunE:  syncAction,
	}

	cmd.Flags().StringVar(&path, "path", "config.json", "Path to the config file")
	return cmd
}

func syncAction(cmd *cobra.Command, args []string) error {
	config.Init()
	dao.Init()
	config := util.ParseJson(path)
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
		if err := dao.CreateOrUpdateTemplateRepository(o.Name, kind, ""); err != nil {
			fmt.Println("Error creating or updating template repository:", err)
			return err
		}
		t := dao.GetTemplateRepository(o.Name)

		var activeVersion model.Version

		for _, version := range o.Version {
			if version.State == "active" {
				activeVersion = version
				break
			}
		}

		if activeVersion.Name == "" || activeVersion.Image == "" {
			continue
		}

		if err := dao.DeprecateTemplates(t.UID); err != nil {
			return err
		}

		if err := dao.CreateOrUpdateTemplate(activeVersion.Name, t.UID, activeVersion.Image, activeVersion.Config, "active", time.Now()); err != nil {
			fmt.Println("Error creating or updating active template:", err)
			return err
		}
	}
	return nil
}
