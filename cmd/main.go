package main

import "iv-sukhanov/finance_tracker/internal/spendings"

func main() {
	servise := spendings.New()
	servise.RunBot()
}
