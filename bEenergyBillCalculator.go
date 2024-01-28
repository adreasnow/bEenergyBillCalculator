package main

import (
	"bufio"
	"flag"
	"fmt"
	"os"
	"slices"
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

type Bill struct {
	cost   float64
	supply float64
	days   int
}

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
					day.gas = 1

				} else if strings.Contains(line, "Anytime Usage,") {
					day.electricity = quantity
				}
				month[key] = day
				if !slices.Contains(dayList, key) {
					dayList = append(dayList, key)
				}

			}
		}
	}
	return month, nil
}

func CalculateCost(month Month, days int) (Bill, Bill, Bill) {
	totalElec := 0.0
	totalWater := 0.0
	supplyElec := 0.0
	supplyWater := 0.0
	supplyGas := 0.0

	dayCounter := 0

	var day Day
	var dayIndex int

	if days == 0 {
		days = len(dayList)
	}

	for i := 0; i <= days-1; i++ {
		dayIndex = (len(dayList) - days) + i
		day = month[dayList[dayIndex]]
		totalElec += (day.electricity * costs.electricity_rate)
		supplyElec += costs.electricity_supply
		totalWater += (day.water * costs.water_rate)
		supplyWater += costs.water_supply
		supplyGas += (float64(day.gas) * costs.gas_supply)
		dayCounter++

	}
	return Bill{cost: totalElec, supply: supplyElec, days: dayCounter},
		Bill{cost: totalWater, supply: supplyWater, days: dayCounter},
		Bill{cost: 0.0, supply: supplyGas, days: dayCounter}
}

func (b Bill) Avg(d int) float64 {
	oneDayCost := (b.cost / float64(b.days))
	oneDaySupply := (b.supply / float64(b.days))
	fmt.Println(oneDaySupply)
	return (oneDayCost + oneDaySupply) * float64(d)
}

func (b Bill) Total() float64 {
	return b.cost + b.supply
}

var costs Costs
var dayList []string

func main() {

	fmt.Printf("--------------------------------------------\n")
	fmt.Println("b.energy Bill Calculator Tool")

	dayList = make([]string, 0)
	month := make(Month, 0)

	electricity_rate := flag.Float64("er", 0.30894, "Electricity Rate")
	water_rate := flag.Float64("wr", 18.15, "Water Rate")
	gas_supply := flag.Float64("gs", 0.286, "Gas supply")
	electricity_supply := flag.Float64("es", 1.08661, "Electricity Supply")
	water_supply := flag.Float64("ws", 0.319, "Water Supply")
	flag.Parse()

	costs = Costs{electricity_rate: *electricity_rate,
		water_rate:         *water_rate,
		gas_supply:         *gas_supply,
		electricity_supply: *electricity_supply,
		water_supply:       *water_supply}

	fmt.Printf("--------------------------------------------\n")
	fmt.Printf("Utility\t\t Supply\t\t  Rate\n")
	fmt.Printf("--------------------------------------------\n")
	fmt.Printf("Electricity\t$%3.5f/day\t$%0.5f/KWh\n", costs.electricity_supply, costs.electricity_rate)
	fmt.Printf("Water\t\t$%3.3f/day\t$%2.2f/KL\n", costs.water_supply, costs.water_rate)
	fmt.Printf("Gas\t\t$%3.3f/day\t\n", costs.gas_supply)
	fmt.Printf("--------------------------------------------\n\n")
	var err error
	fmt.Println("Reading from files:")
	for _, arg := range flag.Args() {
		fmt.Printf("- %s\n", arg)
		month, err = ImportCSV(arg, month)
		if err != nil {
			fmt.Println(err.Error())
			os.Exit(1)
		}
	}
	fmt.Printf("\n")

	electricity, water, gas := CalculateCost(month, 0)
	electricity7d, water7d, gas7d := CalculateCost(month, 7)

	total := electricity.Total() + water.Total() + gas.Total()
	total30 := electricity.Avg(30) + water.Avg(30) + gas.Avg(30)
	total7 := electricity7d.Avg(30) + water7d.Avg(30) + gas7d.Avg(30)

	fmt.Printf("---------------------------------------------------\n")
	fmt.Printf("Utility:\tCost\t 30 day avg\t 30 day avg (last 7 days)\n")
	fmt.Printf("---------------------------------------------------\n")
	fmt.Printf("Electricity:\t$%3.2f\t   $%3.2f\t  $%3.2f\n", electricity.Total(), electricity.Avg(30), electricity7d.Avg(30))
	fmt.Printf("Water:\t\t$%3.2f\t   $%3.2f\t  $%3.2f\n", water.Total(), water.Avg(30), water7d.Avg(30))
	fmt.Printf("Gas:\t\t$%3.2f\t   $%3.2f\t  $%3.2f\n", gas.Total(), gas.Avg(30), gas7d.Avg(30))
	fmt.Printf("---------------------------------------------------\n")
	fmt.Printf("Total: \t\t$%3.2f\t   $%3.2f\t  $%3.2f\n", total, total30, total7)
	fmt.Printf("---------------------------------------------------\n")
}
