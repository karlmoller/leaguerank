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

type RankedLeague []LeagueEntry
type LeaguePoints map[string]int

func Run() {
	teams := make(map[string]int)

	scanner := bufio.NewScanner(os.Stdin)
	for scanner.Scan() {
		match, err := parseMatch(scanner.Text())
		if err != nil {
			fmt.Fprintf(os.Stderr, "error: %v\n", err)
		}
		updateLeague(teams, match)
	}

	ss := rankedLeague(teams)

	for _, entry := range ss {
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

func parseMatch(s string) (Match, error) {
	if len(s) < 1 { // empty line
		return Match{}, nil
	}
	parts := strings.Split(s, ",")
	if len(parts) != 2 {
		panic("invalid input, expected 2 parts")
	}
	team1Info := parts[0]
	team2Info := parts[1]

	team1Name, team1Score, err := splitTeamInfo(team1Info)
	if err != nil {
		return Match{}, err
	}
	team2Name, team2Score, err := splitTeamInfo(team2Info)
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

func splitTeamInfo(teamInfo string) (string, int, error) {
	idxLastSpace := strings.LastIndex(teamInfo, " ")

	teamName := strings.TrimSpace(teamInfo[:idxLastSpace])

	scoreStr := strings.TrimSpace(teamInfo[idxLastSpace:])

	score, err := strconv.Atoi(scoreStr)
	if err != nil {
		return "", 0, err
	}
	if score < 0 {
		return "", 0, fmt.Errorf("invalid negative score: %d", score)
	}
	return teamName, score, nil
}

func updateLeague(league LeaguePoints, match Match) {
	if league == nil {
		league = make(LeaguePoints)
	}

	if match.Team1.Score > match.Team2.Score {
		league[match.Team1.TeamName] += 3
		league[match.Team2.TeamName] += 0
	}
	if match.Team1.Score < match.Team2.Score {
		league[match.Team2.TeamName] += 3
		league[match.Team1.TeamName] += 0
	}
	if match.Team1.Score == match.Team2.Score {
		league[match.Team1.TeamName] += 1
		league[match.Team2.TeamName] += 1
	}
}
