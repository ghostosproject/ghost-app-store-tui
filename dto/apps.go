package dto

type App struct {
	Id       string    `json:"id"`
	Name     string    `json:"name"`
	Icon     string    `json:"icon"`
	Repo     string    `json:"repo"`
	Project  string    `json:"project"`
	Versions []Version `json:"versions"`
}

type Version struct {
	Id      string `json:"id"`
	Version string `json:"version"`
	Url     string `json:"url"`
	Digest  string `json:"digest"`
}
