package cmd

import (
	"fmt"
	"github.com/spf13/cobra"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"log"
	"runtimectl/model"
	client "runtimectl/pkg/k8s"
	"runtimectl/pkg/util"
)

func newExportCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "export",
		Short: "export the runtimes to json",
		RunE:  exportAction,
	}
	cmd.Flags().StringVar(&outputFile, "output", "output.json", "path to the output file")
	cmd.Flags().StringVar(&k8sConfig, "k8s-config", "/root/.kube/config", "the k8s config file path")
	return cmd
}

func exportAction(cmd *cobra.Command, args []string) error {
	k8sClient := client.Init(k8sConfig)
	runtimeList, err := k8sClient.GetAllRuntime()
	if err != nil {
		return err
	}
	for _, r := range runtimeList.Items {
		class, _, _ := unstructured.NestedString(r.Object, "spec", "classRef")
		runtimeClass, err := k8sClient.GetRuntimeClass(class)
		kind, _, _ := unstructured.NestedString(runtimeClass.Object, "spec", "kind")
		version, _, _ := unstructured.NestedString(r.Object, "spec", "version")
		image, _, _ := unstructured.NestedString(r.Object, "spec", "config", "image")
		c, _ := k8sClient.GetRuntimeConfig(r)
		state, _, _ := unstructured.NestedString(r.Object, "spec", "state")
		if err = export(kind, class, version, image, outputFile, c, state); err != nil {
			return err
		}
	}
	return nil
}

func export(kind, name, version, image, path, c, state string) error {
	config := util.ParseJson(path)
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
			runtimeVersions[i].Version = append(runtimeVersions[i].Version, model.Version{
				Name:   version,
				Image:  image,
				Config: c,
				State:  state,
			})
			updated = true
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
					Config: c,
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
