package modrinth

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"
	"time"
)

const (
	baseURL   = "https://api.modrinth.com/v2"
	userAgent = "Higuruma/0.1.0 (github.com/Muchprow/higuruma)"
)

var client = &http.Client{Timeout: 15 * time.Second}

type Project struct {
	ProjectID    string   `json:"project_id"`
	Slug         string   `json:"slug"`
	Title        string   `json:"title"`
	Description  string   `json:"description"`
	IconURL      string   `json:"icon_url"`
	Downloads    int      `json:"downloads"`
	Categories   []string `json:"categories"`
	Loaders      []string `json:"loaders"`
	GameVersions []string `json:"game_versions"`
}

type SearchResult struct {
	Hits  []Project `json:"hits"`
	Total int       `json:"total_hits"`
}

type Version struct {
	ID            string   `json:"id"`
	ProjectID     string   `json:"project_id"`
	Name          string   `json:"name"`
	VersionNumber string   `json:"version_number"`
	GameVersions  []string `json:"game_versions"`
	Loaders       []string `json:"loaders"`
	Files         []File   `json:"files"`
}

type File struct {
	URL      string `json:"url"`
	Filename string `json:"filename"`
	Primary  bool   `json:"primary"`
	Size     int    `json:"size"`
}

func doRequest(endpoint string, result any) error {
	req, err := http.NewRequest("GET", baseURL+endpoint, nil)
	if err != nil {
		return fmt.Errorf("modrinth: create request: %w", err)
	}

	req.Header.Set("User-Agent", userAgent)

	resp, err := client.Do(req)
	if err != nil {
		return fmt.Errorf("modrinth: request failed: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return fmt.Errorf("modrinth: status %d: %s", resp.StatusCode, string(body))
	}

	if err := json.NewDecoder(resp.Body).Decode(result); err != nil {
		return fmt.Errorf("modrinth: decode response: %w", err)
	}

	return nil
}

func Search(query string, limit int, offset int) (*SearchResult, error) {
	if limit <= 0 {
		limit = 20
	}

	params := url.Values{}
	params.Set("query", query)
	params.Set("limit", fmt.Sprintf("%d", limit))
	params.Set("offset", fmt.Sprintf("%d", offset))

	endpoint := "/search?" + params.Encode()

	var result SearchResult
	if err := doRequest(endpoint, &result); err != nil {
		return nil, err
	}

	return &result, nil
}

func SearchWithFilters(query string, gameVersion string, loader string, limit int, offset int) (*SearchResult, error) {
	if limit <= 0 {
		limit = 20
	}

	params := url.Values{}
	params.Set("query", query)
	params.Set("limit", fmt.Sprintf("%d", limit))
	params.Set("offset", fmt.Sprintf("%d", offset))

	if gameVersion != "" {
		params.Set("facets", fmt.Sprintf(`[["versions:%s"],["categories:%s"]]`, gameVersion, loader))
	} else if loader != "" {
		params.Set("facets", fmt.Sprintf(`[["categories:%s"]]`, loader))
	}

	endpoint := "/search?" + params.Encode()

	var result SearchResult
	if err := doRequest(endpoint, &result); err != nil {
		return nil, err
	}

	return &result, nil
}

func GetProject(id string) (*Project, error) {
	var project Project
	if err := doRequest("/project/"+id, &project); err != nil {
		return nil, err
	}
	return &project, nil
}

func GetVersions(projectID string, gameVersion string, loader string) ([]Version, error) {
	params := url.Values{}

	if gameVersion != "" {
		params.Set("game_versions", fmt.Sprintf(`["%s"]`, gameVersion))
	}
	if loader != "" {
		params.Set("loaders", fmt.Sprintf(`["%s"]`, loader))
	}

	endpoint := "/project/" + projectID + "/version"
	if len(params) > 0 {
		endpoint += "?" + params.Encode()
	}

	var versions []Version
	if err := doRequest(endpoint, &versions); err != nil {
		return nil, err
	}

	return versions, nil
}

func DownloadFile(fileURL string, dest string) error {
	req, err := http.NewRequest("GET", fileURL, nil)
	if err != nil {
		return fmt.Errorf("modrinth: create download request: %w", err)
	}

	req.Header.Set("User-Agent", userAgent)

	resp, err := client.Do(req)
	if err != nil {
		return fmt.Errorf("modrinth: download failed: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("modrinth: download status %d", resp.StatusCode)
	}

	out, err := os.Create(dest)
	if err != nil {
		return fmt.Errorf("modrinth: create file: %w", err)
	}
	defer out.Close()

	if _, err := io.Copy(out, resp.Body); err != nil {
		return fmt.Errorf("modrinth: write file: %w", err)
	}

	return nil
}
