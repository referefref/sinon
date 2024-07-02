package main

import (
	"fmt"
	"log"
	"os/exec"
	"syscall"
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

func changePreferences(config *Config) {
	prefs := config.Preferences
	browsers := selectRandomOrHardcoded(prefs.DefaultBrowser.Options, prefs.DefaultBrowser.SelectionMethod)
	for _, browser := range browsers {
		cmd := exec.Command("SetDefaultBrowser.exe", "HKLM", browser)
		err := cmd.Run()
		if err != nil {
			logToFile(config.General.LogFile, fmt.Sprintf("Failed to set default browser: %v", err))
		} else {
			logToFile(config.General.LogFile, fmt.Sprintf("Default browser set to %s", browser))
		}
	}

	images := selectRandomOrHardcoded(prefs.BackgroundImages.Options, prefs.BackgroundImages.SelectionMethod)
	for _, image := range images {
		err := downloadFile(image, "C:\\path\\to\\downloaded\\image.jpg")
		if err != nil {
			logToFile(config.General.LogFile, fmt.Sprintf("Failed to download background image: %v", err))
		}
		cmd := exec.Command("powershell", "-Command", "(New-Object -ComObject WScript.Shell).RegWrite('HKCU\\Control Panel\\Desktop\\Wallpaper', 'C:\\path\\to\\downloaded\\image.jpg')")
		err = cmd.Run()
		if err != nil {
			logToFile(config.General.LogFile, fmt.Sprintf("Failed to set background image: %v", err))
		} else {
			logToFile(config.General.LogFile, "Background image set successfully")
		}
	}

	resolutions := selectRandomOrHardcoded(prefs.ScreenResolutions.Options, prefs.ScreenResolutions.SelectionMethod)
	for _, resolution := range resolutions {
		setScreenResolution(resolution)
	}

	languages := selectRandomOrHardcoded(prefs.Languages.Options, prefs.Languages.SelectionMethod)
	for _, language := range languages {
		cmd := exec.Command("powershell", "-Command", fmt.Sprintf(`Set-WinUILanguageOverride -Language "%s"; Set-WinSystemLocale -SystemLocale "%s"; Set-WinUserLanguageList -LanguageList "%s" -Force; Set-Culture -CultureInfo "%s"; Set-WinHomeLocation -GeoId 12`, language, language, language, language))	
		err := cmd.Run()
		if err != nil {
			logToFile(config.General.LogFile, fmt.Sprintf("Failed to set language to %s: %v", language, err))
		} else {
			logToFile(config.General.LogFile, fmt.Sprintf("Language set to %s", language))
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

