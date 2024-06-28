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
  selection_method: random

websites:
  options:
    - https://www.example.com
    - https://www.google.com
    - https://www.github.com
  selection_method: random

preferences:
  default_browser:
    options: [googlechrome, firefox, edge]
    selection_method: random
  background_images:
    location: "http://example.com/backgrounds"
    type: "http"
    selection_method: random
  screen_resolutions:
    options: ["1920x1080", "1280x720"]
    selection_method: random
  languages:
    options: ["en-US", "es-ES", "fr-FR"]
    selection_method: random

start_menu_items:
  options:
    - name: Chrome
      path: "C:\\Program Files\\Google\\Chrome\\Application\\chrome.exe"
    - name: VLC
      path: "C:\\Program Files\\VideoLAN\\VLC\\vlc.exe"
  selection_method: random

file_operations:
  create_modify_files:
    options:
      - path: "C:\\path\\to\\file1.txt"
        content: "This is a test file."
        use_gpt: false
      - path: "C:\\path\\to\\file2.docx"
        content: "Generate a report on sales data."
        use_gpt: true
        gpt_prompt: "Generate a detailed sales report."
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
      - send_to: "example1@example.com"
        subject: "Test Email 1"
        body: "This is a test email."
        use_gpt: false
      - send_to: "example2@example.com"
        subject: "Test Email 2"
        body: "Generate a thank you email."
        use_gpt: true
        gpt_prompt: "Generate a thank you email."
    selection_method: random

software_management:
  options:
    - install: "some_software"
    - uninstall: "another_software"
  selection_method: random

system_updates:
  options: ["check_updates", "apply_updates"]
  selection_method: random

user_accounts:
  options:
    - name: "User03"
      password: "P@ssw0rd"
      full_name: "Third User"
      description: "Description of this account."
  selection_method: random

network_settings:
  options:
    - ssid: "Network1"
      password: "Password1"
    - ssid: "Network2"
      password: "Password2"
  selection_method: random

system_logs:
  options: ["read_logs", "generate_logs"]
  selection_method: random

media_files:
  options:
    - "C:\\path\\to\\media1.mp4"
    - "C:\\path\\to\\media2.mp3"
  selection_method: random

printing:
  options:
    - "C:\\path\\to\\document1.pdf"
    - "C:\\path\\to\\document2.txt"
  selection_method: random

scheduled_tasks:
  options: ["create_task", "modify_task", "delete_task"]
  selection_method: random

decoy_files:
  location: ["http://example.com/decoy_files"]
  type: "http"
  target_directory: "C:\\path\\to\\downloaded\\decoy_files"

interaction_duration: 3600
action_delay: 5
randomness_factor: 2

openai_api_key: "your-openai-api-key"
```

## Deploying Windows Deception Hosts
Sinon is designed to automate the setup of deception hosts by performing a variety of actions that simulate real user activity. The goal is to create a realistic environment that can deceive potential intruders. The modular and configurable nature of Sinon allows for easy adjustments and randomization, making each deployment unique.

## Steps to Deploy
- Prepare the Windows environment: Ensure that the target Windows machine is ready and accessible.
- Configure Sinon: Edit the config.yaml file to define the desired behaviors and settings.
- Run Sinon: Execute the compiled Sinon binary to start the automation process.
- Monitor and manage: Keep an eye on the deployed deception host and make necessary adjustments to the configuration as needed.

```Note: Since Sinon is a proof-of-concept, it is recommended to use pre-generated content and avoid storing sensitive information like API keys on deception hosts in production environments.```
