package main

import (
	"fmt"
	"log"

	webpush "github.com/SherClockHolmes/webpush-go"
)

func main() {
	priv, pub, err := webpush.GenerateVAPIDKeys()
	if err != nil {
		log.Fatalf("Failed to generate VAPID keys: %v", err)
	}
	fmt.Printf("VAPID_PUBLIC_KEY=%s\n", pub)
	fmt.Printf("VAPID_PRIVATE_KEY=%s\n", priv)
}
