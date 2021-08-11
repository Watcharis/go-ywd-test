package helper

func TwoOldestAges(ages []int) [2]int {
	var res [2]int
	n := 0
	for n < len(ages)-1 {
		// fmt.Println(n)
		if ages[n] > ages[n+1] {
			temp := ages[n]
			ages[n] = ages[n+1]
			ages[n+1] = temp
			n = -1
			// fmt.Println(n)
		}
		n += 1
	}
	res = [2]int{ages[len(ages)-2], ages[len(ages)-1]}
	return res
}

func SortTwoMoreThanValues(ages []int) [2]int {
	var res [2]int
	for i := 0; i < len(ages)-1; i++ {
		if ages[i] > ages[i+1] {
			temp := ages[i]
			ages[i] = ages[i+1]
			ages[i+1] = temp
			i = -1
		}
	}
	// fmt.Println("ages :", ages)
	res = [2]int{ages[len(ages)-2], ages[len(ages)-1]}
	return res
}

func FindDuplicateNumber(data []int) []int {
	var duplicate []int
	for i := 0; i < len(data); i++ {
		for n := 0; n < i; n++ {
			if data[i] == data[n] {
				duplicate = append(duplicate, data[i])
			}
		}
	}
	return duplicate
}
