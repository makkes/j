package j

import (
	"io/ioutil"
	"os"
	"path"
	"path/filepath"
	"sort"
	"strings"

	"github.com/sahilm/fuzzy"
)

type DirExists func(string) bool
type PrefixDir func(string) string

type Jumper struct {
	dirExists DirExists
	prefixDir PrefixDir
}

type Opt func(j *Jumper)

func NewJumper(opts ...Opt) Jumper {
	j := Jumper{
		dirExists: func(p string) bool {
			fi, err := os.Stat(p)
			return err == nil && fi.IsDir()
		},
		prefixDir: func(p string) string {
			wd, err := os.Getwd()
			if err != nil {
				return ""
			}
			fis, err := os.ReadDir(wd)
			if err != nil {
				return ""
			}
			for _, fi := range fis {
				info, err := os.Stat(fi.Name())
				if err != nil {
					return ""
				}
				if info.IsDir() && strings.HasPrefix(fi.Name(), p) {
					return fi.Name()
				}
			}
			return ""
		},
	}
	for _, opt := range opts {
		opt(&j)
	}
	return j
}

type JumpFileEntry struct {
	Path string
	Freq int
}

type JumpFile []JumpFileEntry

// The string to be matched at position i.
func (f JumpFile) String(i int) string {
	return f[i].Path
}

// The length of the source. Typically is the length of the slice of things that you want to match.
func (f JumpFile) Len() int {
	return len(f)
}

type Match struct {
	Path     string
	AbsPath  string
	Freq     int
	Score    int
	Priority bool
}

type Matches []*Match

func (ms Matches) Has(path string) bool {
	for _, m := range ms {
		if m.Path == path {
			return true
		}
	}
	return false
}

type MatchResult struct {
	Matches Matches
	Query   string
}

func NewMatch(path string, abspath string) *Match {
	return &Match{
		Path:     path,
		AbsPath:  abspath,
		Priority: true,
	}
}

// Less reports whether the element with
// index i should sort before the element with index j.
func (m MatchResult) Less(i int, j int) bool {
	if m.Matches[i].Priority && !m.Matches[j].Priority {
		return true
	}
	if m.Matches[j].Priority && !m.Matches[i].Priority {
		return false
	}

	if strings.HasSuffix(m.Matches[i].Path, m.Query) {
		if strings.HasSuffix(m.Matches[j].Path, m.Query) {
			if m.Matches[i].Freq > m.Matches[j].Freq {
				return true
			}
			return m.Matches[i].Score > m.Matches[j].Score
		}
		return true
	}

	if m.Matches[i].Score == m.Matches[j].Score {
		return m.Matches[i].Freq > m.Matches[j].Freq
	}

	return m.Matches[i].Score > m.Matches[j].Score

}

// Swap swaps the elements with indexes i and j.
func (m MatchResult) Swap(i int, j int) {
	m.Matches[i], m.Matches[j] = m.Matches[j], m.Matches[i]
}

func (m MatchResult) Len() int {
	return len(m.Matches)
}

func (m MatchResult) Paths() []string {
	res := make([]string, m.Len())
	for idx, m := range m.Matches {
		res[idx] = m.Path
	}
	return res
}

func (m *MatchResult) StripCommonPrefix() {
	if m.Len() <= 1 {
		return
	}
	prefix := findCommonPrefix(m.Paths())
	lastSlash := strings.LastIndex(prefix, "/")
	if lastSlash != -1 {
		prefix = prefix[0 : lastSlash+1]
	}
	for _, match := range m.Matches {
		match.Path = strings.TrimPrefix(match.Path, prefix)
	}
}

func findCommonPrefix(s []string) string {
	if len(s) == 0 {
		return ""
	}
	sort.Strings(s)
	i := 0
	for ; len(s[0]) > i && len(s[len(s)-1]) > i && s[0][i] == s[len(s)-1][i]; i++ {
	}
	if i == 0 {
		return ""
	}
	return s[0][:i]
}

func appendSlash(s string) string {
	if strings.HasSuffix(s, "/") {
		return s
	}
	return s + "/"
}

func entriesInDir(dir string) []string {
	var res []string
	fis, err := ioutil.ReadDir(dir)
	if err != nil {
		return nil
	}
	for _, fi := range fis {
		if fi.IsDir() {
			res = append(res, path.Join(dir, fi.Name()))
		}
	}
	return res
}

func (j Jumper) FindCandidates(jumpFile JumpFile, query string) MatchResult {
	var matches Matches
	// first, try to find an exact match
	if j.dirExists(query) && !matches.Has(appendSlash(query)) {
		abs, err := filepath.Abs(appendSlash(query))
		if err != nil {
			// TODO: what should we do here?
		}
		matches = append(matches, NewMatch(appendSlash(query), abs))
		if strings.HasSuffix(query, "/") {
			for _, entry := range entriesInDir(query) {
				abs, err := filepath.Abs(query + entry)
				if err != nil {
					continue
				}
				matches = append(matches, NewMatch(entry, abs))
			}
		}
	}

	// now, try to find a prefix match in the cwd
	prefixDirPath := j.prefixDir(query)
	if prefixDirPath != "" && !matches.Has(prefixDirPath+"/") {
		abs, err := filepath.Abs(prefixDirPath)
		if err != nil {
			// TODO: what should we do here?
		}
		matches = append(matches, NewMatch(prefixDirPath+"/", abs))
	}

	stringMatches := fuzzy.FindFrom(query, jumpFile)
	for _, match := range stringMatches {
		// only add the match if the entry exists in the filesystem and is a directory
		if j.dirExists(match.Str) && !matches.Has(match.Str) {
			matches = append(matches, &Match{
				Path:    match.Str,
				AbsPath: match.Str,
				Freq:    jumpFile[match.Index].Freq,
				Score:   match.Score,
			})
		}
	}
	res := MatchResult{Matches: matches, Query: query}
	sort.Sort(res)
	return res
}
