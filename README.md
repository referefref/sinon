# Sinon - Modular Windows Burn-In Automation with Generative AI for Deception

![image](https://github.com/referefref/sinon/assets/56499429/9c5395ac-2034-4bcb-8214-c3690013d521)

Sinon is a modular tool for automatic burn-in of Windows-based deception hosts that aims to reduce the difficulty of orchestrating deception hosts at scale whilst enabling diversity and randomness through generative capabilities. It has been created as a proof-of-concept and is not intended for production deception environments. It would likely be better suited to having content pre-generated and built into a one-time script, as we wouldn't want to be storing secrets like OpenAI API keys on a decoy or deception host.

## Features

- Generative content including files, emails, and so on using OpenAI API (Configured for GPT-4o)
- Randomness factor - select from list in config, or follow config completely
- Temporal randomness - set delay to execution and delay between events including randomness factor

Sinon performs the following functions, as determined by a config file:

- **Install Applications**: Automatically install applications from a predefined list using Chocolatey.
- **Browse Websites**: Automatically open a list of websites to simulate user activity.
- **Change Preferences**: Modify system preferences such as default browser, background images, screen resolutions, and system languages.
- **Add Start Menu Items**: Add shortcuts to specified applications in the start menu.
- **Create and Modify Files**: Generate and modify text files with the option to use OpenAI GPT-4 for content generation.
- **Send Emails**: Send emails with the option to use OpenAI GPT-4 for content generation.
- **Download Decoy Files**: Download files from specified URLs to simulate decoy file activity.
- **Manage Software**: Install or uninstall software applications using predefined commands.
- **Perform System Updates**: Execute system update commands.
- **Manage User Accounts**: Create and manage user accounts with specified attributes.
- **Manage Network Settings**: Configure Wi-Fi network connections using SSID and password.
- **Open Media Files**: Open media files such as images, videos, and audio files.
- **Print Documents**: Print specified text documents.
- **Create Scheduled Tasks**: Schedule tasks to run specified commands at defined times.
- **Simulate User Interaction**: Control the duration and delay of interactions with randomness.
- **Create Lures**: Generate various types of lures to deceive intruders.
  - Credential pairs
  - SSH keys
  - Website URLs
  - Registry keys
  - CSV documents
  - API keys
  - LNK files (shortcuts)
- **Monitor File System**: Watch specified paths for file system events such as modifications and log these events.
- **Redis Connectivity**: Send session metadata to a Redis server for centralized logging and analysis.

## Usage

1. **Clone the repository:**
    ```sh
    git clone https://github.com/yourusername/sinon.git
    cd sinon
    ```

2. **Configure the application:**
   - Modify the `config.yaml` file to suit your needs. See the [Config Items](#config-items) section for details.

3. **Build the application:**
    ```sh
    go build -o sinon
    # building for windows on linux: GOOS=windows GOARCH=amd64 go build -o sinon.exe
    ```

4. **Deploy the application to your target machine:**
   - This could be accomplished many ways, you may want to burn it in to an image, use SCCM/Intune etc.

## Config Items

The `config.yaml` file contains all the configuration options for Sinon. Here is an example configuration file with explanations:

```yaml
applications:
  options:
    - googlechrome
    - firefox
    - notepadplusplus
    - vlc
  selection_method: random

websites:
  options:
    - https://www.google.com
    - https://www.wikipedia.org
    - https://www.github.com
  selection_method: random

preferences:
  default_browser:
    options:
      - "C:\\Program Files\\Google\\Chrome\\Application\\chrome.exe"
      - "C:\\Program Files\\Mozilla Firefox\\firefox.exe"
    selection_method: random
  background_images:
    location: "C:\\Users\\user\\Pictures"
    type: http
    selection_method: random
    options:
      - https://example.com/background1.jpg
      - https://example.com/background2.jpg
  screen_resolutions:
    options:
      - "1920x1080"
      - "1366x768"
    selection_method: random
  languages:
    options:
      - en-US
      - es-ES
    selection_method: random

start_menu_items:
  options:
    - "C:\\Program Files\\Google\\Chrome\\Application\\chrome.exe"
    - "C:\\Program Files\\Mozilla Firefox\\firefox.exe"
  selection_method: random

file_operations:
  create_modify_files:
    - path: "C:\\Users\\user\\Documents\\example.txt"
      content: "This is an example text file."
      use_gpt: false
      gpt_prompt: ""

email_operations:
  google_account:
    email: "user@gmail.com"
    password: "password"
  microsoft_account:
    email: "user@outlook.com"
    password: "password"
  send_receive:
    - send_to: "recipient@example.com"
      subject: "Test Email"
      body: "This is a test email."
      use_gpt: true
      gpt_prompt: "Write a friendly email to a colleague."

software_management:
  options:
    - upgrade all
    - uninstall vlc
  selection_method: random

system_updates:
  method: install_all
  specific_updates:
    - KB123456
    - KB789012
  selection_method: random
  hide_updates:
    - KB654321
    - KB210987

user_accounts:
  - name: user1
    password: password1
    full_name: User One
    description: First user account
  - name: user2
    password: password2
    full_name: User Two
    description: Second user account

network_settings:
  - ssid: ExampleSSID
    password: examplepassword

system_logs:
  options:
    - Application
    - System
  selection_method: random

media_files:
  location: "C:\\Users\\user\\Videos"
  type: http
  selection_method: random
  options:
    - https://example.com/video1.mp4
    - https://example.com/video2.mp4

printing:
  options:
    - "C:\\Users\\user\\Documents\\print_me.txt"
  selection_method: random

scheduled_tasks:
  options:
    - name: Task1
      path: "C:\\Windows\\System32\\notepad.exe"
      schedule: "daily"
      start_time: "14:00"
    - name: Task2
      path: "C:\\Windows\\System32\\calc.exe"
      schedule: "weekly"
      start_time: "10:00"
  selection_method: random

decoy_files:
  sets:
    - location:
        - "https://example.com/decoy1.txt"
        - "https://example.com/decoy2.txt"
      type: http
      target_directory:
        - "C:\\Users\\user\\Documents"
      selection_method: random

lures:
  - name: CredentialLure
    type: credential_pair
    location: "C:\\Users\\user\\Desktop\\credential.txt"
    generation_params:
      length: 12
    generative_type: golang
    openai_prompt: ""
  - name: SSLLure
    type: ssh_key
    location: "C:\\Users\\user\\Desktop\\id_rsa"
    generation_params: {}
    generative_type: golang
    openai_prompt: ""
  - name: URLLure
    type: website_url
    location: "C:\\Users\\user\\Desktop\\phishing_link.url"
    generation_params:
      base_url: "https://malicious.example.com"
    generative_type: golang
    openai_prompt: ""
  - name: RegistryLure
    type: registry_key
    location: "HKEY_CURRENT_USER\\Software\\ExampleKey"
    generation_params:
      registry_key_type: "REG_SZ"
      registry_key_value: "ExampleValue"
    generative_type: golang
    openai_prompt: ""
  - name: CSVLure
    type: csv
    location: "C:\\Users\\user\\Desktop\\financial_records.csv"
    generation_params:
      document_content: "Date,Amount,Description\n2024-01-01,1000,Salary"
    generative_type: golang
    openai_prompt: ""
  - name: APIKeyLure
    type: api_key
    location: "C:\\Users\\user\\Desktop\\api_key.txt"
    generation_params:
      api_key_format: "uuid"
    generative_type: golang
    openai_prompt: ""
  - name: LNKLure
    type: lnk
    location: "C:\\Users\\user\\Desktop\\shortcut.lnk"
    generation_params:
      target_path: "C:\\Windows\\System32\\notepad.exe"
    generative_type: golang
    openai_prompt: ""

general:
  redis:
    ip: "127.0.0.1"
    port: 6379
  log_file: "C:\\Users\\user\\sinon.log"
  openai_api_key: "your_openai_api_key"
  interaction_duration: 60
  action_delay: 5
  randomness_factor: 2

```

## Deploying Windows Deception Hosts
Sinon is designed to automate the setup of deception hosts by performing a variety of actions that simulate real user activity. The goal is to create a realistic environment that can deceive potential intruders. The modular and configurable nature of Sinon allows for easy adjustments and randomization, making each deployment unique.

## Steps to Deploy
- Prepare the Windows environment: Ensure that the target Windows machine is ready and accessible.
- Configure Sinon: Edit the config.yaml file to define the desired behaviors and settings.
- Run Sinon: Execute the compiled Sinon binary to start the automation process.
- Monitor and manage: Keep an eye on the deployed deception host and make necessary adjustments to the configuration as needed.

```Note: Since Sinon is a proof-of-concept, it is recommended to use pre-generated content and avoid storing sensitive information like API keys on deception hosts in production environments.```
