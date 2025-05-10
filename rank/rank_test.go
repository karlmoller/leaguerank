package rank

import "testing"

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
			wantErr: true,
		},
		{
			name: "Ending with space",
			args: args{
				teamInfo:   "Some Team Name 1 ",
				lineNumber: 1,
			},
			wantErr: true,
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
