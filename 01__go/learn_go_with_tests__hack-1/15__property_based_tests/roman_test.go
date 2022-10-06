package property_based_tests

import "testing"

func TestRomanNumerals(t *testing.T) {
    got := ConvertToRoman(1)
    want := "I"
    if got != want {
        t.Errorf("got %q, want %q", got, want)
    }
}

func TestRomanNumeralsTable(t *testing.T) {
    tests := []struct {
        name   string
        input  int
        output string
    }{
        {"1 get converted to I",
            1,
            "I"},
        {"2 get converted to II",
            2,
            "II"},
        {"3 get converted to III",
            3,
            "III"},
        {"4 gets converted to IV (can't repeat more than 3 times)", 4, "IV"},
        {"5 gets converted to V", 5, "V"},
        {"9 gets converted to IX", 9, "IX"},
        {"10 gets converted to X", 10, "X"},
        {"14 gets converted to XIV", 14, "XIV"},
        {"18 gets converted to XVIII", 18, "XVIII"},
        {"20 gets converted to XX", 20, "XX"},
        {"39 gets converted to XXXIX", 39, "XXXIX"},
    }
    for _, test := range tests {
        t.Run(test.name, func(t *testing.T) {
            got := ConvertToRoman(test.input)
            if got != test.output {
                t.Errorf("want %q, got %q", test.output, got)
            }
        })
    }
}
