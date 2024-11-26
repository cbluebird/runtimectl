package cmd

import (
	"encoding/json"
	"fmt"
	"github.com/spf13/cobra"
	"io/ioutil"
	"log"
	"runtimectl/model"
	"runtimectl/util"
)

var (
	kind    string
	name    string
	version string
	image   string
)

func newGenCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "gen",
		Short: "gen the json for database",
		RunE:  genAction,
	}

	cmd.Flags().StringVar(&kind, "kind", "", "Kind of runtime (framework, language, custom, os)")
	cmd.Flags().StringVar(&name, "name", "", "Name of the runtime")
	cmd.Flags().StringVar(&version, "version", "", "Version of the runtime")
	cmd.Flags().StringVar(&image, "image", "", "Image of the runtime")

	return cmd
}

func genAction(cmd *cobra.Command, args []string) error {
	config := util.ParseJson()
	if config == nil {
		return fmt.Errorf("failed to parse config.json")
	}

	runtimeVersions := getRuntimeVersions(config, kind)
	if runtimeVersions == nil {
		return fmt.Errorf("invalid kind: %s", kind)
	}

	updated := false
	for i, rv := range runtimeVersions {
		if rv.Name == name {
			for j, v := range rv.Version {
				if v.Name == version {
					runtimeVersions[i].Version[j].Image = image
					updated = true
					break
				}
			}
			if !updated {
				runtimeVersions[i].Version = append(runtimeVersions[i].Version, model.Version{
					Name:   version,
					Image:  image,
					Config: "{ports:\n  [ {containerPort: 22,\n  name: devbox-ssh-port,\n protocol: TCP}]\n  appPorts:\n  [{port: 8080,\n  name: devbox-app-port,\n  protocol: TCP}]\n  user: devbox,\n  workingDir: /home/devbox/project,\n  releaseCommand:\n    [/bin/bash\n , -c]\n  releaseArgs:\n   [/home/devbox/project/entrypoint.sh]}",
				})
				updated = true
			}
			break
		}
	}

	if !updated {
		log.Println("Creating new runtime")
		runtimeVersions = append(runtimeVersions, model.RuntimeVersion{
			Name: name,
			Version: []model.Version{
				{
					Name:   version,
					Image:  image,
					Config: "{ports:\n  [ {containerPort: 22,\n  name: devbox-ssh-port,\n protocol: TCP}]\n  appPorts:\n  [{port: 8080,\n  name: devbox-app-port,\n  protocol: TCP}]\n  user: devbox,\n  workingDir: /home/devbox/project,\n  releaseCommand:\n    [/bin/bash\n , -c]\n  releaseArgs:\n   [/home/devbox/project/entrypoint.sh]}",
				},
			},
		})
	}

	// Update the config object with the modified runtimeVersions
	switch kind {
	case "framework":
		config.Runtime.Framework = runtimeVersions
	case "language":
		config.Runtime.Language = runtimeVersions
	case "custom":
		config.Runtime.Custom = runtimeVersions
	case "os":
		config.Runtime.OS = runtimeVersions
	}

	saveConfig(config)
	return nil
}

func getRuntimeVersions(config *model.Config, kind string) []model.RuntimeVersion {
	switch kind {
	case "framework":
		return config.Runtime.Framework
	case "language":
		return config.Runtime.Language
	case "custom":
		return config.Runtime.Custom
	case "os":
		return config.Runtime.OS
	default:
		return nil
	}
}

func saveConfig(config *model.Config) {
	file, _ := json.MarshalIndent(config, "", "  ")
	_ = ioutil.WriteFile("config.json", file, 0644)
}
