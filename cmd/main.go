package main

import (
	"bufio"
	"fmt"
	"github.com/buguzei/go-task/internal/config"
	"github.com/buguzei/go-task/internal/models"
	postgres "github.com/buguzei/go-task/internal/repo"
	"github.com/pressly/goose"
	"log"
	"os"
	"strconv"
	"strings"
)

func main() {
	// init configs
	cfg, err := config.InitConfig()
	if err != nil {
		log.Fatal(err)
	}

	// init postgres
	pg := postgres.NewPostgres(cfg.DB)
	defer func() {
		err = pg.DB.Close()
		if err != nil {
			log.Fatal(err)
		}
	}()

	// make migrations
	err = goose.Up(pg.DB, "./migrations/")
	if err != nil {
		log.Fatal(err)
	}

	ordersID, err := readInput()
	if err != nil {
		log.Fatal(err)
	}

	outputMap := make(map[string][]models.Product)

	for _, orderID := range ordersID {
		products, err := pg.GetOrderProducts(orderID)
		if err != nil {
			log.Println(err)
		}

		for _, product := range products {
			outputMap[product.MainRack] = append(outputMap[product.MainRack], product)
		}
	}

	writeOutput(ordersID, outputMap)
}

func writeOutput(input []int, outputMap map[string][]models.Product) {
	output := "=+=+=+=\n"

	var addStr string
	for _, num := range input {
		strNum := strconv.Itoa(num)

		addStr += strNum
	}

	output += addStr + "\n"
	addStr = ""

	for rack, products := range outputMap {
		addStr += fmt.Sprintf("===Стеллаж %s\n", rack)
		for _, product := range products {
			if product.SecondaryRacks != nil {
				addStr += fmt.Sprintf("%s(id=%d)\nзаказ %d, %d шт\nдоп стеллаж: ", product.Name, product.ID, product.OrderID, product.Amount)
				for i, secondaryRack := range product.SecondaryRacks {
					if i == len(product.SecondaryRacks)-1 {
						addStr += secondaryRack + "\n\n"
						continue
					}

					addStr += secondaryRack + ","
				}

				continue
			}

			addStr += fmt.Sprintf("%s(id=%d)\nзаказ %d, %d шт\n\n", product.Name, product.ID, product.OrderID, product.Amount)
		}
	}

	output += addStr

	fmt.Println(addStr)
}

func readInput() ([]int, error) {
	r := bufio.NewReader(os.Stdin)
	input, _ := r.ReadString('\n')

	strs := strings.Split(input, ",")
	nums := make([]int, 0)

	for _, str := range strs {
		num, err := strconv.Atoi(strings.TrimSpace(str))
		if err != nil {
			return nil, fmt.Errorf("atoi error: %w", err)
		}
		nums = append(nums, num)
	}

	return nums, nil
}
