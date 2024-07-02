package main

import (
	"io/ioutil"

	"gopkg.in/yaml.v2"
)

type Config struct {
	Applications       StringSelection        `yaml:"applications"`
	Websites           StringSelection        `yaml:"websites"`
	Preferences        Preferences            `yaml:"preferences"`
	StartMenuItems     StringSelection        `yaml:"start_menu_items"`
	FileOperations     FileOperations         `yaml:"file_operations"`
	EmailOperations    EmailOperations        `yaml:"email_operations"`
	SoftwareManagement StringSelection        `yaml:"software_management"`
	SystemUpdates      SystemUpdates          `yaml:"system_updates"`
	UserAccounts       []UserAccount          `yaml:"user_accounts"`
	NetworkSettings    []Network              `yaml:"network_settings"`
	SystemLogs         StringSelection        `yaml:"system_logs"`
	MediaFiles         MediaLocation          `yaml:"media_files"`
	Printing           StringSelection        `yaml:"printing"`
	ScheduledTasks     ScheduledTaskSelection `yaml:"scheduled_tasks"`
	DecoyFiles         DecoyFiles             `yaml:"decoy_files"`
	Lures              []LureConfig           `yaml:"lures"`
	General            GeneralConfig          `yaml:"general"`
}

type StringSelection struct {
	Options         []string `yaml:"options"`
	SelectionMethod string   `yaml:"selection_method"`
}

type ScheduledTaskSelection struct {
	Options         []ScheduledTask `yaml:"options"`
	SelectionMethod string          `yaml:"selection_method"`
}

type Preferences struct {
	DefaultBrowser    StringSelection `yaml:"default_browser"`
	BackgroundImages  MediaLocation   `yaml:"background_images"`
	ScreenResolutions StringSelection `yaml:"screen_resolutions"`
	Languages         StringSelection `yaml:"languages"`
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
	Sets []DecoyFileSet `yaml:"sets"`
}

type DecoyFileSet struct {
	Location        []string `yaml:"location"`
	Type            string   `yaml:"type"`
	TargetDirectory []string `yaml:"target_directory"`
	SelectionMethod string   `yaml:"selection_method"`
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

type ScheduledTask struct {
	Name      string `yaml:"name"`
	Path      string `yaml:"path"`
	Schedule  string `yaml:"schedule"`
	StartTime string `yaml:"start_time"`
}

type LureConfig struct {
	Name             string                 `yaml:"name"`
	Type             string                 `yaml:"type"`
	Location         string                 `yaml:"location"`
	GenerationParams map[string]interface{} `yaml:"generation_params"`
	GenerativeType   string                 `yaml:"generative_type"`
	OpenaiPrompt     string                 `yaml:"openai_prompt,omitempty"`
	Responder        string                 `yaml:"responder,omitempty"`
}

type SystemUpdates struct {
	Method          string   `yaml:"method"`
	SpecificUpdates []string `yaml:"specific_updates,omitempty"`
	SelectionMethod string   `yaml:"selection_method"`
	HideUpdates     []string `yaml:"hide_updates,omitempty"`
}

type GeneralConfig struct {
	Redis struct {
		IP   string `yaml:"ip"`
		Port int    `yaml:"port"`
	} `yaml:"redis"`
	LogFile            string   `yaml:"log_file,omitempty"`
	OpenaiApiKey       string   `yaml:"openai_api_key,omitempty"`
	InteractionDuration int      `yaml:"interaction_duration"`
	ActionDelay         int      `yaml:"action_delay"`
	RandomnessFactor    int      `yaml:"randomness_factor"`
	Usernames           []string `yaml:"usernames"`
	SelectionMethod     string   `yaml:"selection_method"`
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
