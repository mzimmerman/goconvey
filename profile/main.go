// will run `go test` on github.com/smartystreets/goconvey/examples
// with any arguments in the `.goconvey` profile for that project.
package main

import (
	"fmt"
	"io/ioutil"
	"os/exec"
	"strings"
)

func main() {
	rawProfile, err := ioutil.ReadFile("../examples/.goconvey")
	if err != nil {
		panic(err)
	}
	profile := string(rawProfile)
	fmt.Println(profile)
	lines := strings.Split(profile, "\n")
	arguments := []string{"test", "github.com/smartystreets/goconvey/examples", "-v"}

	for number, line := range lines {
		if number == 0 && strings.ToLower(line) == "ignore" {
			fmt.Println("Ignoring package")
			return
		}

		if len(line) == 0 {
			continue
		} else if strings.HasPrefix(line, "#") {
			continue
		} else if strings.HasPrefix(line, "//") {
			continue
		}

		arguments = append(arguments, line)
	}

	fmt.Println("go", arguments)

	output, err := exec.Command("go", arguments...).CombinedOutput()
	if err != nil {
		fmt.Println("ERROR:", err.Error())
	}
	fmt.Println(string(output))
}
