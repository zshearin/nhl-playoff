package main

import (
	"fmt"
	"math/rand"
	"os"
	"sort"
	"strconv"
	"strings"
	"time"
)

//Matchup denotes a matchup
type Matchup struct {
	LowerSeed  int
	HigherSeed int
}

func main() {

	rand.Seed(time.Now().UnixNano())

	args := os.Args[1:]
	if len(args) > 0 {
		if args[0] == "s" {
			printSeriesLengthPredictions()
			return
		}
	}

	runSimulation()
}



func runSimulation() {
        conferences := []string{"East", "West"}

        var finals []int
        for _, confName := range conferences {
                conferenceSeedWinner := runConferenceSimulation(confName)

                finals = append(finals, conferenceSeedWinner)

        }

        eastWins := FindWinner(finals[0], finals[1])
        if eastWins {
                fmt.Println("East wins, winner seed: " + strconv.Itoa(finals[0]))
        } else {
                fmt.Println("West wins, winner seed: " + strconv.Itoa(finals[1]))
        }

}



/*
Percentage of Best of 7 series that went different lengths:
4		5		6		7
---------------------------------------------------------
0.178378378	0.255405405	0.316216216	0.25
cumulative sum: .433784		.75		1 
*/

func getNumberOfGames(randNum float64) int {

	if randNum < .178378 {
		return 4
	} else if randNum < .433784 {
		return 5
	} else if randNum < .75 {
		return 6
	} else {
		return 7
	}

}


func printSeriesLengthPredictions() {
        var confSeeds []int
        for i := 0; i < 8; i++ {
                seed := i + 1
                confSeeds = append(confSeeds, seed)

        }

        r1m := GetMatchups(confSeeds)
        PrintRoundMatchups(r1m, 1)
        for i := 0; i < len(r1m); i++ {
                num := getNumberOfGames(getRandomNumber())
                fmt.Println("Matchup " + strconv.Itoa(i+1)+ " number of games: " + strconv.Itoa(num))
        }
}



//FindWinner returns true if eastSeed wins, false if westSeed wins
func FindWinner(eastSeed, westSeed int) bool {

	rand.Seed(time.Now().UnixNano())
	randNumber := getRandomNumber()

	return randNumber <= .5
}

func runConferenceSimulation(region string) int {
	fmt.Println("==========================================================")
	underline := strings.Repeat("-", 16+len(region))
	fmt.Printf("%s region results: \n", region)
	fmt.Println(underline)

	var confSeeds []int
	for i := 0; i < 8; i++ {
		seed := i + 1
		confSeeds = append(confSeeds, seed)

	}

	r1m := GetMatchups(confSeeds)
	PrintRoundMatchups(r1m, 1)
	for i := 0; i < len(r1m); i++ {
		num := getNumberOfGames(getRandomNumber())
		fmt.Println("Matchup " + strconv.Itoa(i+1)+ " number of games: " + strconv.Itoa(num))
	}

	r1r := GetMatchupResults(r1m)

	//return r1r
	sort.Ints(r1r)

	r2m := GetMatchups(r1r)
	PrintRoundMatchups(r2m, 2)

	r2r := GetMatchupResults(r2m)
	r3m := GetMatchups(r2r)
	PrintRoundMatchups(r3m, 3)

	r3r := GetMatchupResults(r3m)
	//r4m := GetMatchups(r3r)

	//r4r := GetMatchupResults(r4m)
	winner := r3r[0]
	fmt.Println("Conference winner: " + strconv.Itoa(winner))
	return winner
}

//PrintRoundMatchups prints the matchups for a round
func PrintRoundMatchups(matchups []Matchup, roundNumber int) {

	fmt.Printf("Round %d matchups: ", roundNumber)

	for _, value := range matchups {
		fmt.Printf("%d vs. %d   ", value.LowerSeed, value.HigherSeed)
	}
	fmt.Printf("\n")
	//	fmt.Println(matchups)

}

//GetMatchupResults gets the result array from the provided matchup
func GetMatchupResults(matchups []Matchup) []int {

	var results []int

	for _, matchup := range matchups {

		lowerSeed := matchup.LowerSeed
		higherSeed := matchup.HigherSeed

		result := notAnUpset(lowerSeed, higherSeed)

		if result {
			results = append(results, lowerSeed)
		} else {
			results = append(results, higherSeed)
		}
	}

	return results

}

//GetMatchups gets the matchups from the int array
func GetMatchups(round1Results []int) []Matchup {

	//
	var matchups []Matchup

	for len(round1Results) > 1 {

		var first, last int

		first, round1Results = round1Results[0], round1Results[1:]

		last, round1Results = round1Results[len(round1Results)-1], round1Results[:len(round1Results)-1]

		var curMatchup Matchup
		if first < last {
			curMatchup.LowerSeed = first
			curMatchup.HigherSeed = last
		} else {
			curMatchup.LowerSeed = last
			curMatchup.HigherSeed = first
		}

		matchups = append(matchups, curMatchup)
	}

	return matchups

}

//returns true if results expected (not an upset)
func notAnUpset(num1, num2 int) bool {

	//make num1 the smaller number and num2 the bigger number (seeding processing later)
	if num1 > num2 {
		temp := num1
		num1 = num2
		num2 = temp
	}

	//Generate threshold and random number to assess against (probability based on seed value)
	threshold := calculateThreshold(num1, num2)
	random := getRandomNumber()
	return random < threshold

}

func getRandomNumber() float64 {
	random := rand.Float64()
	return random
}

//Assess random value relative to threshold
// For a 1 and 16 seed matchup, the 1 seed has a 16/17 chance
// of winning and the 16 seed has a 1/17 chance
func calculateThreshold(num1, num2 int) float64 {
	sum := num1 + num2
	if num1 > num2 {
		num2 = num1
	}
	threshold := float64(num2) / float64(sum)
	return threshold
}
