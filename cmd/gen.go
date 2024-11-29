package cmd

import (
	"encoding/json"
	"fmt"
	"github.com/spf13/cobra"
	"io/ioutil"
	"log"
	"runtimectl/model"
	"runtimectl/pkg/util"
)

func newGenCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "gen",
		Short: "gen the json in config.json ",
		RunE:  genAction,
	}

	cmd.Flags().StringVar(&kind, "kind", "", "Kind of runtime (Framework, Language, Custom, OS)")
	cmd.Flags().StringVar(&name, "name", "", "Name of the runtime")
	cmd.Flags().StringVar(&version, "version", "", "Version of the runtime")
	cmd.Flags().StringVar(&image, "image", "", "Image of the runtime")
	cmd.Flags().StringVar(&path, "path", "config.json", "Path to the config file")
	cmd.Flags().BoolVar(&active, "active", true, "Force update the runtime")

	return cmd
}

func genAction(cmd *cobra.Command, args []string) error {
	return gen(kind, name, version, image, path, active)
}

func gen(kind, name, version, image, path string, active bool) error {
	config := util.ParseJson(path)
	if config == nil {
		return fmt.Errorf("failed to parse config.json")
	}

	runtimeVersions := getRuntimeVersions(config, kind)
	if runtimeVersions == nil {
		return fmt.Errorf("invalid kind: %s", kind)
	}

	updated := false
	state := "deprecated"
	if active {
		state = "active"
	}

	for i, rv := range runtimeVersions {
		if rv.Name == name {
			for j, v := range rv.Version {
				if v.Name == version {
					if active {
						if v.Image != image {
							runtimeVersions[i].Version[j].State = "deprecated"
						} else {
							updated = true
							runtimeVersions[i].Version[j].State = state
						}
					} else {
						if v.Image == image {
							runtimeVersions[i].Version[j].State = state
							updated = true
						}
					}
				}
			}

			if !updated {
				runtimeVersions[i].Version = append(runtimeVersions[i].Version, model.Version{
					Name:   version,
					Image:  image,
					Config: formatConfig(),
					State:  state,
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
					Config: formatConfig(),
					State:  state,
				},
			},
		})
	}

	// Update the config object with the modified runtimeVersions
	switch kind {
	case "Framework":
		config.Runtime.Framework = runtimeVersions
	case "Language":
		config.Runtime.Language = runtimeVersions
	case "Custom":
		config.Runtime.Custom = runtimeVersions
	case "OS":
		config.Runtime.OS = runtimeVersions
	}

	saveConfig(path, config)
	return nil
}

func formatConfig() string {
	config := map[string]interface{}{
		"ports": []map[string]interface{}{
			{"containerPort": 22, "name": "devbox-ssh-port", "protocol": "TCP"},
		},
		"appPorts": []map[string]interface{}{
			{"port": 8080, "name": "devbox-app-port", "protocol": "TCP"},
		},
		"user":           "devbox",
		"workingDir":     "/home/devbox/project",
		"releaseCommand": []string{"/bin/bash", "-c"},
		"releaseArgs":    []string{"/home/devbox/project/entrypoint.sh"},
	}

	configBytes, _ := json.Marshal(config)
	return string(configBytes)
}

func getRuntimeVersions(config *model.Config, kind string) []model.RuntimeVersion {
	switch kind {
	case "Framework":
		return config.Runtime.Framework
	case "Language":
		return config.Runtime.Language
	case "Custom":
		return config.Runtime.Custom
	case "OS":
		return config.Runtime.OS
	default:
		return nil
	}
}

func saveConfig(path string, config *model.Config) {
	file, _ := json.MarshalIndent(config, "", "  ")
	_ = ioutil.WriteFile(path, file, 0644)
}
