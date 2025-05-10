package rank

import (
	"bufio"
	"fmt"
	"os"
	"sort"
	"strconv"
	"strings"
)

type TeamScore struct {
	TeamName string
	Score    int
}

type Match struct {
	Team1 TeamScore
	Team2 TeamScore
}

type LeagueEntry struct {
	TeamName string
	Points   int
	Rank     int
}

const (
	win  = 3
	tie  = 1
	loss = 0
)

type RankedLeague []LeagueEntry
type LeaguePoints map[string]int

func Run() {
	teams := make(map[string]int)

	scanner := bufio.NewScanner(os.Stdin)
	line := 1

	for scanner.Scan() {
		match, err := parseMatch(scanner.Text(), line)
		if err != nil {
			if errors.Is(err, errEmptyLine) {
				// empty line, do nothing, lenient on whitespace
				continue
			}
			fmt.Fprintf(os.Stderr, "error: %v\n", err)
			os.Exit(1)
		}
		err = updateLeague(teams, match)
		if err != nil {
			fmt.Fprintf(os.Stderr, "error: %v\n", err)
			os.Exit(1)
		}
		line++
	}

	ss := rankedLeague(teams)

	printRankedLeague(ss)
}

func printRankedLeague(league RankedLeague) {
	for _, entry := range league {
		if entry.Points == 1 {
			fmt.Printf("%d. %s, %d pt\n", entry.Rank, entry.TeamName, entry.Points)
			continue
		}
		fmt.Printf("%d. %s, %d pts\n", entry.Rank, entry.TeamName, entry.Points)
	}
}

func rankedLeague(league LeaguePoints) RankedLeague {
	if len(league) == 0 {
		return nil // do nothing
	}

	ss := make(RankedLeague, 0, len(league))
	for teamName, points := range league {
		ss = append(ss, LeagueEntry{
			TeamName: teamName,
			Points:   points,
		})
	}

	sort.Slice(ss, func(i, j int) bool {
		if ss[i].Points != ss[j].Points {
			return ss[i].Points > ss[j].Points
		}
		return ss[i].TeamName < ss[j].TeamName
	})

	currentRank := 1

	// handle the first case separately to keep the loop cleaner
	//fmt.Printf("%d. %s, %d pts\n", currentRank, ss[0].TeamName, ss[0].Points)

	ss[0].Rank = currentRank

	// rest of the teams
	for i := 1; i < len(ss); i++ {
		if ss[i].Points == ss[i-1].Points {
			// equal points, same rank
			ss[i].Rank = currentRank
		} else {
			// if not in the same rank, use alice idx i + 1 as the rank (counting from 1)
			currentRank = i + 1
			ss[i].Rank = currentRank
		}
	}
	return ss
}

var errEmptyLine error = fmt.Errorf("empty line")

func parseMatch(s string, lineNumber int) (Match, error) {
	if len(s) < 1 { // empty line, allow, does nothing
		return Match{}, errEmptyLine
	}
	parts := strings.Split(s, ",")
	if len(parts) != 2 {
		return Match{}, fmt.Errorf("line: %d, invalid match format: %s", lineNumber, s)
	}
	team1Info := parts[0]
	team2Info := parts[1]

	team1Name, team1Score, err := splitTeamInfo(team1Info, lineNumber)
	if err != nil {
		return Match{}, err
	}
	team2Name, team2Score, err := splitTeamInfo(team2Info, lineNumber)
	if err != nil {
		return Match{}, err
	}
	return Match{
		Team1: TeamScore{
			TeamName: team1Name,
			Score:    team1Score,
		},
		Team2: TeamScore{
			TeamName: team2Name,
			Score:    team2Score,
		},
	}, nil
}

func splitTeamInfo(teamInfo string, lineNumber int) (string, int, error) {
	// len has to be at least 1 char per name, comma, and 1 digit score each.
	if len(teamInfo) < 5 {
		return "", 0, fmt.Errorf("line: %d, invalid format: %s", lineNumber, teamInfo)
	}

	teamInfo = strings.TrimSpace(teamInfo)

	idxLastSpace := strings.LastIndex(teamInfo, " ")

	teamName := strings.TrimSpace(teamInfo[:idxLastSpace])

	if teamName == "," {
		return "", 0, fmt.Errorf("line: %d, invalid empty team name: %s", lineNumber, teamName)
	}

	scoreStr := strings.TrimSpace(teamInfo[idxLastSpace:])

	score, err := strconv.Atoi(scoreStr)
	if err != nil {
		return "", 0, err
	}
	if score < 0 {
		return "", 0, fmt.Errorf("line: %d, invalid negative score: %d", lineNumber, score)
	}
	if len(teamName) < 1 {
		return "", 0, fmt.Errorf("line: %d, invalid empty team name: %s", lineNumber, teamName)
	}
	return teamName, score, nil
}

func updateLeague(league LeaguePoints, match Match) error {
	if league == nil {
		return fmt.Errorf("league is nil")
	}

	if match.Team1.Score > match.Team2.Score {
		league[match.Team1.TeamName] += win
		league[match.Team2.TeamName] += loss
	}
	if match.Team1.Score < match.Team2.Score {
		league[match.Team2.TeamName] += win
		league[match.Team1.TeamName] += loss
	}
	if match.Team1.Score == match.Team2.Score {
		league[match.Team1.TeamName] += tie
		league[match.Team2.TeamName] += tie
	}
	return nil
}
