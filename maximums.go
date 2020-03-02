package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"strings"
)

type belt struct {
	Name string `json:"name"`
	Key  string `json:"key_name"`
	Rate int    `json:"rate"`
}

type building struct {
	Name     string `json:"name"`
	Key      string `json:"key_name"`
	Category string `json:"category"`
	Power    int    `json:"power"`
	Max      int    `json:"max"`
}

type miner struct {
	Name     string `json:"name"`
	Key      string `json:"key_name"`
	Category string `json:"category"`
	Rate     int    `json:"base_rate"`
	Power    int    `json:"power"`
}

type item struct {
	Name  string `json:"name"`
	Key   string `json:"key_name"`
	Tier  int    `json:"tier"`
	Stack int    `json:"stack_size"`
}

func (i item) String() string {
	return fmt.Sprintf("<%q (%s)>", i.Name, i.Key)
}

// products and ingredients coded as: key, count
type recipe struct {
	Name        string           `json:"name"`
	Key         string           `json:"key_name"`
	Category    string           `json:"category"`
	Time        int              `json:"time"`
	Ingredients [][2]interface{} `json:"ingredients"`
	Product     [2]interface{}   `json:"product"`
}

type data struct {
	Belts     []belt     `json:"belts"`
	Buildings []building `json:"buildings"`
	Miners    []miner    `json:"miners"`
	Items     []item     `json:"items"`
	Recipes   []recipe   `json:"recipes"`
}

const includeAlts = false

func main() {
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

	// Reverse-map indices for items.
	index := make(map[string]int)
	for i, item := range dat.Items {
		index[item.Key] = i
	}

	// Build matrix from recipes
	var recipeConversions [][]float64
	var recipesUsed []recipe

	width := len(dat.Items)
	for _, r := range dat.Recipes {
		if strings.HasPrefix(r.Name, "Alternate") {
			continue
		}
		row := make([]float64, width)
		norm := 60.0 / float64(r.Time)
		for _, ing := range r.Ingredients {
			key := ing[0].(string)
			count := ing[1].(float64)

			row[index[key]] = -norm * count
		}

		prodKey := r.Product[0].(string)
		prodCount := r.Product[1].(float64)

		row[index[prodKey]] = norm * prodCount

		recipeConversions = append(recipeConversions, row)
		recipesUsed = append(recipesUsed, r)
	}

	for _, x := range recipeConversions {
		for _, v := range x {
			fmt.Printf("%7.2f ", v)
		}
		fmt.Println()
	}

}
