package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"gopkg.in/yaml.v2"
	"io"
	"io/ioutil"
	"log"
	"math/rand"
	"net/http"
	"net/smtp"
	"os"
	"os/exec"
	"path/filepath"
	"syscall"
	"time"
	"unsafe"
)

const (
	CCHDEVICENAME                 = 32
	CCHFORMNAME                   = 32
	ENUM_CURRENT_SETTINGS  uint32 = 0xFFFFFFFF
	ENUM_REGISTRY_SETTINGS uint32 = 0xFFFFFFFE
	DISP_CHANGE_SUCCESSFUL uint32 = 0
	DISP_CHANGE_RESTART    uint32 = 1
	DISP_CHANGE_FAILED     uint32 = 0xFFFFFFFF
	DISP_CHANGE_BADMODE    uint32 = 0xFFFFFFFE
)

type DEVMODE struct {
	DmDeviceName       [CCHDEVICENAME]uint16
	DmSpecVersion      uint16
	DmDriverVersion    uint16
	DmSize             uint16
	DmDriverExtra      uint16
	DmFields           uint32
	DmOrientation      int16
	DmPaperSize        int16
	DmPaperLength      int16
	DmPaperWidth       int16
	DmScale            int16
	DmCopies           int16
	DmDefaultSource    int16
	DmPrintQuality     int16
	DmColor            int16
	DmDuplex           int16
	DmYResolution      int16
	DmTTOption         int16
	DmCollate          int16
	DmFormName         [CCHFORMNAME]uint16
	DmLogPixels        uint16
	DmBitsPerPel       uint32
	DmPelsWidth        uint32
	DmPelsHeight       uint32
	DmDisplayFlags     uint32
	DmDisplayFrequency uint32
	DmICMMethod        uint32
	DmICMIntent        uint32
	DmMediaType        uint32
	DmDitherType       uint32
	DmReserved1        uint32
	DmReserved2        uint32
	DmPanningWidth     uint32
	DmPanningHeight    uint32
}

type Config struct {
	Applications       Selection       `yaml:"applications"`
	Websites           Selection       `yaml:"websites"`
	Preferences        Preferences     `yaml:"preferences"`
	StartMenuItems     Selection       `yaml:"start_menu_items"`
	FileOperations     FileOperations  `yaml:"file_operations"`
	EmailOperations    EmailOperations `yaml:"email_operations"`
	SoftwareManagement Selection       `yaml:"software_management"`
	SystemUpdates      Selection       `yaml:"system_updates"`
	UserAccounts       []UserAccount   `yaml:"user_accounts"`
	NetworkSettings    []Network       `yaml:"network_settings"`
	SystemLogs         Selection       `yaml:"system_logs"`
	MediaFiles         MediaLocation   `yaml:"media_files"`
	Printing           Selection       `yaml:"printing"`
	ScheduledTasks     Selection       `yaml:"scheduled_tasks"`
	DecoyFiles         DecoyFiles      `yaml:"decoy_files"`
	InteractionDuration int            `yaml:"interaction_duration"`
	ActionDelay         int            `yaml:"action_delay"`
	RandomnessFactor    int            `yaml:"randomness_factor"`
	OpenAIAPIKey        string         `yaml:"openai_api_key"`
}

type Selection struct {
	Options         []string `yaml:"options"`
	SelectionMethod string   `yaml:"selection_method"`
}

type Preferences struct {
	DefaultBrowser    Selection     `yaml:"default_browser"`
	BackgroundImages  MediaLocation `yaml:"background_images"`
	ScreenResolutions Selection     `yaml:"screen_resolutions"`
	Languages         Selection     `yaml:"languages"`
}

type FileOperations struct {
	CreateModifyFiles []FileOperation `yaml:"create_modify_files"`
}

type FileOperation struct {
	Path      string `yaml:"path"`
	Content   string `yaml:"content"`
	UseGPT    bool   `yaml:"use_gpt"`
	GPTPrompt string `yaml:"gpt_prompt"`
}

