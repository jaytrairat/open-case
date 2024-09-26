package main

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"

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
	var caseID int

	var rootCmd = &cobra.Command{
		Use:   "openfolder",
		Short: "Open a case folder in the file explorer based on year and case",
		Run: func(cmd *cobra.Command, args []string) {
			if year == "" || caseID == 0 {
				fmt.Println("Error: Both year (-y) and case (-c) must be provided.")
				return
			}

			caseIDStr := fmt.Sprintf("%03d", caseID)

			targetPath := filepath.Join(basePath, year, fmt.Sprintf("F-%s-%s", year, caseIDStr))

			err := checkFolder(targetPath)
			fmt.Println(targetPath)
			if err != nil {
				fmt.Printf("Error: %v\n", err)
				return
			}

			// Open the folder
			err = openFolder(targetPath)
			if err != nil {
				fmt.Println("Error:", err)
			} else {
				fmt.Println("Folder opened successfully:", targetPath)
			}
		},
	}

	rootCmd.Flags().StringVarP(&year, "year", "y", "", "Year of the case (required)")
	rootCmd.Flags().IntVarP(&caseID, "case", "c", 0, "Case ID (3-digit number) (required)")

	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
	}
}
