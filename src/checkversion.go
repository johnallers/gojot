package sdees

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/exec"
	"path"
	"path/filepath"
	"strconv"
	"strings"
	"time"
)

type GithubJson struct {
	URL             string `json:"url"`
	AssetsURL       string `json:"assets_url"`
	UploadURL       string `json:"upload_url"`
	HTMLURL         string `json:"html_url"`
	ID              int    `json:"id"`
	TagName         string `json:"tag_name"`
	TargetCommitish string `json:"target_commitish"`
	Name            string `json:"name"`
	Draft           bool   `json:"draft"`
	Author          struct {
		Login             string `json:"login"`
		ID                int    `json:"id"`
		AvatarURL         string `json:"avatar_url"`
		GravatarID        string `json:"gravatar_id"`
		URL               string `json:"url"`
		HTMLURL           string `json:"html_url"`
		FollowersURL      string `json:"followers_url"`
		FollowingURL      string `json:"following_url"`
		GistsURL          string `json:"gists_url"`
		StarredURL        string `json:"starred_url"`
		SubscriptionsURL  string `json:"subscriptions_url"`
		OrganizationsURL  string `json:"organizations_url"`
		ReposURL          string `json:"repos_url"`
		EventsURL         string `json:"events_url"`
		ReceivedEventsURL string `json:"received_events_url"`
		Type              string `json:"type"`
		SiteAdmin         bool   `json:"site_admin"`
	} `json:"author"`
	Prerelease  bool      `json:"prerelease"`
	CreatedAt   time.Time `json:"created_at"`
	PublishedAt time.Time `json:"published_at"`
	Assets      []struct {
		URL      string `json:"url"`
		ID       int    `json:"id"`
		Name     string `json:"name"`
		Label    string `json:"label"`
		Uploader struct {
			Login             string `json:"login"`
			ID                int    `json:"id"`
			AvatarURL         string `json:"avatar_url"`
			GravatarID        string `json:"gravatar_id"`
			URL               string `json:"url"`
			HTMLURL           string `json:"html_url"`
			FollowersURL      string `json:"followers_url"`
			FollowingURL      string `json:"following_url"`
			GistsURL          string `json:"gists_url"`
			StarredURL        string `json:"starred_url"`
			SubscriptionsURL  string `json:"subscriptions_url"`
			OrganizationsURL  string `json:"organizations_url"`
			ReposURL          string `json:"repos_url"`
			EventsURL         string `json:"events_url"`
			ReceivedEventsURL string `json:"received_events_url"`
			Type              string `json:"type"`
			SiteAdmin         bool   `json:"site_admin"`
		} `json:"uploader"`
		ContentType        string    `json:"content_type"`
		State              string    `json:"state"`
		Size               int       `json:"size"`
		DownloadCount      int       `json:"download_count"`
		CreatedAt          time.Time `json:"created_at"`
		UpdatedAt          time.Time `json:"updated_at"`
		BrowserDownloadURL string    `json:"browser_download_url"`
	} `json:"assets"`
	TarballURL string `json:"tarball_url"`
	ZipballURL string `json:"zipball_url"`
	Body       string `json:"body"`
}

func CheckNewVersion(program string, version string, osType string) {
	dir, err := filepath.Abs(filepath.Dir(program))
	if err != nil {
		log.Fatal("Could not get filepath: " + err.Error())
	}
	logger.Debug("Current executable path: %s", dir)

	newVersion, versionName := checkGithub(version)
	if !newVersion {
		return
	}
	var yesnoall string
	fmt.Printf("\nVersion %s is available. Download? (y/n) ", versionName)
	fmt.Scanln(&yesnoall)
	if yesnoall == "n" {
		return
	}
	downloadVersion := versionName + "/sdees_" + osType + ".zip"
	os.Remove("sdees_" + osType + ".zip")
	fmt.Printf("\nDownloading %s...", downloadVersion)
	cmd := exec.Command("wget", "https://github.com/schollz/sdees/releases/download/"+downloadVersion)
	_, err = cmd.Output()
	if err != nil {
		logger.Error("Problem downloading, do you have internet?")
		log.Fatal(err)
	}

	logger.Debug("Removing old version: %s", path.Join(dir, program))
	err = os.Remove(path.Join(dir, program))
	if err != nil {
		logger.Error("Problem removing file, do you need sudo?")
		log.Fatal(err)
	}

	logger.Debug("Unzipping new version")
	cmd = exec.Command("unzip", "sdees_"+osType+".zip")
	_, err = cmd.Output()
	if err != nil {
		logger.Error("Problem unzipping, do you have zip?")
		log.Fatal(err)
	}

	logger.Debug("Cleaning...")
	os.Remove("sdees_" + osType + ".zip")
	fmt.Printf("\n\nsdees Version %s installed!\n", versionName)
	os.Exit(0)
}

func checkGithub(version string) (bool, string) {
	url := "https://api.github.com/repos/schollz/sdees/releases/latest"
	r, err := http.Get(url)
	if err != nil {
		log.Fatal(err)
	}
	defer r.Body.Close()
	var j GithubJson
	err = json.NewDecoder(r.Body).Decode(&j)
	if err != nil {
		log.Fatal(err)
	}
	newVersion := j.TagName
	versions := strings.Split(newVersion, ".")
	if len(versions) != 3 {
		return false, ""
	}
	majorMinorWeb := []int{}
	for i := 0; i < 3; i++ {
		i, _ := strconv.Atoi(versions[i])
		majorMinorWeb = append(majorMinorWeb, i)
	}

	versions = strings.Split(version, ".")
	if len(versions) != 3 {
		return false, ""
	}
	majorMinor := []int{}
	for i := 0; i < 3; i++ {
		i, _ := strconv.Atoi(versions[i])
		majorMinor = append(majorMinor, i)
	}

	newVersionAvailable := false
	for i := range majorMinor {
		if majorMinor[i] > majorMinorWeb[i] {
			break
		}
		newVersionAvailable = true
	}

	return newVersionAvailable, newVersion
}
