package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/hashicorp/vault-client-go"
	"github.com/hashicorp/vault-client-go/schema"
	"github.com/joho/godotenv"
)

func main() {
	ctx := context.Background()

	err := godotenv.Load(".env")

	if err != nil {
		log.Fatal(err)
	}

	client, err := vault.New(
		vault.WithAddress(os.Getenv("VAULTLOCATION")),
		vault.WithRequestTimeout(30*time.Second),
	)

	if err != nil {
		log.Fatal(err)
	}

	resp, err := client.Auth.AppRoleLogin(
		ctx,
		schema.AppRoleLoginRequest{
			RoleId:   os.Getenv("ROLE_ID"),
			SecretId: os.Getenv("SECRET_ID"),
		},
	)

	if err != nil {
		log.Fatal(err)
	}

	if err := client.SetToken(resp.Auth.ClientToken); err != nil {
		log.Fatal(err)
	}

	for {

		choice := ""
		fmt.Print("Enter stop to exit or press Enter to continue: ")
		fmt.Scanln(&choice)

		s, err := client.Secrets.KvV2Read(ctx, "textsecret", vault.WithMountPath("golangtest"))

		if err != nil {
			log.Fatal(err)
		}

		log.Println("secret: ", s.Data.Data["name"])
	}
}
