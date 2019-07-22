package version

import "net/http"

type versionRes struct {
	Service   string `json:"service"`
	GoVersion string `json:"GOVERSION"`
	GitCommit string `json:"commit"`
	GOOS      string `json:"GOOS"`
	GOARCH    string `json:"GOARCH"`
}

func (res versionRes) Code() int {
	return http.StatusOK
}

func (res versionRes) Headers() map[string]string {
	return map[string]string{}
}

func (res versionRes) Empty() bool {
	return false
}
