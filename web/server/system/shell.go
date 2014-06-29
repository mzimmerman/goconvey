package system

import (
	"os"
	"path/filepath"
	"strings"

	"github.com/smartystreets/goconvey/web/server/contract"
)

type Shell struct {
	executor    Executor
	profiles    Profiles
	coverage    bool
	gobin       string
	reportsPath string
}

func (self *Shell) GoTest(directory, packageName string) (output string, err error) {
	self.profiles.Refresh(directory)

	if self.profiles.IsIgnored(directory) {
		return "", contract.ProfileSkipsPackage
	}

	output, err = self.compilePackageDependencies(directory)
	if err == nil {
		output, err = self.goTest(directory, packageName)
	}

	return
}

func (self *Shell) compilePackageDependencies(directory string) (output string, err error) {
	return self.executor.Execute(directory, self.gobin, "test", "-i")
}

func (self *Shell) goTest(directory, packageName string) (output string, err error) {
	if !self.coverage {
		return self.runWithoutCoverage(directory, packageName)
	}

	return self.tryRunWithCoverage(directory, packageName)
}

func (self *Shell) tryRunWithCoverage(directory, packageName string) (output string, err error) {
	coverageReport := self.composeCoverageReportPath(packageName)
	output, err = self.runWithCoverage(directory, packageName, coverageReport+".txt")

	if err != nil && self.coverage {
		output, err = self.runWithoutCoverage(directory, packageName)
	} else if self.coverage {
		self.generateCoverageReports(directory, coverageReport+".txt", coverageReport+".html")
	}
	return
}

func (self *Shell) composeCoverageReportPath(packageName string) string {
	reportFilename := strings.Replace(packageName, "/", "-", -1)
	reportPath := filepath.Join(self.reportsPath, reportFilename)
	return reportPath
}

func (self *Shell) runWithCoverage(directory, packageName, coverageReport string) (string, error) {
	arguments := []string{"test", "-v", "-covermode=set", "-coverprofile=" + coverageReport}
	arguments = append(arguments, self.jsonFlag(directory, packageName)...)
	arguments = append(arguments, self.profiles.GoTestFlags(directory)...) // TODO: Filter out -coverprofile flag if it is redefined in the profile.
	return self.executor.Execute(directory, self.gobin, arguments...)
}

func (self *Shell) runWithoutCoverage(directory, packageName string) (string, error) {
	arguments := []string{"test", "-v"}
	arguments = append(arguments, self.jsonFlag(directory, packageName)...)
	arguments = append(arguments, self.profiles.GoTestFlags(directory)...) // TODO: Filter out any coverage flags specified in the profile.
	return self.executor.Execute(directory, self.gobin, arguments...)
}

func (self *Shell) jsonFlag(directory, packageName string) []string {
	imports, err := self.executor.Execute(directory, self.gobin, "list", "-f", "'{{.TestImports}}'", packageName)
	if !strings.Contains(imports, goconveyDSLImport) && err == nil {
		return []string{}
	}
	return []string{"-json"}
}

func (self *Shell) generateCoverageReports(directory, coverageReport, html string) {
	self.executor.Execute(directory, self.gobin, "tool", "cover", "-html="+coverageReport, "-o", html)
}

func (self *Shell) Getenv(key string) string {
	return os.Getenv(key)
}

func (self *Shell) Setenv(key, value string) error {
	if self.Getenv(key) != value {
		return os.Setenv(key, value)
	}
	return nil
}

func NewShell(executor Executor, profiles Profiles, gobin string, cover bool, reports string) *Shell {
	self := new(Shell)
	self.executor = executor
	self.profiles = profiles
	self.gobin = gobin
	self.coverage = cover
	self.reportsPath = reports
	return self
}

const goconveyDSLImport = "github.com/smartystreets/goconvey/convey " // note the trailing space: we don't want to target packages nested in the /convey package.
