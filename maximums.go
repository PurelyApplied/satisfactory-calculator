package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"strings"
)

const (
	powerKey = "power"
	shardKey = "overclock-shard"
)

type belt struct {
	Name string
	Key  string `json:"key_name"`
	Rate float64
}

type building struct {
	Name     string
	Key      string `json:"key_name"`
	Category string
	Power    int
	Inputs   int `json:"max"`
}

type miner struct {
	Name     string
	Key      string `json:"key_name"`
	Category string
	Rate     float64 `json:"base_rate"`
	Power    int
}

type item struct {
	Name  string
	Key   string `json:"key_name"`
	Tier  int
	Stack int `json:"stack_size"`
}

func (i item) String() string {
	return fmt.Sprintf("<%q (%s)>", i.Name, i.Key)
}

// products and ingredients coded as: key, count
type recipe struct {
	Name        string
	Key         string `json:"key_name"`
	Category    string
	Time        float64
	Ingredients [][2]interface{}
	Product     [2]interface{}
}

type resource struct {
	Key      string `json:"key_name"`
	Category string
}

type data struct {
	Belts     []belt
	Buildings []building
	Miners    []miner
	Items     []item
	Recipes   []recipe
	Resources []resource
}

func main() {
	do(false)
}

func unpack(i [2]interface{}) (string, float64) {
	return i[0].(string), i[1].(float64)
}

func do(includeAlts bool) {
	dat := loadData()
	index := buildIndex(dat)

	// Build matrix from recipes
	var matrix [][]float64
	width := len(index)

	for _, r := range dat.Recipes {
		if includeAlts && strings.HasPrefix(r.Name, "Alternate") {
			continue
		}

		row := make([]float64, width)
		// Normalize to per minute
		norm := 60.0 / r.Time
		for _, ing := range r.Ingredients {
			key, count := unpack(ing)

			row[index[key]] = -norm * count
		}

		prodKey, prodCount := unpack(r.Product)

		row[index[prodKey]] = norm * prodCount

		row[index[powerKey]] = -getPowerCost(r, dat)

		matrix = append(matrix, row)
	}

	for _, x := range matrix {
		for _, v := range x {
			fmt.Printf("%7.2f ", v)
		}
		fmt.Println()
	}

}

func getPowerCost(r recipe, dat data) float64 {
	cat := r.Category
	for _, b := range dat.Buildings {
		if cat == b.Category {
			return float64(b.Power)
		}
	}
	panic("Failed to find building for power cost.")
}

func buildIndex(dat data) map[string]int {
	// Matrix will include resources, items, and custom columns for overclocking resources and power
	i := 0
	index := make(map[string]int)

	for _, reso := range dat.Resources {
		if _, ok := index[reso.Key]; !ok {
			index[reso.Key] = i
			i++
		}
	}

	for _, item := range dat.Items {
		if _, ok := index[item.Key]; !ok {
			index[item.Key] = i
			i++
		}
	}

	index[powerKey] = i
	i++
	index[shardKey] = i

	return index
}

func loadData() data {
	bytes, err := ioutil.ReadFile("data/data.json")
	if err != nil {
		log.Fatal(err)
	} else {
		fmt.Println("Loaded: ", string(bytes))
	}
	var dat data
	err = json.Unmarshal(bytes, &dat)
	if err != nil {
		log.Fatal(err)
	}
	return dat
}
