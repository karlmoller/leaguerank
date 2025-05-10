package rank

import (
	"reflect"
	"testing"
)

func Test_splitTeamInfo(t *testing.T) {
	type args struct {
		teamInfo   string
		lineNumber int
	}
	tests := []struct {
		name      string
		args      args
		wantName  string
		wantScore int
		wantErr   bool
	}{
		{
			name: "Long Team Name",
			args: args{
				teamInfo:   "A Very Long Team Name.. Formerly Known as FC Awesome 5",
				lineNumber: 1,
			},
			wantName:  "A Very Long Team Name.. Formerly Known as FC Awesome",
			wantScore: 5,
			wantErr:   false,
		},
		{
			name: "Missing Score",
			args: args{
				teamInfo:   "Some Team Name",
				lineNumber: 1,
			},
			wantErr: true,
		},
		{
			name: "Negative Score",
			args: args{
				teamInfo:   "Some Team Name, -1",
				lineNumber: 1,
			},
			wantErr: true,
		},
		{
			name: "Empty Team Name",
			args: args{
				teamInfo:   ", 1",
				lineNumber: 1,
			},
			wantErr: true,
		},
		{
			name: "Starting with Space",
			args: args{
				teamInfo:   "     Some Team Name 1",
				lineNumber: 1,
			},
			wantErr:   false,
			wantName:  "Some Team Name",
			wantScore: 1,
		},
		{
			name: "Ending with space",
			args: args{
				teamInfo:   "Some Team Name 1 ",
				lineNumber: 1,
			},
			wantErr:   false,
			wantName:  "Some Team Name",
			wantScore: 1,
		},
		{
			name: "Ending with comma",
			args: args{
				teamInfo:   "Some Team Name, 1, ",
				lineNumber: 1,
			},
			wantErr: true,
		},
		{
			name: "Empty line",
			args: args{
				teamInfo: "",
			},
			wantErr: true,
		},
		{
			name: "invalid score, alpha",
			args: args{
				teamInfo:   "Some Team Name, a",
				lineNumber: 1,
			},
			wantErr: true,
		},
		{
			name: "invalid score, float",
			args: args{
				teamInfo:   "Some Team Name, 0.5",
				lineNumber: 1,
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, got1, err := splitTeamInfo(tt.args.teamInfo, tt.args.lineNumber)
			if (err != nil) != tt.wantErr {
				t.Errorf("splitTeamInfo() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.wantName {
				t.Errorf("splitTeamInfo() got = %v, want %v", got, tt.wantName)
			}
			if got1 != tt.wantScore {
				t.Errorf("splitTeamInfo() got1 = %v, want %v", got1, tt.wantScore)
			}
		})
	}
}

func Test_parseMatch(t *testing.T) {
	type args struct {
		s          string
		lineNumber int
	}
	tests := []struct {
		name    string
		args    args
		want    Match
		wantErr bool
	}{
		{
			name: "Valid Match",
			args: args{
				s:          "Team A 1, Team B 2",
				lineNumber: 1,
			},
			want: Match{
				Team1: TeamScore{
					TeamName: "Team A",
					Score:    1,
				},
				Team2: TeamScore{
					TeamName: "Team B",
					Score:    2,
				},
			},
			wantErr: false,
		},
		{
			name: "Invalid Match -- no Score Team B",
			args: args{
				s:          "Team A 1, Team B",
				lineNumber: 1,
			},
			wantErr: true,
		},
		{
			name: "Missing team",
			args: args{
				s:          "Foo 1,",
				lineNumber: 1,
			},
			wantErr: true,
		},
		{
			name: "Empty line",
			args: args{
				s: "",
			},
			wantErr: false, // empty line is valid, does nothing
		},
		{
			name: "Missing team -- first part",
			args: args{
				s: ", Bar 5",
			},
			wantErr: true,
		},
		{
			name: "Too many parts",
			args: args{
				s: "Foo 1, Bar 2, Baz 3",
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := parseMatch(tt.args.s, tt.args.lineNumber)
			if (err != nil) != tt.wantErr {
				t.Errorf("parseMatch() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("parseMatch() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_updateLeague(t *testing.T) {
	league := make(LeaguePoints)

	match := Match{
		Team1: TeamScore{
			TeamName: "Team A",
			Score:    1,
		},
		Team2: TeamScore{
			TeamName: "Team B",
			Score:    2,
		},
	}
	err := updateLeague(league, match)
	if err != nil {
		t.Errorf("error: %v", err)
	}

	if league == nil {
		t.Errorf("league is nil")
	}
	if league["Team A"] != 0 {
		t.Errorf("league[Team A] = %d, want 0", league["Team A"])
	}
	if league["Team B"] != 3 {
		t.Errorf("league[Team B] = %d, want 0", league["Team B"])
	}

	match2 := Match{
		Team1: TeamScore{
			TeamName: "Team A",
			Score:    2,
		},
		Team2: TeamScore{
			TeamName: "Team B",
			Score:    1,
		},
	}

	err = updateLeague(league, match2)
	if err != nil {
		t.Errorf("error: %v", err)
	}
	if league["Team A"] != 3 {
		t.Errorf("league[Team A] = %d, want 3", league["Team A"])
	}
	if league["Team B"] != 3 {
		t.Errorf("league[Team B] = %d, want 3", league["Team B"])
	}

	matchTied := Match{
		Team1: TeamScore{
			TeamName: "Team A",
			Score:    2,
		},
		Team2: TeamScore{
			TeamName: "Team B",
			Score:    2,
		},
	}

	err = updateLeague(league, matchTied)
	if err != nil {
		t.Errorf("error: %v", err)
	}
	if league["Team A"] != 4 {
		t.Errorf("league[Team A] = %d, want 4", league["Team A"])
	}
	if league["Team B"] != 4 {
		t.Errorf("league[Team B] = %d, want 4", league["Team B"])
	}
}