type EmailOperations struct {
	GoogleAccount    EmailAccount     `yaml:"google_account"`
	MicrosoftAccount EmailAccount     `yaml:"microsoft_account"`
	SendReceive      []EmailOperation `yaml:"send_receive"`
}

type EmailAccount struct {
	Email    string `yaml:"email"`
	Password string `yaml:"password"`
}

type EmailOperation struct {
	SendTo    string `yaml:"send_to"`
	Subject   string `yaml:"subject"`
	Body      string `yaml:"body"`
	UseGPT    bool   `yaml:"use_gpt"`
	GPTPrompt string `yaml:"gpt_prompt"`
}

type MediaLocation struct {
	Location        string   `yaml:"location"`
	Type            string   `yaml:"type"`
	SelectionMethod string   `yaml:"selection_method"`
	Options         []string `yaml:"options"`
}

type DecoyFiles struct {
	Location        []string `yaml:"location"`
	Type            string   `yaml:"type"`
	TargetDirectory string   `yaml:"target_directory"`
}

type UserAccount struct {
	Name        string `yaml:"name"`
	Password    string `yaml:"password"`
	FullName    string `yaml:"full_name"`
	Description string `yaml:"description"`
}

type Network struct {
	SSID     string `yaml:"ssid"`
	Password string `yaml:"password"`
}

func loadConfig(filename string) (*Config, error) {
	data, err := ioutil.ReadFile(filename)
	if err != nil {
		return nil, err
	}
	var config Config
	err = yaml.Unmarshal(data, &config)
	if err != nil {
		return nil, err
	}
	return &config, nil
}

func selectRandomOrHardcoded(options []string, method string) []string {
	if method == "hardcoded" {
		return options
	}
	rand.Seed(time.Now().UnixNano())
	return []string{options[rand.Intn(len(options))]}
}

func downloadFile(url string, dest string) error {
	resp, err := http.Get(url)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	out, err := os.Create(dest)
	if err != nil {
		return err
	}
	defer out.Close()

	_, err = io.Copy(out, resp.Body)
	if err != nil {
		return err
	}
	return nil
}

func installApplications(config *Config) {
	apps := selectRandomOrHardcoded(config.Applications.Options, config.Applications.SelectionMethod)
	for _, app := range apps {
		cmd := exec.Command("choco", "install", app, "-y")
		err := cmd.Run()
		if err != nil {
			log.Printf("Failed to install %s: %v", app, err)
		} else {
			log.Printf("Installed %s", app)
		}
	}
}

func browseWebsites(config *Config) {
	websites := selectRandomOrHardcoded(config.Websites.Options, config.Websites.SelectionMethod)
	for _, website := range websites {
		cmd := exec.Command("cmd", "/C", "start", website)
		err := cmd.Run()
		if err != nil {
			log.Printf("Failed to browse %s: %v", website, err)
		} else {
			log.Printf("Browsing %s", website)
		}
		time.Sleep(time.Duration(config.ActionDelay+rand.Intn(2*config.RandomnessFactor+1)-config.RandomnessFactor) * time.Second)
	}
}

func changePreferences(config *Config) {
	prefs := config.Preferences
	browsers := selectRandomOrHardcoded(prefs.DefaultBrowser.Options, prefs.DefaultBrowser.SelectionMethod)
	for _, browser := range browsers {
		cmd := exec.Command("cmd", "/C", "start", browser)
		err := cmd.Run()
		if err != nil {
			log.Printf("Failed to set default browser: %v", err)
		}
	}

	images := selectRandomOrHardcoded(prefs.BackgroundImages.Options, prefs.BackgroundImages.SelectionMethod)
	for _, image := range images {
		err := downloadFile(image, "C:\\path\\to\\downloaded\\image.jpg")
		if err != nil {
			log.Printf("Failed to download background image: %v", err)
		}
		cmd := exec.Command("powershell", "-Command", "(New-Object -ComObject WScript.Shell).RegWrite('HKCU\\Control Panel\\Desktop\\Wallpaper', 'C:\\path\\to\\downloaded\\image.jpg')")
		err = cmd.Run()
		if err != nil {
			log.Printf("Failed to set background image: %v", err)
		}
	}

	resolutions := selectRandomOrHardcoded(prefs.ScreenResolutions.Options, prefs.ScreenResolutions.SelectionMethod)
	for _, resolution := range resolutions {
		setScreenResolution(resolution)
	}

	languages := selectRandomOrHardcoded(prefs.Languages.Options, prefs.Languages.SelectionMethod)
	for _, language := range languages {
		cmd := exec.Command("powershell", "-Command", fmt.Sprintf(`Set-WinUILanguageOverride -Language "%s"`, language))
		err := cmd.Run()
		if err != nil {
			log.Printf("Failed to set language to %s: %v", language, err)
		} else {
			log.Printf("Language set to %s", language)
		}
	}
}

