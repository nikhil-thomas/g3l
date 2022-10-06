package property_based_tests

import (
    "strings"
)

func ConvertToRoman(val int) string {
    var result strings.Builder
    for val > 0 {
        switch {
        case val > 4:
            result.WriteString("V")
            val -= 5
        case val > 3:
            result.WriteString("IV")
            val -= 4
        default:
            result.WriteString("I")
            val--
        }
    }
    return result.String()
}
