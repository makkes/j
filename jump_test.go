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
			name: "matches with same suffix are ranked by frequency",
			jumpFile: JumpFile{
				JumpFileEntry{
					Path: "/home/max/dev/weaveworks/weave-gitops",
					Freq: 55,
				},
				JumpFileEntry{
					Path: "/home/max/private/e13/infra/weave-gitops",
					Freq: 36,
				},
			},
			query: "weave-gitops",
			res: MatchResult{
				Matches: []*Match{
					{Path: "/home/max/dev/weaveworks/weave-gitops"},
					{Path: "/home/max/private/e13/infra/weave-gitops"},
				},
			},
		},
		{
			name: "better score always ranks higher",
			jumpFile: JumpFile{
				JumpFileEntry{
					Path: "/home/max/dev/weaveworks/cluster-bootstrap-controller",
					Freq: 8,
				},
				JumpFileEntry{
					Path: "/home/max/dev/weaveworks/cluster-controller",
					Freq: 5,
				},
			},
			query: "cluster-c",
			res: MatchResult{
				Matches: []*Match{
					{Path: "/home/max/dev/weaveworks/cluster-controller"},
					{Path: "/home/max/dev/weaveworks/cluster-bootstrap-controller"},
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
					{Path: "/Users/tester/dev/playground"},
					{Path: "/Users/playground/whateva"},
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
