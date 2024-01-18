package main

import (
	"bufio"
	"flag"
	"fmt"
	"os"
	"strconv"
	"strings"
)

type Day struct {
	water       float64
	electricity float64
	gas         int
}

type Costs struct {
	electricity_rate   float64
	water_rate         float64
	gas_supply         float64
	electricity_supply float64
	water_supply       float64
}

type Month map[string]Day
type Bill float64

func ReadStrings(fileName string) ([]string, error) {
	file, err := os.Open(fileName)
	if err != nil {
		return nil, err
	}

	defer file.Close()

	input := make([]string, 0)
	scanner := bufio.NewScanner(file)

	for scanner.Scan() {
		input = append(input, scanner.Text())
	}

	if err := scanner.Err(); err != nil {
		return nil, err
	}
	return input, err
}

func ImportCSV(csv string, month Month) (Month, error) {
	lines, err := ReadStrings(csv)
	if err != nil {
		return month, err
	} else {
		var key string
		var quantity float64
		var day Day
		var splitLine []string

		for _, line := range lines {
			if strings.Contains(line, "Hot Water,") || strings.Contains(line, "Anytime Usage,") {
				splitLine = strings.Split(line, ",")
				key = splitLine[1]
				quantity, err = strconv.ParseFloat(splitLine[2], 64)
				if err != nil {
					return month, err
				}

				if _, ok := month[key]; !ok {
					day = Day{}
				} else {
					day = month[key]
				}

				if strings.Contains(line, "Hot Water,") {
					day.water = quantity
					if day.water > 0 {
						day.gas = 1
					}

				} else if strings.Contains(line, "Anytime Usage,") {
					day.electricity = quantity
					if day.water > 0 {
						day.gas = 1
					}
				}
				month[key] = day
				dayList = append(dayList, key)
			}
		}
	}
	return month, nil
}

func CalculateCost(month Month) (Bill, Bill, Bill, int) {
	totalElec := 0.0
	totalWater := 0.0
	totalGas := 0.0
	days := 0

	for _, day := range month {
		if day.electricity != 0.0 {
			totalElec += (day.electricity * costs.electricity_rate) + costs.electricity_supply
		}
		if day.water != 0.0 {
			totalWater += (day.water * costs.water_rate) + costs.water_supply
		}
		if day.gas != 0 {
			totalGas += (float64(day.gas) * costs.gas_supply)
		}
		if day.gas == 1 {
			days++
		}

	}
	return Bill(totalElec), Bill(totalWater), Bill(totalGas), days
}

func CalculateCostDays(month Month, days int) (Bill, Bill, Bill, int) {
	totalElec := 0.0
	totalWater := 0.0
	totalGas := 0.0
	var day Day
	var dayIndex int

	for i := 0; i <= days; i++ {
		dayIndex = (len(dayList) - 1 - days) + i
		day = month[dayList[dayIndex]]
		if day.electricity != 0.0 {
			totalElec += (day.electricity * costs.electricity_rate) + costs.electricity_supply
		}
		if day.water != 0.0 {
			totalWater += (day.water * costs.water_rate) + costs.water_supply
		}
		if day.gas != 0 {
			totalGas += (float64(day.gas) * costs.gas_supply)
		}
	}

	return Bill(totalElec), Bill(totalWater), Bill(totalGas), days
}

func (b Bill) _30DayAvg(d int) float64 {

	return (float64(b) / float64(d)) * 30.0
}

var costs Costs
var dayList []string

func main() {
	fmt.Println("-------------------------------")
	fmt.Println("b.energy Bill Calculator Tool")
	fmt.Println("-------------------------------\n")
	costs = Costs{}
	dayList = make([]string, 0)
	month := make(Month, 0)

	costs.electricity_rate = *flag.Float64("er", 0.30894, "er")
	costs.water_rate = *flag.Float64("wr", 18.15, "wr")
	costs.gas_supply = *flag.Float64("gs", 0.286, "gs")
	costs.electricity_supply = *flag.Float64("es", 1.08661, "es")
	costs.water_supply = *flag.Float64("ws", 0.319, "ws")
	flag.Parse()

	var err error
	fmt.Println("Reading from files:")
	for _, arg := range flag.Args() {
		fmt.Println(arg)
		month, err = ImportCSV(arg, month)
		if err != nil {
			fmt.Println(err.Error())
			os.Exit(1)
		}
	}
	fmt.Printf("\n")

	electricity, water, gas, days := CalculateCost(month)
	electricity7d, water7d, gas7d, _ := CalculateCostDays(month, 7)

	total := electricity + water + gas
	_30dayTotal := electricity._30DayAvg(days) + water._30DayAvg(days) + gas._30DayAvg(days)
	_30dayTotal7d := electricity7d._30DayAvg(7) + water7d._30DayAvg(7) + gas7d._30DayAvg(7)

	fmt.Printf("---------------------------------------------------\n")
	fmt.Printf("Utility:\tCost\t 30 day avg\t 7 day avg\n")
	fmt.Printf("---------------------------------------------------\n")
	fmt.Printf("Electricity:\t$%3.2f\t   $%3.2f\t  $%3.2f\n", electricity, electricity._30DayAvg(days), electricity7d._30DayAvg(7))
	fmt.Printf("Water:\t\t$%3.2f\t   $%3.2f\t  $%3.2f\n", water, water._30DayAvg(days), water7d._30DayAvg(7))
	fmt.Printf("Gas:\t\t$%3.2f\t   $%3.2f\t  $%3.2f\n", gas, gas._30DayAvg(days), gas7d._30DayAvg(7))
	fmt.Printf("---------------------------------------------------\n")
	fmt.Printf("Total: \t\t$%3.2f\t   $%3.2f\t  $%3.2f\n", total, _30dayTotal, _30dayTotal7d)
	fmt.Printf("---------------------------------------------------\n")
}
