package build

import (
	"io"
	"net/http"
)

type GitHubJarTemplate struct{}

func (t *GitHubJarTemplate) Jar() (io.ReadCloser, error) {

	// Get the JAR url
	jarUrl := PluginJarUrl()

	// Download the JAR file
	res, err := http.Get(jarUrl)
	if err != nil {
		return nil, err
	}

	// Return the response body
	return res.Body, nil
}

func PluginJarUrl() string {
	return "https://github.com/projopenrealms/bukkit-runtime/releases/latest/download/bukkit-runtime.jar"
}