func setScreenResolution(resolution string) {
	width, height := parseResolution(resolution)
	user32dll := syscall.NewLazyDLL("user32.dll")
	procEnumDisplaySettingsW := user32dll.NewProc("EnumDisplaySettingsW")
	procChangeDisplaySettingsW := user32dll.NewProc("ChangeDisplaySettingsW")

	devMode := new(DEVMODE)
	ret, _, _ := procEnumDisplaySettingsW.Call(uintptr(unsafe.Pointer(nil)),
		uintptr(ENUM_CURRENT_SETTINGS), uintptr(unsafe.Pointer(devMode)))

	if ret == 0 {
		log.Println("Couldn't extract display settings.")
		return
	}

	newMode := *devMode
	newMode.DmPelsWidth = uint32(width)
	newMode.DmPelsHeight = uint32(height)
	ret, _, _ = procChangeDisplaySettingsW.Call(uintptr(unsafe.Pointer(&newMode)), uintptr(0))

	switch ret {
	case uintptr(DISP_CHANGE_SUCCESSFUL):
		log.Println("Successfully changed the display resolution.")
	case uintptr(DISP_CHANGE_RESTART):
		log.Println("Restart required to apply the resolution changes.")
	case uintptr(DISP_CHANGE_BADMODE):
		log.Println("The resolution is not supported by the display.")
	case uintptr(DISP_CHANGE_FAILED):
		log.Println("Failed to change the display resolution.")
	}
}

func parseResolution(resolution string) (int, int) {
	var width, height int
	fmt.Sscanf(resolution, "%dx%d", &width, &height)
	return width, height
}

func addStartMenuItems(config *Config) {
	items := selectRandomOrHardcoded(config.StartMenuItems.Options, config.StartMenuItems.SelectionMethod)
	for _, item := range items {
		log.Printf("Adding start menu item %s", item)
	}
}

func createAndModifyFiles(config *Config) {
	for _, fileOp := range config.FileOperations.CreateModifyFiles {
		var content string
		if fileOp.UseGPT {
			content = generateContentUsingGPT(config.OpenAIAPIKey, fileOp.GPTPrompt)
		} else {
			content = fileOp.Content
		}
		err := ioutil.WriteFile(fileOp.Path, []byte(content), 0644)
		if err != nil {
			log.Printf("Failed to create/modify file %s: %v", fileOp.Path, err)
		} else {
			log.Printf("Created/modified file %s", fileOp.Path)
		}
	}
}

func generateContentUsingGPT(apiKey string, prompt string) string {
	prompt = fmt.Sprintf("This is where the prompt goes from the config file input: %s", prompt)

	log.Printf(prompt)

	// Prepare the JSON payload for the request
	payload := map[string]interface{}{
		"model":       "gpt-4o",
		"prompt":      prompt,
		"max_tokens":  2048,
		"temperature": 0.3,
	}
	payloadBytes, err := json.Marshal(payload)
	if err != nil {
		log.Fatalf("Error marshalling payload: %v", err)
	}

	// Create a new request
	req, err := http.NewRequest("POST", "https://api.openai.com/v1/completions", bytes.NewReader(payloadBytes))
	if err != nil {
		log.Fatalf("Error creating request: %v", err)
	}

	// Set headers
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+apiKey)

	// Execute the request
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		log.Fatalf("Error making request to OpenAI: %v", err)
	}
	defer resp.Body.Close()

	// Read the response body
	respBody, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Fatalf("Error reading response body: %v", err)
	}

	// Parse the response
	var response struct {
		Choices []struct {
			Text string `json:"text"`
		} `json:"choices"`
	}
	if err := json.Unmarshal(respBody, &response); err != nil {
		log.Fatalf("Error parsing response JSON: %v", err)
	}

	generatedContent := response.Choices[0].Text

	return generatedContent
}

