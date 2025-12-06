package apps

import (
	"encoding/json"
	"fmt"
	"ghost-app-store-tui/dto"
	"io"
	"net/http"
)

type App struct {
	Id       string       `json:"id"`
	Name     string       `json:"name"`
	Versions []AppVersion `json:"versions"`
}

type AppVersion struct {
	Id      string `json:"id"`
	Version string `json:"version"`
	Icon    string `json:"icon"`
}

func GetApps() []dto.App {
	resp, err := http.Get(fmt.Sprintf("http://localhost:8084/v1/apps"))
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		panic(err)
	}

	var appBody []dto.App

	err = json.Unmarshal(body, &appBody)
	if err != nil {
		panic(err)
	}

	return appBody
}
