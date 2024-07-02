package main

import (
	"fmt"
	"io"
	"math/rand"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"time"
	"net/url"

	"github.com/pkg/sftp"
	"golang.org/x/crypto/ssh"
)

func selectRandomOrHardcoded(options []string, method string) []string {
	if len(options) == 0 {
		logToFile("config.General.LogFile", "No options provided for selection.")
		return options
	}
	if method == "hardcoded" {
		return options
	}
	rand.Seed(time.Now().UnixNano())
	return []string{options[rand.Intn(len(options))]}
}

func downloadFile(url, dest string) error {
	resp, err := http.Get(url)
	if err != nil {
		logToFile("config.General.LogFile", fmt.Sprintf("Failed to download HTTP file %s: %v", url, err))
		return err
	}
	defer resp.Body.Close()

	out, err := os.Create(dest)
	if err != nil {
		logToFile("config.General.LogFile", fmt.Sprintf("Failed to create file %s: %v", dest, err))
		return err
	}
	defer out.Close()

	_, err = io.Copy(out, resp.Body)
	if err != nil {
		logToFile("config.General.LogFile", fmt.Sprintf("Failed to copy HTTP file %s: %v", url, err))
	}
	return err
}

func downloadSMBFile(smbPath, dest string) error {
	cmd := exec.Command("powershell", "-Command", fmt.Sprintf("New-PSDrive -Name S -PSProvider FileSystem -Root \\\\%s; Copy-Item S:\\%s -Destination %s; Remove-PSDrive -Name S", smbPath, filepath.Base(dest), dest))
	output, err := cmd.CombinedOutput()
	if err != nil {
		logToFile("config.General.LogFile", fmt.Sprintf("Failed to download SMB file %s: %v, output: %s", smbPath, err, string(output)))
	}
	return err
}

func downloadNFSFile(nfsPath, dest string) error {
	cmd := exec.Command("powershell", "-Command", fmt.Sprintf("New-PSDrive -Name N -PSProvider FileSystem -Root %s; Copy-Item N:\\%s -Destination %s; Remove-PSDrive -Name N", nfsPath, filepath.Base(dest), dest))
	output, err := cmd.CombinedOutput()
	if err != nil {
		logToFile("config.General.LogFile", fmt.Sprintf("Failed to download NFS file %s: %v, output: %s", nfsPath, err, string(output)))
	}
	return err
}

func downloadSFTPFile(sftpPath, dest, username, password string) error {
	config := &ssh.ClientConfig{
		User: username,
		Auth: []ssh.AuthMethod{
			ssh.Password(password),
		},
		HostKeyCallback: ssh.InsecureIgnoreHostKey(),
	}
	client, err := ssh.Dial("tcp", sftpPath, config)
	if err != nil {
		logToFile("config.General.LogFile", fmt.Sprintf("Failed to connect to SFTP server %s: %v", sftpPath, err))
		return err
	}
	defer client.Close()

	sftpClient, err := sftp.NewClient(client)
	if err != nil {
		logToFile("config.General.LogFile", fmt.Sprintf("Failed to create SFTP client %s: %v", sftpPath, err))
		return err
	}
	defer sftpClient.Close()

	source, err := sftpClient.Open(filepath.Base(sftpPath))
	if err != nil {
		logToFile("config.General.LogFile", fmt.Sprintf("Failed to open source SFTP file %s: %v", sftpPath, err))
		return err
	}
	defer source.Close()

	output, err := os.Create(dest)
	if err != nil {
		logToFile("config.General.LogFile", fmt.Sprintf("Failed to create destination file %s: %v", dest, err))
		return err
	}
	defer output.Close()

	_, err = io.Copy(output, source)
	if err != nil {
		logToFile("config.General.LogFile", fmt.Sprintf("Failed to copy SFTP file %s: %v", sftpPath, err))
	}
	return err
}

func downloadDecoyFiles(config *Config) {
	for _, set := range config.DecoyFiles.Sets {
		locations := selectRandomOrHardcoded(set.Location, set.SelectionMethod)
		targetDirs := selectRandomOrHardcoded(set.TargetDirectory, set.SelectionMethod)

		for _, location := range locations {
			for _, targetDir := range targetDirs {
				filename := filepath.Base(location)
				targetPath := filepath.Join(targetDir, filename)
				var err error
				switch {
				case strings.HasPrefix(location, "http"):
					err = downloadFile(location, targetPath)
				case strings.HasPrefix(location, "smb"):
					err = downloadSMBFile(location, targetPath)
				case strings.HasPrefix(location, "nfs"):
					err = downloadNFSFile(location, targetPath)
				case strings.HasPrefix(location, "sftp"):
					// Assume SFTP URL format: sftp://username:password@host:port/path
					url, parseErr := url.Parse(location)
					if parseErr != nil {
						logToFile("config.General.LogFile", fmt.Sprintf("Failed to parse SFTP URL %s: %v", location, parseErr))
						continue
					}
					username := url.User.Username()
					password, _ := url.User.Password()
					err = downloadSFTPFile(url.Host, targetPath, username, password)
				default:
					logToFile("config.General.LogFile", fmt.Sprintf("Unsupported protocol for decoy file %s", location))
				}
				if err != nil {
					logToFile("config.General.LogFile", fmt.Sprintf("Failed to download decoy file %s: %v", location, err))
				} else {
					logToFile("config.General.LogFile", fmt.Sprintf("Downloaded decoy file %s", location))
				}
			}
		}
	}
}
