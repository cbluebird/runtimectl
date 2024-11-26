package model

type Version struct {
	Name   string `json:"name"`
	Image  string `json:"image"`
	Config string `json:"config"`
}

type RuntimeVersion struct {
	Name    string    `json:"name"`
	Version []Version `json:"version"`
}

type Runtime struct {
	Framework []RuntimeVersion `json:"framework"`
	Language  []RuntimeVersion `json:"language"`
	Custom    []RuntimeVersion `json:"custom"`
	OS        []RuntimeVersion `json:"os"`
}

type Config struct {
	Runtime Runtime `json:"runtime"`
}
