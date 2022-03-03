package wallet

import (
    "errors"
    "fmt"
)

var (
    ErrInsufficientFunds = errors.New("cannot withdraw, insufficient funds")
)

type Bitcoins int64

func (b Bitcoins) String() string {
    return fmt.Sprintf("%d BTC", b)
}

type Wallet struct {
    coins Bitcoins
}

func (w *Wallet) Deposit(coins Bitcoins) {
    w.coins += coins
    fmt.Printf("address from method %v, %v\n", &w, &w.coins)
    fmt.Printf("coins: %d\n", w.coins)
    fmt.Printf("coins: %s\n", w.coins)
}

func (w *Wallet) Balance() Bitcoins {
    return w.coins
}

func (w *Wallet) Withdraw(coins Bitcoins) error {
    if w.coins < coins {
        return ErrInsufficientFunds
    }
    w.coins -= coins
    return nil
}
