package utils

import (
	"fmt"
	"log"

	"golang.org/x/net/context"

	firebase "firebase.google.com/go"
	"firebase.google.com/go/auth"

	"google.golang.org/api/option"
)

func CreateFb() {
	opt := option.WithCredentialsFile("D:/Go_Project/src/go-contacts/utils/go-contacts-ea3f8-firebase-adminsdk-ezmfl-408175d86d.json")
	ctx := context.Background()
	app, err := firebase.NewApp(ctx, nil, opt)
	if err != nil {
		fmt.Printf("error initializing app: %v", err)
		// return app, fmt.Errorf("error initializing app: %v", err)
	}
	client, err := app.Auth(ctx)
	if err != nil {
		fmt.Printf("error initializing app: %v", err)
		// return app, fmt.Errorf("error initializing app: %v", err)
	}
	params := (&auth.UserToCreate{}).
		PhoneNumber("+919898447464")
	u, err := client.CreateUser(ctx, params)
	if err != nil {
		log.Fatalf("error creating user: %v\n", err)
	}
	log.Printf("Successfully created user: %v\n", u)
}
