package iteration

//const repeatCount = 5

func Repeat(s string, count int) string {
    repeated_string := ""
    for i := 0; i < count; i++ {
        repeated_string += s
    }
    return repeated_string
}
