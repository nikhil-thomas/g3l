package arrays_and_slices

func Sum(numbers []int) int {
    sum := 0
    for _, v := range numbers {
        sum += v
    }
    return sum
}

func SumAll(numbersToSum ...[]int) []int {
    sums := []int{}
    //sums := make([]int, len(numbersToSum))
    for _, numbers := range numbersToSum {
        //for i, numbers := range numbersToSum {
        sum := 0
        for _, number := range numbers {
            sum += number
        }
        sums = append(sums, sum)
        //sums[i] = sum
    }
    return sums
}

func SumAllTails(numbers ...[]int) []int {
    tailSums := []int{}
    for _, nums := range numbers {
        if len(nums) <= 1 {
            tailSums = append(tailSums, 0)
            continue
        }
        tailSums = append(tailSums, Sum(nums[1:]))
    }
    return tailSums
}