func sendEmails(config *Config) {
	for _, emailOp := range config.EmailOperations.SendReceive {
		var body string
		if emailOp.UseGPT {
			body = generateContentUsingGPT(config.OpenAIAPIKey, emailOp.GPTPrompt)
		} else {
			body = emailOp.Body
		}
		sendEmail(config.EmailOperations.GoogleAccount, emailOp.SendTo, emailOp.Subject, body)
	}
}

func sendEmail(account EmailAccount, to, subject, body string) {
	from := account.Email
	password := account.Password

	msg := "From: " + from + "\n" +
		"To: " + to + "\n" +
		"Subject: " + subject + "\n\n" +
		body

	err := smtp.SendMail("smtp.gmail.com:587",
		smtp.PlainAuth("", from, password, "smtp.gmail.com"),
		from, []string{to}, []byte(msg))

	if err != nil {
		log.Printf("Failed to send email: %v", err)
	} else {
		log.Printf("Email sent to %s", to)
	}
}

func downloadDecoyFiles(config *Config) {
	urls := selectRandomOrHardcoded(config.DecoyFiles.Location, config.DecoyFiles.Type)
	for _, url := range urls {
		err := downloadFile(url, filepath.Join(config.DecoyFiles.TargetDirectory, filepath.Base(url)))
		if err != nil {
			log.Printf("Failed to download decoy file %s: %v", url, err)
		} else {
			log.Printf("Downloaded decoy file %s", url)
		}
	}
}

func manageSoftware(config *Config) {
	operations := selectRandomOrHardcoded(config.SoftwareManagement.Options, config.SoftwareManagement.SelectionMethod)
	for _, operation := range operations {
		cmd := exec.Command("choco", operation, "-y")
		err := cmd.Run()
		if err != nil {
			log.Printf("Failed to perform software management operation %s: %v", operation, err)
		} else {
			log.Printf("Performed software management operation %s", operation)
		}
	}
}

func performSystemUpdates(config *Config) {
	updates := config.SystemUpdates.Options
	for _, update := range updates {
		cmd := exec.Command("powershell", "-Command", update)
		err := cmd.Run()
		if err != nil {
			log.Printf("Failed to perform system update %s: %v", update, err)
		} else {
			log.Printf("Performed system update %s", update)
		}
	}
}

func manageUserAccounts(config *Config) {
	for _, account := range config.UserAccounts {
		cmd := exec.Command("powershell", "-Command", fmt.Sprintf(`$Password = "%s" | ConvertTo-SecureString -AsPlainText -Force; $params = @{Name = "%s"; Password = $Password; FullName = "%s"; Description = "%s"}; New-LocalUser @params`, account.Password, account.Name, account.FullName, account.Description))
		err := cmd.Run()
		if err != nil {
			log.Printf("Failed to manage user account %s: %v", account.Name, err)
		} else {
			log.Printf("Managed user account %s", account.Name)
		}
	}
}

