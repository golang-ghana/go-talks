package main

import (
	"context"
	"fmt"
)

type ctxKey int

const (
	_ ctxKey = iota
	userKey
	locationKey
)

func main() {
	ctx := context.Background()

	ctx = context.WithValue(ctx, userKey, "Jerry")
	ctx = context.WithValue(ctx, locationKey, "Riff Hills")

	user := ctx.Value(userKey).(string)
	location := ctx.Value(locationKey).(string)

	fmt.Println("User: ", user, "Location: ", location)
}
