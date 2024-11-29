package cmd

import (
	"github.com/spf13/cobra"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	client "runtimectl/pkg/k8s"
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
		active, _, _ := unstructured.NestedString(r.Object, "spec", "state")
		class, _, _ := unstructured.NestedString(r.Object, "spec", "classRef")
		runtimeClass, err := k8sClient.GetRuntimeClass(class)
		kind, _, _ := unstructured.NestedString(runtimeClass.Object, "spec", "kind")
		version, _, _ := unstructured.NestedString(r.Object, "spec", "version")
		image, _, _ := unstructured.NestedString(r.Object, "spec", "config", "image")
		state := "active" == active
		if err = gen(kind, class, version, image, outputFile, state); err != nil {
			return err
		}
	}
	return nil
}
