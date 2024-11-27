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
	Framework []RuntimeVersion `json:"Framework"`
	Language  []RuntimeVersion `json:"Language"`
	Custom    []RuntimeVersion `json:"Custom"`
	OS        []RuntimeVersion `json:"OS"`
}

type Config struct {
	Runtime Runtime `json:"runtime"`
}
