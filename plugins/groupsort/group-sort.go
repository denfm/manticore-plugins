package main

import "C"
import (
	"regexp"
	"sort"
	"strconv"
	"strings"
)

//export group_sort_ver
func group_sort_ver() int {
	return 9
}

//export group_sort_init
func group_sort_init(init *SPH_UDF_INIT, args *SPH_UDF_ARGS, errmsg *ERR_MSG) int32 {
	if args.arg_count < 2 || args.arg_count > 3 || args.arg_type(0) != SPH_UDF_TYPE_STRING {
		return errmsg.say("Please call valid format: " +
			"GROUP_SORT(group_concat(id):required, group_concat(status):required, group_concat(similar.<id>):optional)")
	}
	return 0
}

func ParseInteger(str string) ([]int, error) {
	var items []int

	for _, v := range strings.Split(str, ",") {
		id, err := strconv.Atoi(v)
		if err != nil {
			return items, err
		}
		items = append(items, id)
	}
	return items, nil
}

func ParseSimilarPosition(str string) ([]int, error) {
	var items []int

	for _, v := range strings.Split(str, ",") {
		parsePosition := regexp.MustCompile(`(?m)^0+`).ReplaceAllString(v, "")
		if parsePosition == "" {
			parsePosition = "0"
		}

		position, err := strconv.Atoi(parsePosition)
		if err != nil {
			return items, err
		}

		items = append(items, position)
	}
	return items, nil
}

func Sorting(productsIdsMap []int, statusIdsMap []int, similarMap []int) []string {
	productsValueMap := map[float64]int{}
	lenSimilar := len(similarMap)

	var sortData []float64
	var newProductsSort []string

	for k, productId := range productsIdsMap {
		var value float64
		if lenSimilar > 0 {
			value = float64(statusIdsMap[k]) + float64(similarMap[k])/1000
		} else {
			value = float64(statusIdsMap[k])
		}

		sortData = append(sortData, value)
		productsValueMap[value] = productId
	}

	sort.Float64s(sortData)
	for _, value := range sortData {
		newProductsSort = append(newProductsSort, strconv.Itoa(productsValueMap[value]))
	}

	return newProductsSort
}

//export group_sort
func group_sort(init *SPH_UDF_INIT, args *SPH_UDF_ARGS, errf *ERR_FLAG) uintptr {
	productsIdsMap, err := ParseInteger(args.stringval(0))
	if err != nil {
		return args.return_string(args.stringval(0))
	}

	statusIdsMap, err := ParseInteger(args.stringval(1))
	if err != nil {
		return args.return_string(args.stringval(1))
	}

	var similarMap []int
	if args.Arg_count() == 3 {
		similarMap, err = ParseSimilarPosition(args.stringval(2))
		if err != nil {
			return args.return_string(args.stringval(1))
		}
	}

	return args.return_string(strings.Join(Sorting(productsIdsMap, statusIdsMap, similarMap), ","))
}
