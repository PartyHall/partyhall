package services

import (
	"fmt"
	"io"
	"strings"

	"github.com/partyhall/partyhall/logs"
)

func InjectHtmlMode(mode string) string {
	script := fmt.Sprintf(`<script>window.SOCKET_TYPE = '%v';</script>`, mode)

	//#region Load the original HTML file
	if WEBAPP_FS == nil {
		logs.Errorf("Can't inject mode as the webapp is not built")
		return ""
	}

	file, err := (*WEBAPP_FS).Open("index.html")
	if err != nil {
		logs.Errorf("Can't inject mode: failed to open index.html")
		return ""
	}

	html, err := io.ReadAll(file)
	if err != nil {
		logs.Errorf("Can't inject mode: failed to read index.html")
	}

	//#endregion

	//#region Injecting the socket type
	content := string(html)
	content = strings.ReplaceAll(content, "<!-- INJECT MODE HERE -->", script)
	//#endregion

	return content
}
