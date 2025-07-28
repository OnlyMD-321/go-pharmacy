package firebase

import (
	"context"
	"log"

	firebase "firebase.google.com/go/v4"
	"google.golang.org/api/option"

	"github.com/OnlyMD-321/go-pharmacy/internal/config"
)

var App *firebase.App

func InitFirebase() {
	opt := option.WithCredentialsFile(config.AppConfig.FirebaseCredentialsPath)

	app, err := firebase.NewApp(context.Background(), nil, opt)
	if err != nil {
		log.Fatalf("Failed to initialize Firebase app: %v", err)
	}

	App = app
	log.Println("Firebase initialized")
}
