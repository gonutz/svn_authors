// svn_authors is a tool to help migrate SVN repositories to Git on Windows.
// There is a manual at:
//
//     https://gist.github.com/NathanSweet/7327535
//
// which  uses awk to extract a list of authors from the SVN log. Instead of
// using awk on Windows, this script does the job as well. It takes the output
// of
//
//     svn log -q
//
// on stdin and outputs a lexicographically sorted list of authors on stdout.
// This was tested with SVN version 1.12.2. It might not work for other versions
// if the log output format is different.

package main

import (
	"fmt"
	"io/ioutil"
	"os"
	"sort"
	"strings"
)

func main() {
	names := make(map[string]bool)

	findAuthor := func(line string) {
		// This line starts like this:
		//     r45 | author_name | ...
		// We want to extract the author_name.
		i := strings.Index(line, "| ")
		if i == -1 {
			return
		}
		line = line[i+2:]
		i = strings.Index(line, " |")
		if i == -1 {
			return
		}
		line = line[:i]
		names[line] = true
	}

	data, err := ioutil.ReadAll(os.Stdin)
	check(err)
	lines := strings.Split(strings.Replace(string(data), "\r", "", -1), "\n")
	for i := range lines {
		if lines[i] == "------------------------------------------------------------------------" {
			// The next line is the header of a new commit message, it has the
			// author in the second column.
			findAuthor(lines[i+1])
		}
	}

	var authors []string
	for name, _ := range names {
		authors = append(authors, name)
	}
	sort.Strings(authors)
	for _, name := range authors {
		fmt.Println(name, "=", name, "<"+name+">")
	}
}

func check(err error) {
	if err != nil {
		panic(err)
	}
}
