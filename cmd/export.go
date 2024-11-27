package cmd

import (
	"encoding/json"
	"github.com/spf13/cobra"
	"os"
	"runtimectl/config"
	"runtimectl/dao"
)

func newExportCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "export",
		Short: "export the template uuid to the json file",
		RunE:  exportAction,
	}

	cmd.Flags().StringVar(&outputFile, "output", "output.json", "Path to the output file")

	return cmd
}

func exportAction(cmd *cobra.Command, args []string) error {
	config.Init()
	dao.Init()
	file, err := os.Create(outputFile)
	if err != nil {
		return err
	}
	defer file.Close()

	data, err := dao.GetTemplates()
	if err != nil {
		return err
	}
	encoder := json.NewEncoder(file)
	encoder.SetIndent("", "  ")
	if err := encoder.Encode(data); err != nil {
		return err
	}
	return nil
}
