package main

import "C"
import (
	"sort"
	"strconv"
	"strings"
)

//export column_sort_ver
func column_sort_ver() int {
	return 9
}

//export column_sort_init
func column_sort_init(init *SPH_UDF_INIT, args *SPH_UDF_ARGS, errmsg *ERR_MSG) int32 {
	if args.arg_count < 1 || args.arg_count > 2 || args.arg_type(0) != SPH_UDF_TYPE_STRING {
		return errmsg.say("Please call valid format: COLUMN_SORT(column string, type string(asc|desc))")
	}

	if args.arg_count == 2 && args.arg_type(1) != SPH_UDF_TYPE_STRING {
		return errmsg.say("The second argument can only take the int value \"asc\" (default) or \"desc\". " +
			"Where desc - sort in descending order, asc - in ascending order")
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

func Sorting(nums []int, target int) []string {
	if target == 0 {
		sort.Ints(nums)
	} else {
		sort.Sort(sort.Reverse(sort.IntSlice(nums)))
	}

	var res []string
	for _, value := range nums {
		res = append(res, strconv.Itoa(value))
	}

	return res
}

//export column_sort
func column_sort(init *SPH_UDF_INIT, args *SPH_UDF_ARGS, errf *ERR_FLAG) uintptr {
	nums, err := ParseInteger(args.stringval(0))
	if err != nil {
		return args.return_string("0")
	}

	var target int
	if args.Arg_count() == 2 && args.stringval(1) == "desc" {
		target = 1
	} else {
		target = 0
	}

	return args.return_string(Sorting(nums, target)[0])
}
