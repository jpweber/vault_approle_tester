/*
* @Author: Jim Weber
* @Date:   2016-09-09 10:01:50
* @Last Modified by:   Jim Weber
* @Last Modified time: 2016-09-09 10:46:02
 */

package main

import (
	"flag"
	"log"
	"net/http"

	vaultapi "github.com/hashicorp/vault/api"
)

type VaultConfig struct {
	Host         string
	RoleID       string
	SecretID     string
	CubHoleToken string
	ActiveToken  string
}

func main() {

	// read in cli arguments for
	// vault host
	// cubby hole token
	// role-id or role name
	vaultHost := flag.String("host", "", "Hostname of Vault Server")
	roleID := flag.String("role", "", "Role ID For Your Application")
	cubHoleToken := flag.String("cubby", "", "Cubby Hole Token")

	// Once all flags are declared, call `flag.Parse()`
	// to execute the command-line parsing.
	flag.Parse()

	// save the cli arguments in the config struct
	vaultConfig := VaultConfig{
		Host:         *vaultHost,
		RoleID:       *roleID,
		CubHoleToken: *cubHoleToken,
	}

	// TODO: @debug
	log.Printf("%+v", vaultConfig)

	// init vault client config
	httpClient := &http.Client{}
	clientConfig := vaultapi.Config{
		Address:    "https://" + vaultConfig.Host + ":8200",
		HttpClient: httpClient,
		MaxRetries: 3,
	}

	// intialize vault client
	client, err := vaultapi.NewClient(&clientConfig)
	if err != nil {
		log.Println(err)
	}

	// read value from cubby hole
	// the token received from the cubby hole is the "secret-id"
	client.SetToken(vaultConfig.CubHoleToken)
	secret, err := client.Logical().Read("cubbyhole/response")
	if err != nil {
		log.Println(err)
	}
	// TODO: @debug
	log.Printf("%+v", secret)

	// login to vault with role-id and secret-id
	// the response will contain a token.
	// this token will be used or all further secret requests

	// make request for the dummy hello world credentials

}
