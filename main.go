package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"strconv"
	"strings"
	"time"
)

type rawTransaction struct {
	Owner       string   `json:"owner"`
	BlockHeight uint64   `json:"blockHeight"`
	BlockHash   string   `json:"blockHash"`
	BlockTime   uint64   `json:"blockTime"`
	Type        string   `json:"type"`
	RewardType  string   `json:"rewardType"`
	PoolID      uint64   `json:"poolID"`
	Amounts     []string `json:"amounts"`
}

type reward struct {
	BlockTime uint64 `json:"blockTime"`
	Amounts   map[string]float64
}

type dateReward struct {
	Date    time.Time `json:"date"`
	Amounts map[string]float64
}

func getTransactionsByBlockTime(t []rawTransaction) []reward {
	// create rewards map to store results
	rewards := make(map[uint64]map[string]float64)

	for i := range t {
		if t[i].Type == "Rewards" || t[i].Type == "Commission" || t[i].Type == "receive" {
			// split amount and currency from string, also only check for first amount (since there are never more)
			amountCurrency := strings.Split(t[i].Amounts[0], "@")
			currency := amountCurrency[1]
			amount, err := strconv.ParseFloat(amountCurrency[0], 64)
			if err != nil {
				panic(err)
			}

			// check if there are no entries for this Blocktime and create it if necessary
			if len(rewards[t[i].BlockTime]) == 0 {
				rewards[t[i].BlockTime] = make(map[string]float64)
			}

			// check if there are entries for this currency for this blocktime and add if already exists
			if rewards[t[i].BlockTime][currency] != 0 {
				rewards[t[i].BlockTime][currency] += amount
				continue
			}

			// set equal to amount if no entry exists
			rewards[t[i].BlockTime][currency] = amount
		}
	}

	// create struct array to hold final results
	var r []reward

	// loop through rewards and assemble into struct array
	for i := range rewards {
		r = append(r, reward{
			BlockTime: i,
			Amounts:   rewards[i],
		})
	}

	// return rewards
	return r
}

func getTransactionsByDay(r []reward) {

	// create rewards map to store results
	dailyrewards := make(map[string]map[string]float64)

	for i := range r {
		tm := time.Unix(int64(r[i].BlockTime), 0)
		day := tm.Format("2006-01-02")

		// check if there are no entries for this day and create it if necessary
		if len(dailyrewards[day]) == 0 {
			dailyrewards[day] = make(map[string]float64)
		}

		// add all currency amounts to its correct map
		for currency, amount := range r[i].Amounts {
			dailyrewards[day][currency] += amount
		}
	}

	fmt.Println(dailyrewards["2022-03-10"])

	// prettify.PrintArray(dailyrewards)
	// fmt.Println(len(dailyrewards))
	// fmt.Println(dailyrewards["2022-03-10"])
}

func main() {
	var t []rawTransaction

	// Open our jsonFile
	jsonFile, err := os.Open("transactions.json")

	// if we os.Open returns an error then handle it
	if err != nil {
		fmt.Println(err)
	}

	// defer the closing of our jsonFile so that we can parse it later on
	defer jsonFile.Close()

	// read our opened jsonFile as a byte array.
	byteValue, _ := ioutil.ReadAll(jsonFile)
	json.Unmarshal(byteValue, &t)

	fmt.Printf("Successfully read %d transactions.\n", len(t))

	btt := getTransactionsByBlockTime(t)

	// prettify.PrintArray(btt)
	getTransactionsByDay(btt)
	// getRewards(t)
}
