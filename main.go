package main

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"time"

	"github.com/spf13/cobra"
)

const basePath = `\\192.168.9.10\Case Archive\Case-Forensic`

func openFolder(folderPath string) error {
	var cmd *exec.Cmd

	switch runtime.GOOS {
	case "windows":
		cmd = exec.Command("explorer", folderPath)
	case "darwin":
		cmd = exec.Command("open", folderPath)
	case "linux":
		cmd = exec.Command("xdg-open", folderPath)
	default:
		return fmt.Errorf("unsupported platform")
	}

	return cmd.Start()
}

func checkFolder(folderPath string) error {
	info, err := os.Stat(folderPath)
	if os.IsNotExist(err) {
		return fmt.Errorf("folder does not exist")
	}
	if err != nil {
		return fmt.Errorf("error accessing folder: %v", err)
	}
	if !info.IsDir() {
		return fmt.Errorf("the target path is not a folder")
	}
	return nil
}

func main() {
	var year string

	var rootCmd = &cobra.Command{
		Use:   "open-case [caseId]",
		Short: "Open a case folder in the file explorer based on year and case",
		Args:  cobra.MinimumNArgs(1),
		Run: func(cmd *cobra.Command, args []string) {
			caseId := fmt.Sprintf("%03s", args[0])

			targetPath := filepath.Join(basePath, year)
			if caseId != "dir" {
				targetPath = filepath.Join(targetPath, fmt.Sprintf("F-%s-%s", year, caseId))
			}

			err := checkFolder(targetPath)
			if err != nil {
				fmt.Printf("Error: %v\n", err)
				return
			}

			err = openFolder(targetPath)
			if err != nil {
				fmt.Println("Error:", err)
			} else {
				fmt.Println("Folder opened successfully:", targetPath)
			}
		},
	}

	rootCmd.Flags().StringVarP(&year, "year", "y", fmt.Sprintf("%d", time.Now().Year()), "Year of the case (required)")

	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
	}
}
