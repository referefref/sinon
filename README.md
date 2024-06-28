# Sinon - Modular Windows Burn-In Automation with Generative AI for Deception

![image](https://github.com/referefref/sinon/assets/56499429/9c5395ac-2034-4bcb-8214-c3690013d521)

Sinon is a modular tool for automatic burn-in of Windows-based deception hosts that aims to reduce the difficulty of orchestrating deception hosts at scale whilst enabling diversity and randomness through generative capabilities. It has been created as a proof-of-concept and is not intended for production deception environments. It would likely be better suited to having content pre-generated and built into a one-time script, as we wouldn't want to be storing secrets like OpenAI API keys on a decoy or deception host.

## Features

- Generative content including files, emails, and so on using OpenAI API
- Randomness factor - select from list in config, or follow config completely
- Temporal randomness - set delay to execution and delay between events including randomness factor

Sinon performs the following functions, as determined by a config file:

- **Install Applications:** Automatically installs specified applications using Chocolatey.
- **Browse Websites:** Opens and browses a list of websites.
- **Change Preferences:** Sets default browsers, background images, screen resolutions, and languages.
- **Add Start Menu Items:** Adds specified items to the start menu.
- **Create and Modify Files:** Generates and modifies files, optionally using GPT-4 for content generation.
- **Send Emails:** Sends emails using specified Google or Microsoft accounts.
- **Download Decoy Files:** Downloads decoy files from specified sources.
- **Manage Software:** Performs software management tasks such as installation and uninstallation.
- **Perform System Updates:** Executes system update commands.
- **Manage User Accounts:** Creates and manages user accounts.
- **Manage Network Settings:** Configures Wi-Fi networks.
- **Open Media Files:** Opens media files (images, video, and audio).
- **Print Documents:** Sends documents to print.
- **Create Scheduled Tasks:** Creates and manages scheduled tasks.

## Usage

1. **Clone the repository:**
    ```sh
    git clone https://github.com/yourusername/sinon.git
    cd sinon
    ```

2. **Configure the application:**
   - Modify the `config.yaml` file to suit your needs. See the [Config Items](#config-items) section for details.

3. **Build and run the application:**
    ```sh
    go build -o sinon
    ./sinon
    ```

## Config Items

The `config.yaml` file contains all the configuration options for Sinon. Here is an example configuration file with explanations:

```yaml
applications:
  options:
    - googlechrome
    - vlc
    - 7zip
    - notepadplusplus
    - git
    - firefox
    - winscp
    - slack
  selection_method: random

websites:
  options:
    - https://www.example.com
    - https://www.google.com
    - https://www.github.com
    - https://www.stackoverflow.com
    - https://www.reddit.com
    - https://www.wikipedia.org
    - https://www.medium.com
    - https://news.ycombinator.com
  selection_method: random

preferences:
  default_browser:
    options: [googlechrome, firefox, edge]
    selection_method: random
  background_images:
    location: "http://example.com/backgrounds"
    type: "http"
    selection_method: random
    options:
      - http://example.com/backgrounds/image1.jpg
      - http://example.com/backgrounds/image2.jpg
      - http://example.com/backgrounds/image3.jpg
  screen_resolutions:
    options: ["1920x1080", "1280x720", "1366x768", "1440x900"]
    selection_method: random
  languages:
    options: ["en-US", "es-ES", "fr-FR", "de-DE", "zh-CN"]
    selection_method: random

start_menu_items:
  options:
    - name: Google Chrome
      path: "C:\\Program Files\\Google\\Chrome\\Application\\chrome.exe"
    - name: VLC Media Player
      path: "C:\\Program Files\\VideoLAN\\VLC\\vlc.exe"
    - name: Notepad++
      path: "C:\\Program Files\\Notepad++\\notepad++.exe"
    - name: Git Bash
      path: "C:\\Program Files\\Git\\git-bash.exe"
  selection_method: random

file_operations:
  create_modify_files:
    options:
      - path: "C:\\Users\\Public\\Documents\\report.txt"
        content: "This is a test report file."
        use_gpt: false
      - path: "C:\\Users\\Public\\Documents\\meeting_notes.txt"
        content: "Generate a summary of the last meeting."
        use_gpt: true
        gpt_prompt: "Generate a detailed summary of the last team meeting discussing project milestones and deadlines."
      - path: "C:\\Users\\Public\\Documents\\financial_analysis.txt"
        content: "Analyze the financial data for Q1."
        use_gpt: true
        gpt_prompt: "Analyze the financial data for Q1, focusing on revenue, expenses, and profit margins."
    selection_method: random

email_operations:
  google_account:
    email: "user@gmail.com"
    password: "password"
  microsoft_account:
    email: "user@outlook.com"
    password: "password"
  send_receive:
    options:
      - send_to: "colleague@example.com"
        subject: "Project Update"
        body: "Please find the latest update on the project attached."
        use_gpt: false
      - send_to: "manager@example.com"
        subject: "Weekly Report"
        body: "Generate a weekly report on the team's performance."
        use_gpt: true
        gpt_prompt: "Generate a weekly report on the team's performance, including completed tasks, ongoing projects, and any blockers."
    selection_method: random

software_management:
  options:
    - install: "vlc"
    - uninstall: "notepadplusplus"
    - install: "git"
    - uninstall: "slack"
  selection_method: random

system_updates:
  options: ["Get-WindowsUpdate -Install", "Install-WindowsUpdate -AcceptAll"]
  selection_method: hardcoded

user_accounts:
  options:
    - name: "User01"
      password: "Password123!"
      full_name: "First User"
      description: "First test user account."
    - name: "User02"
      password: "Password456!"
      full_name: "Second User"
      description: "Second test user account."
    - name: "User03"
      password: "Password789!"
      full_name: "Third User"
      description: "Third test user account."
  selection_method: random

network_settings:
  options:
    - ssid: "HomeNetwork"
      password: "HomePassword123"
    - ssid: "OfficeNetwork"
      password: "OfficePassword456"
    - ssid: "GuestNetwork"
      password: "GuestPassword789"
  selection_method: random

system_logs:
  options: ["Get-EventLog -LogName System", "Get-EventLog -LogName Application"]
  selection_method: hardcoded

media_files:
  location: "http://example.com/media_files"
  type: "http"
  selection_method: random
  options:
    - http://example.com/media_files/sample_video.mp4
    - http://example.com/media_files/sample_music.mp3
    - http://example.com/media_files/sample_image.jpg

printing:
  options:
    - "C:\\Users\\Public\\Documents\\document1.pdf"
    - "C:\\Users\\Public\\Documents\\document2.txt"
  selection_method: random

scheduled_tasks:
  options:
    - name: "ExampleTask"
      path: "C:\\Path\\To\\Executable.exe"
      schedule: "daily"
      start_time: "12:00"
    - name: "ExampleTask2"
      path: "C:\\Path\\To\\AnotherExecutable.exe"
      schedule: "weekly"
      start_time: "08:00"
  selection_method: random

decoy_files:
  location: ["http://jamesbrine.com.au/passwords.xlsx"]
  type: "http"
  target_directory: "C:\\Users\\Public\\DecoyFiles"

interaction_duration: 3600
action_delay: 5
randomness_factor: 2

openai_api_key: ""
```

## Deploying Windows Deception Hosts
Sinon is designed to automate the setup of deception hosts by performing a variety of actions that simulate real user activity. The goal is to create a realistic environment that can deceive potential intruders. The modular and configurable nature of Sinon allows for easy adjustments and randomization, making each deployment unique.

## Steps to Deploy
- Prepare the Windows environment: Ensure that the target Windows machine is ready and accessible.
- Configure Sinon: Edit the config.yaml file to define the desired behaviors and settings.
- Run Sinon: Execute the compiled Sinon binary to start the automation process.
- Monitor and manage: Keep an eye on the deployed deception host and make necessary adjustments to the configuration as needed.

```Note: Since Sinon is a proof-of-concept, it is recommended to use pre-generated content and avoid storing sensitive information like API keys on deception hosts in production environments.```
