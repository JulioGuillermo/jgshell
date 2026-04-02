package statusbar

import "charm.land/lipgloss/v2"

func GetOSIcon(os string) string {
	//              
	switch os {
	case "linux":
		return ""
	case "android":
		return ""
	case "ios", "mac":
		return ""
	case "freebsd", "openbsd", "netbsd":
		return ""
	case "windows":
		return ""
	default:
		return ""
	}
}

func GetOS(os string) string {
	color := "#ff9800"
	switch os {
	case "linux":
		color = "#ffffff"
	case "android":
		color = "#00ff00"
	case "ios", "mac":
		color = "#aaaaaa"
	case "freebsd", "openbsd", "netbsd":
		color = "#ff3333"
	case "windows":
		color = "#0088ff"
	}
	return lipgloss.NewStyle().
		Foreground(lipgloss.Color(color)).
		Render(GetOSIcon(os))
}
