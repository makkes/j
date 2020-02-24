package j

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestMatching(t *testing.T) {
	tests := []struct {
		name     string
		jumpFile JumpFile
		query    string
		res      MatchResult
	}{
		{
			name: "two matching suffixes",
			jumpFile: JumpFile{
				JumpFileEntry{
					Path: "/Users/tester/dev/playground",
					Freq: 5,
				},
				JumpFileEntry{
					Path: "/Users/tester/dev/testing/playground",
					Freq: 25,
				},
			},
			query: "playground",
			res: MatchResult{
				Matches: []*Match{
					&Match{
						Path: "/Users/tester/dev/testing/playground",
					},
					&Match{
						Path: "/Users/tester/dev/playground",
					},
				},
			},
		},
		{
			name: "suffix match wins",
			jumpFile: JumpFile{
				JumpFileEntry{
					Path: "/Users/tester/dev/playground",
					Freq: 5,
				},
				JumpFileEntry{
					Path: "/Users/playground/whateva",
					Freq: 25,
				},
			},
			query: "playground",
			res: MatchResult{
				Matches: []*Match{
					&Match{
						Path: "/Users/tester/dev/playground",
					},
					&Match{
						Path: "/Users/playground/whateva",
					},
				},
			},
		},
		{
			name: "frequency wins",
			jumpFile: JumpFile{
				JumpFileEntry{
					Path: "/a/path/to/thissuper/dir",
					Freq: 5,
				},
				JumpFileEntry{
					Path: "/a/path/to/another/super/dir",
					Freq: 25,
				},
			},
			query: "super",
			res: MatchResult{
				Matches: []*Match{
					&Match{
						Path: "/a/path/to/another/super/dir",
					},
					&Match{
						Path: "/a/path/to/thissuper/dir",
					},
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			jumper := NewJumper(func(j *Jumper) {
				j.dirExists = func(p string) bool {
					return p != tt.query
				}
			})
			res := jumper.FindCandidates(tt.jumpFile, tt.query)
			require.Equal(t, tt.res.Len(), res.Len(), "Result length didn't match expected length")
			for i := 0; i < res.Len(); i++ {
				require.Equal(t, tt.res.Matches[i].Path, res.Matches[i].Path)
			}
		})
	}
}
