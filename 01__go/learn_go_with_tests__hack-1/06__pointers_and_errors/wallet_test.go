package wallet

import (
    "fmt"
    "testing"
)

func TestWallet(t *testing.T) {

    t.Run("test deposit", func(t *testing.T) {

        wallet := Wallet{}
        wallet.Deposit(10)
        want := Bitcoins(10)
        fmt.Printf("address from test %v, %v\n", &wallet, &wallet.coins)
        assertBalance(t, wallet, want)
    })

    t.Run("test withdraw", func(t *testing.T) {
        wallet := Wallet{coins: 10}
        err := wallet.Withdraw(5)
        assertNoError(t, err)
        want := Bitcoins(5)
        assertBalance(t, wallet, want)
    })

    t.Run("withdraw insufficient funds", func(t *testing.T) {
        startingBalance := Bitcoins(10)
        wallet := Wallet{startingBalance}
        err := wallet.Withdraw(Bitcoins(100))
        assertBalance(t, wallet, startingBalance)
        assertError(t, err, ErrInsufficientFunds)
    })
}

func assertBalance(t testing.TB, wallet Wallet, want Bitcoins) {
    t.Helper()
    got := wallet.Balance()
    if got != want {
        t.Errorf("got %d want %d", got, want)
    }
}

func assertError(t testing.TB, got, want error) {
    t.Helper()
    fmt.Printf("err %v\n", got)
    t.Logf("err %v\n", got)
    if got == nil {
        t.Fatal("wanted and error but didn't get one")
    }
    if got != want {
        t.Errorf("got %q, want %q", got, want)
    }
}

func assertNoError(t testing.TB, got error) {
    t.Helper()
    if got != nil {
        t.Fatalf("expected no error got one: %v", got)
    }
}
