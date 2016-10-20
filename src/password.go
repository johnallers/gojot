package sdees

import (
	"fmt"
	"os"
	"strings"
	"syscall"
	"time"

	"golang.org/x/crypto/ssh/terminal"
)

// PromptPassword prompts for password and tests against the file in input,
// use "" for no file, in which a new password will be generated
func PromptPassword(gitfolder string) string {
	password1 := "1"
	textToTest, _ := GetTextOfOne(gitfolder, "master", ".key.gpg")
	if len(textToTest) == 0 {
		fmt.Printf("Getting new password\n")
		password2 := "2"
		for password1 != password2 {
			fmt.Printf("Enter new password for: ")
			bytePassword, _ := terminal.ReadPassword(int(os.Stdin.Fd()))
			password1 = strings.TrimSpace(string(bytePassword))
			fmt.Printf("\nEnter password again: ")
			bytePassword2, _ := terminal.ReadPassword(int(syscall.Stdin))
			password2 = strings.TrimSpace(string(bytePassword2))
			if password1 != password2 {
				fmt.Println("\nPasswords do not match.")
			}
		}
		Passphrase = password1

		logger.Debug("It seems key doesn't exist yet, making it")
		WriteToMaster(gitfolder, ".key", RandStringBytesMaskImprSrc(100, time.Now().UnixNano()))

	} else {
		logger.Debug("Testing with master:key.gpg")
		passwordAccepted := false
		for passwordAccepted == false {
			fmt.Printf("\nEnter password: ")
			bytePassword, _ := terminal.ReadPassword(int(os.Stdin.Fd()))
			password1 = strings.TrimSpace(string(bytePassword))
			_, err := DecryptString(textToTest, password1)
			if err == nil {
				passwordAccepted = true
			} else {
				fmt.Println("\nPasswords do not match.")
				logger.Debug("Got error: %s", err.Error())
			}
		}
	}
	fmt.Println("")
	return password1
}
