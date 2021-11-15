package main

import "gitlab.com/aosorgin/pgutil/cmd/pgquerycmd"

func main() {
	if err := pgquerycmd.RootCmd.Execute(); err != nil {
		panic(err)
	}
}
