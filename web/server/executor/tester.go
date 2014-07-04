package executor

import (
	"log"

	"github.com/smartystreets/goconvey/web/server/contract"
)

type ConcurrentTester struct {
	shell     contract.Shell
	batchSize int
}

func (self *ConcurrentTester) SetBatchSize(batchSize int) {
	self.batchSize = batchSize
	log.Printf("Now configured to test %d packages concurrently.\n", self.batchSize)
}

func (self *ConcurrentTester) TestAll(folders []*contract.Package) {
	if self.batchSize == 1 {
		self.executeSynchronously(folders)
	} else {
		newCuncurrentCoordinator(folders, self.batchSize, self.shell).ExecuteConcurrently()
	}
}

func (self *ConcurrentTester) executeSynchronously(folders []*contract.Package) {
	for _, folder := range folders {
		if !folder.Active {
			log.Printf("Skipping execution: %s\n", folder.Name)
			continue
		}

		log.Printf("Executing tests: %s\n", folder.Name)

		folder.Output, folder.Error = self.shell.GoTest(folder.Path, folder.Name)
		if folder.Error == contract.ProfileSkipsPackage { // TODO: unit test
			log.Println("Marking folder as ignored according to profile directive:", folder.Name)
			folder.Active = false
		} else if folder.Error != nil && folder.Output == "" {
			log.Println("Unexpected error at", folder.Path)
			panic(folder.Error)
		}
	}
}

func NewConcurrentTester(shell contract.Shell) *ConcurrentTester {
	self := new(ConcurrentTester)
	self.shell = shell
	self.batchSize = defaultBatchSize
	return self
}

const defaultBatchSize = 10
