package cmd

import (
	"fmt"
	"os/exec"

	"github.com/spf13/cobra"
)

var gpgCmd = "gpg"
var zipCmd = "zip"
var rmCmd = "rm"
var defaultUser = "testUser"

var backupCmd = &cobra.Command{
	Use:   "backup",
	Short: "Run a backup process",
	Long: `Run a backup process`,
	Run: func(cmd *cobra.Command, args []string) {
		checkForGpgKeys()

		filesToBackup := args

		backups := backup(filesToBackup)
		encrypt(backups)
		remove(backups)

		fmt.Println(backups)
	},
}

func checkForGpgKeys() {
	checkKeysCmd := exec.Command(gpgCmd, "--list-keys", defaultUser)
	stdout, err := checkKeysCmd.Output()

	if (string(stdout) == "") {
		genKeyCmd := exec.Command(gpgCmd, "--quick-generate-key", defaultUser)
		
		_, err = genKeyCmd.Output()

		if err != nil {
			fmt.Println(err.Error())
			return
		}
	}
}

func backup(files []string) []string {
	archivedFiles := []string {}

	for _, file := range files {
		archiveName := file + ".zip"

		archiveCmd := exec.Command(zipCmd, "-r", archiveName, file)
		_, err := archiveCmd.Output()

		if err != nil {
			fmt.Println(err.Error())
		} else {
			archivedFiles = append(archivedFiles, archiveName)
		}
	}

	return archivedFiles
}

func encrypt(files []string) []string {
	encryptedFiles := []string {}

	for _, file := range files {
		encryptedFile := file + ".gpg"

		encryptCmd := exec.Command(gpgCmd, "-r", defaultUser, "-e", file)
		_, err := encryptCmd.Output()

		if err != nil {
			fmt.Println(err.Error())
		} else {
			encryptedFiles = append(encryptedFiles, encryptedFile)
		}
	}

	return encryptedFiles
}

func remove(files []string) {
	for _, file := range files {
		removeCmd := exec.Command(rmCmd, file)
		_, err := removeCmd.Output()

		if err != nil {
			fmt.Println(err.Error())
		}
	}
}

func init() {
	rootCmd.AddCommand(backupCmd)

	// backupCmd.PersistentFlags().String("foo", "", "A help for foo")

	// backupCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