func manageNetworkSettings(config *Config) {
	for _, network := range config.NetworkSettings {
		xmlContent := fmt.Sprintf(`
		<?xml version="1.0"?>
		<WLANProfile xmlns="http://www.microsoft.com/networking/WLAN/profile/v1">
		    <name>%s</name>
		    <SSIDConfig>
		        <SSID>
		            <name>%s</name>
		        </SSID>
		    </SSIDConfig>
		    <connectionType>ESS</connectionType>
		    <connectionMode>auto</connectionMode>
		    <MSM>
		        <security>
		            <authEncryption>
		                <authentication>WPA2PSK</authentication>
		                <encryption>AES</encryption>
		                <useOneX>false</useOneX>
		            </authEncryption>
		            <sharedKey>
		                <keyType>passPhrase</keyType>
		                <protected>false</protected>
		                <keyMaterial>%s</keyMaterial>
		            </sharedKey>
		        </security>
		    </MSM>
		    <MacRandomization xmlns="http://www.microsoft.com/networking/WLAN/profile/v3">
		        <enableRandomization>false</enableRandomization>
		    </MacRandomization>
		</WLANProfile>`, network.SSID, network.SSID, network.Password)

		xmlFilePath := fmt.Sprintf("%s_profile.xml", network.SSID)
		err := ioutil.WriteFile(xmlFilePath, []byte(xmlContent), 0644)
		if err != nil {
			log.Printf("Failed to write XML file for network %s: %v", network.SSID, err)
			continue
		}

		cmd := exec.Command("netsh", "wlan", "add", "profile", fmt.Sprintf("filename=\"%s\"", xmlFilePath))
		err = cmd.Run()
		if err != nil {
			log.Printf("Failed to add profile for network %s: %v", network.SSID, err)
			continue
		}

		cmd = exec.Command("netsh", "wlan", "connect", fmt.Sprintf("name=\"%s\"", network.SSID))
		err = cmd.Run()
		if err != nil {
			log.Printf("Failed to connect to network %s: %v", network.SSID, err)
		} else {
			log.Printf("Connected to network %s", network.SSID)
		}
	}
}

func openMediaFiles(config *Config) {
	mediaFiles := selectRandomOrHardcoded(config.MediaFiles.Options, config.MediaFiles.SelectionMethod)
	for _, mediaFile := range mediaFiles {
		cmd := exec.Command("cmd", "/C", "start", mediaFile)
		err := cmd.Run()
		if err != nil {
			log.Printf("Failed to open media file %s: %v", mediaFile, err)
		} else {
			log.Printf("Opened media file %s", mediaFile)
		}
	}
}

func printDocuments(config *Config) {
	documents := selectRandomOrHardcoded(config.Printing.Options, config.Printing.SelectionMethod)
	for _, document := range documents {
		cmd := exec.Command("powershell", "-Command", "Start-Process", "notepad.exe", document, "-Verb", "Print")
		err := cmd.Run()
		if err != nil {
			log.Printf("Failed to print document %s: %v", document, err)
		} else {
			log.Printf("Printed document %s", document)
		}
	}
}

func createScheduledTasks(config *Config) {
	tasks := selectRandomOrHardcoded(config.ScheduledTasks.Options, config.ScheduledTasks.SelectionMethod)
	for _, task := range tasks {
		cmd := exec.Command("powershell", "-Command", "schtasks", "/create", task)
		err := cmd.Run()
		if err != nil {
			log.Printf("Failed to create scheduled task %s: %v", task, err)
		} else {
			log.Printf("Created scheduled task %s", task)
		}
	}
}

func main() {
	config, err := loadConfig("config.yaml")
	if err != nil {
		log.Fatalf("Error loading config: %v", err)
	}

	time.Sleep(time.Duration(config.InteractionDuration) * time.Second)

	delayWithRandomness := func(delay, randomnessFactor int) {
		randomDelay := delay + rand.Intn(2*randomnessFactor+1) - randomnessFactor
		time.Sleep(time.Duration(randomDelay) * time.Second)
	}

	interactions := []func(*Config){
		installApplications,
		browseWebsites,
		changePreferences,
		addStartMenuItems,
		createAndModifyFiles,
		sendEmails,
		downloadDecoyFiles,
		manageSoftware,
		performSystemUpdates,
		manageUserAccounts,
		manageNetworkSettings,
		openMediaFiles,
		printDocuments,
		createScheduledTasks,
	}

	for _, interaction := range interactions {
		interaction(config)
		delayWithRandomness(config.ActionDelay, config.RandomnessFactor)
	}

	log.Println("All interactions completed.")
}
