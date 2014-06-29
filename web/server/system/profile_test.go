package system

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
