package main

import (
	"fmt"
	"log"
	"os"
	"sort"
	"strconv"
	"strings"
)

func main() {
	//Variables declaration block
	var filePath string
	var err error

	//In the case the terminal's args len are 1 (default) we will take the default file path
	if len(os.Args) == 1 {
		filePath = "matches.txt"
	} else {
		filePath = os.Args[1]
	}

	//This function makes the analysis
	err = StartTournament(filePath)
	if err != nil {
		log.Fatalf("error when running tournament: %v\n", err)
	}
}

func StartTournament(filePath string) error {
	var err error

	var teamsTable map[string]int

	var bodyBytes []byte
	var matches []string

	teamsTable = make(map[string]int)

	//The file is read
	bodyBytes, err = os.ReadFile(filePath)
	if err != nil {
		return fmt.Errorf("error when reading file: %v", err)
	}

	//The file body's is splited by the newline character and we get the matches
	matches = strings.Split(string(bodyBytes), "\n")

	//Iterate over each mach and analalyze to determine the result
	for _, match := range matches {
		if len(match) == 0 {
			continue
		}
		teamsTable, err = AnalyzeMatch(match, teamsTable)
		if err != nil {
			return fmt.Errorf("file syntax error: %v", err)
		}
	}

	//Sort the table
	SortScores(teamsTable)

	return nil
}

// This function makes the operation to determine the winner/draw
func AnalyzeMatch(match string, teamsTable map[string]int) (map[string]int, error) {
	var err error
	var separator []string

	var lastIndex int

	var team1 string
	var team2 string

	var score1 int
	var score2 int

	//Convert "Team1 2, Team2 0" into ["Team1 2", " Team2 0"]
	separator = strings.Split(match, ",")

	//If the len of the separator is different from 2, it means that the format is not correct
	if len(separator) != 2 {
		return map[string]int{}, fmt.Errorf("format expected -> <str team1> <int score1>, <str team2> <int score1>\nGot: %v", match)
	}

	//Deletes unnecessary blankspaces, e.g: "   Team1  3  " -> "Team1 3"
	team1 = strings.Join(strings.Fields(separator[0]), " ")
	team2 = strings.Join(strings.Fields(separator[1]), " ")

	//Look for the last blankspace that is the one that separates the team from the scoring
	lastIndex = strings.LastIndex(team1, " ")
	if lastIndex == -1 {
		return map[string]int{}, fmt.Errorf("team format expected -> <str team1> <int score1>\nGot: %v", team1)
	}

	//Making a slicing operation we take the score value and assign it into a int variable validating it
	score1, err = strconv.Atoi(team1[lastIndex+1:])
	if err != nil {
		return map[string]int{}, fmt.Errorf("team format expected -> <str team1> <int score1>\nGot: %v", team1)
	}

	//The rest of the string is the team's name
	team1 = team1[:lastIndex]

	//Repeat the process with the second team
	lastIndex = strings.LastIndex(team2, " ")
	if lastIndex == -1 {
		return map[string]int{}, fmt.Errorf("team format expected -> <str team2> <int score2>\nGot: %v", team2)
	}

	score2, err = strconv.Atoi(team2[lastIndex+1:])
	if err != nil {
		return map[string]int{}, fmt.Errorf("team format expected -> <str team2> <int score2>\nGot: %v", team2)
	}

	team2 = team2[:lastIndex]

	//Validating scoring to determine wich team is the winner or if it was a draw
	if score1 > score2 {
		teamsTable[team1] += 3
		teamsTable[team2] += 0 //This might look redudant but it is written because it will add the team's name as a key in the map
	} else if score2 > score1 {
		teamsTable[team2] += 3
		teamsTable[team1] += 0
	} else {
		teamsTable[team1] += 1
		teamsTable[team2] += 1
	}

	return teamsTable, nil
}

// This function sorts the tournament's table
func SortScores(teamsTable map[string]int) {
	var teamsPositions []string = make([]string, len(teamsTable))
	var i int = 0

	//Assign the teams's names to a slice
	for key := range teamsTable {
		teamsPositions[i] = key
		i++
	}

	//This method sorts the slice by a private function validator
	sort.Slice(teamsPositions, func(i int, j int) bool {
		//In the case where two teams have the same points, the next undraw method is the alphabetic order
		if teamsTable[teamsPositions[i]] == teamsTable[teamsPositions[j]] {
			return teamsPositions[i] < teamsPositions[j] //Determine alphabetic order
		}

		return teamsTable[teamsPositions[i]] > teamsTable[teamsPositions[j]] //Determine which team has more points
	})

	//The table is printed
	i = 1
	for j := 0; j < len(teamsPositions); j++ {
		fmt.Printf("%v. %v, %v pts\n", i, teamsPositions[j], teamsTable[teamsPositions[j]])

		//This condition is for assigning the same position to the teams that have the same points
		if j < len(teamsPositions)-1 {
			if teamsTable[teamsPositions[j+1]] < teamsTable[teamsPositions[j]] {
				i = j + 2
			}
		}
	}
}
