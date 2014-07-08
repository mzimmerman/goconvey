package system

import (
	"fmt"
	"testing"

	. "github.com/smartystreets/goconvey/convey"
)

func TestProfileParsing(t *testing.T) {
	for i, test := range testCases {
		if test.SKIP {
			SkipConvey(fmt.Sprintf("Profile Parsing, Test Case #%d: %s (SKIPPED)", i, test.description), t, nil)
		} else {
			Convey(fmt.Sprintf("Profile Parsing, Test Case #%d: %s", i, test.description), t, func() {
				ignored, testArgs := parseProfile(test.input)

				So(ignored, ShouldEqual, test.resultIgnored)
				So(testArgs, ShouldResemble, test.resultTestArgs)
			})
		}
	}
}

var testCases = []struct {
	SKIP           bool
	description    string
	input          string
	resultIgnored  bool
	resultTestArgs []string
}{
	{
		SKIP:           false,
		description:    "Blank profile",
		input:          "",
		resultIgnored:  false,
		resultTestArgs: []string{},
	},
	{
		SKIP:           false,
		description:    "All lines are blank or whitespace",
		input:          "\n \n \t\t\t  \n \n \n",
		resultIgnored:  false,
		resultTestArgs: []string{},
	},
	{
		SKIP:           false,
		description:    "Ignored package, no args included",
		input:          "IGNORE\n-timeout=4s",
		resultIgnored:  true,
		resultTestArgs: []string{},
	},
	{
		SKIP:          false,
		description:   "Ignore directive is commented, all args are included",
		input:         "#IGNORE\n-timeout=4s\n-parallel=5",
		resultIgnored: false,
		resultTestArgs: []string{
			"-timeout=4s",
			"-parallel=5",
		},
	},
	{
		SKIP:          false,
		description:   "No ignore directive, all args are included",
		input:         "-run=TestBlah\n-timeout=42s",
		resultIgnored: false,
		resultTestArgs: []string{
			"-run=TestBlah",
			"-timeout=42s",
		},
	},
	{
		SKIP:          false,
		description:   "Some args are commented, therefore ignored",
		input:         "-run=TestBlah\n#-timeout=42s",
		resultIgnored: false,
		resultTestArgs: []string{
			"-run=TestBlah",
		},
	},
	{
		SKIP:           false,
		description:    "All args are commented, therefore all are ignored",
		input:          "#-run=TestBlah\n//-timeout=42",
		resultIgnored:  false,
		resultTestArgs: []string{},
	},
	{
		SKIP:           false,
		description:    "We ignore certain flags like -v and -cover* because they are specified by the shell",
		input:          "-v\n-cover\n-coverprofile=blah.out",
		resultIgnored:  false,
		resultTestArgs: []string{},
	},
}

/////////////////////////////////////////////////////////////////////////

type FakeProfiles struct {
	recordedRefreshes []string
	recordedIgnored   []string
	providedFlags     map[string][]string
}

func (self *FakeProfiles) Refresh(path string)        {}
func (self *FakeProfiles) IsIgnored(path string) bool { return false }
func (self *FakeProfiles) GoTestFlags(path string) []string {
	return []string{}
}

func NewFakeProfiles() *FakeProfiles {
	self := new(FakeProfiles)
	self.recordedRefreshes = []string{}
	self.recordedIgnored = []string{}
	self.providedFlags = map[string][]string{}
	return self
}
