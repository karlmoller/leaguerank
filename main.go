package main

import (
	"bufio"
	"fmt"
	"os"
	"sort"
	"strconv"
	"strings"
)

type Score struct {
	TeamName string
	Points   int
}

type Match struct {
	Team1 Score
	Team2 Score
}

type League map[string]int

func main() {
	teams := make(map[string]int)

	scanner := bufio.NewScanner(os.Stdin)
	for scanner.Scan() {
		match, err := parseMatch(scanner.Text())
		if err != nil {
			fmt.Fprintf(os.Stderr, "error: %v\n", err)
		}
		updateLeague(teams, match)
	}
	ss := make([]Score, 0, len(teams))
	for teamName, points := range teams {
		ss = append(ss, Score{
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

	if len(ss) < 1 {
		os.Exit(0)
	}

	currentRank := 1

	// handle the first case separately to keep the loop cleaner
	fmt.Printf("%d. %s %d pts\n", currentRank, ss[0].TeamName, ss[0].Points)

	// rest of the teams
	for i := 1; i < len(ss); i++ {
		if ss[i].Points == ss[i-1].Points {
			// equal points, same rank
			fmt.Printf("%d. %s %d pts\n", currentRank, ss[i].TeamName, ss[i].Points)
		} else {
			// if not in the same rank, use alice idx i + 1 as the rank (counting from 1)
			currentRank = i + 1
			fmt.Printf("%d. %s %d pts\n", currentRank, ss[i].TeamName, ss[i].Points)
		}
	}
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
		Team1: Score{
			TeamName: team1Name,
			Points:   team1Score,
		},
		Team2: Score{
			TeamName: team2Name,
			Points:   team2Score,
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

func updateLeague(league League, match Match) {
	if match.Team1.Points > match.Team2.Points {
		league[match.Team1.TeamName] += 3
		league[match.Team2.TeamName] += 0
	}
	if match.Team1.Points < match.Team2.Points {
		league[match.Team2.TeamName] += 3
		league[match.Team1.TeamName] += 0
	}
	if match.Team1.Points == match.Team2.Points {
		league[match.Team1.TeamName] += 1
		league[match.Team2.TeamName] += 1
	}
}
