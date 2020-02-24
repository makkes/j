package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"strconv"

	"github.com/makkes/j"

	"github.com/spf13/pflag"
)

func writeHistory(history map[string]int) error {
	marshalledHistory, err := json.Marshal(history)
	if err != nil {
		return err
	}
	if err := ioutil.WriteFile(os.Getenv("HOME")+"/.jump.json", marshalledHistory, 0644); err != nil {
		return err
	}
	return nil
}

func readHistory() (map[string]int, j.JumpFile, error) {
	jumpBytes, err := ioutil.ReadFile(os.Getenv("HOME") + "/.jump.json")
	if err != nil {
		return nil, nil, err
	}
	var jumpFileMap map[string]int
	if err := json.Unmarshal(jumpBytes, &jumpFileMap); err != nil {
		return nil, nil, err
	}
	var jumpFile j.JumpFile
	for path, freq := range jumpFileMap {
		jumpFile = append(jumpFile, j.JumpFileEntry{
			Path: path,
			Freq: freq,
		})
	}

	return jumpFileMap, jumpFile, nil
}

func main() {
	complete := pflag.Bool("complete", false, "Complete STR and show all possible entries that match. This option does not update J's database.")
	debug := pflag.BoolP("debug", "d", false, "Print debug information in conjunction with --complete")
	pflag.Parse()
	query := ""
	if len(pflag.Args()) > 0 {
		query = pflag.Args()[0]
	}

	if query == "" {
		fmt.Println(os.Getenv("HOME"))
		return
	}

	jumpFileMap, jumpFile, err := readHistory()
	if err != nil {
		panic(err)
	}
	matchResult := j.NewJumper().FindCandidates(jumpFile, query)
	if matchResult.Len() == 0 {
		if *complete {
			os.Exit(0)
		}
		fmt.Fprintf(os.Stderr, "Error: No matching directory found for '%s'\n", query)
		os.Exit(1)
	}

	if *complete {
		matchResult.StripCommonPrefix()
		if matchResult.Len() == 1 && !*debug {
			// don't print the leading order number since the shell will complete the single match in-place
			fmt.Printf("%s\n", matchResult.Matches[0].Path)
			return
		}
		for idx, match := range matchResult.Matches {
			if *debug {
				fmt.Printf("%d %d %s\n", match.Freq, match.Score, match.Path)
			} else {
				fmt.Printf("%0*d %s\n", len(strconv.Itoa(matchResult.Len()))-1, idx, match.Path)
			}
		}
	} else {
		fmt.Printf("%s\n", matchResult.Matches[0].Path)
		jumpFileMap[matchResult.Matches[0].AbsPath]++
		writeHistory(jumpFileMap)
	}

}
