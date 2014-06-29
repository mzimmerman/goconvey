package system

import (
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
	"time"
)

////////////////////////////////////////////////////////////////////////

type Profiles interface {
	Refresh(path string)
	IsIgnored(path string) bool
	GoTestFlags(path string) []string
}

////////////////////////////////////////////////////////////////////////

type ProfileCache struct {
	listing map[string]*PackageProfile
}

func (self *ProfileCache) Refresh(path string) {
	profile, exists := self.listing[path]
	if !exists {
		profile = NewPackageProfile(path)
		self.listing[path] = profile
	}
	profile.Refresh()
}

func (self *ProfileCache) IsIgnored(path string) bool {
	profile, exists := self.listing[path]
	return exists && profile.IsIgnored
}

func (self *ProfileCache) GoTestFlags(path string) []string {
	profile := self.listing[path]
	if profile == nil {
		return []string{}
	}

	destination := make([]string, len(profile.GoTestFlags))
	copy(destination, profile.GoTestFlags)
	return destination
}

func NewProfileCache() *ProfileCache {
	self := new(ProfileCache)
	self.listing = map[string]*PackageProfile{}
	return self
}

/////////////////////////////////////////////////////////////////////

type PackageProfile struct {
	IsIgnored   bool
	GoTestFlags []string

	lastModified time.Time
	folder       string
	path         string
}

func (self *PackageProfile) Refresh() {
	self.path = identifyProfile(self.folder)

	info, err := os.Stat(self.path)
	if err != nil {
		return
	}

	modified := info.ModTime()
	if modified == self.lastModified {
		return
	}

	self.lastModified = modified
	profile, err := ioutil.ReadFile(self.path)
	if err == nil {
		self.IsIgnored, self.GoTestFlags = parseProfile(string(profile))
	}
}

func NewPackageProfile(folder string) *PackageProfile {
	self := new(PackageProfile)
	self.folder = folder
	self.GoTestFlags = []string{}
	self.Refresh()
	return self
}
func identifyProfile(folder string) string {
	files, err := ioutil.ReadDir(folder)
	if err == nil {
		for _, file := range files {
			name := file.Name()
			if strings.HasSuffix(name, ".goconvey") {
				return filepath.Join(folder, name)
			}
		}
	}
	return filepath.Join(folder, ".goconvey")
}

// TODO: unit test
func parseProfile(profile string) (IsIgnored bool, arguments []string) {
	lines := strings.Split(profile, "\n")
	arguments = []string{}

	for _, line := range lines {
		if len(arguments) == 0 && strings.ToLower(line) == "ignore" {
			return true, []string{}

		} else if len(strings.TrimSpace(line)) == 0 {
			continue

		} else if strings.HasPrefix(line, "#") {
			continue

		} else if strings.HasPrefix(line, "//") {
			continue

		} else if strings.HasPrefix(line, "-cover") {
			continue // TODO: enable custom coverage flags...

		} else if line == "-v" {
			continue // Verbose mode is always enabled so there is no need to record it here.

		}

		arguments = append(arguments, line)
	}

	return false, arguments
}
