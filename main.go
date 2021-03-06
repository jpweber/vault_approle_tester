/*
* @Author: Jim Weber
* @Date:   2016-09-09 10:01:50
* @Last Modified by:   Jim Weber
* @Last Modified time: 2016-09-12 17:20:47
 */

package main

import (
	"encoding/json"
	"flag"
	"log"
	"net/http"
	"time"

	vaultapi "github.com/hashicorp/vault/api"
)

type VaultConfig struct {
	Host         string
	RoleID       string
	SecretID     string
	CubHoleToken string
	ActiveToken  string
}

type CubbyHoleResponse struct {
	RequestID     string            `json:"request_id"`
	LeaseID       string            `json:"lease_id"`
	Renewable     bool              `json:"renewable"`
	LeaseDuration int               `json:"lease_duration"`
	Data          map[string]string `json:"data"`
}

func main() {

	// read in cli arguments for
	// vault host
	// cubby hole token
	// role-id or role name
	vaultHost := flag.String("host", "", "Hostname of Vault Server")
	roleID := flag.String("role", "", "Role ID For Your Application")
	cubHoleToken := flag.String("cubby", "", "Cubby Hole Token")
	secretPath := flag.String("secret", "", "Path to the scret you wish to read from as test verification")
	demo := flag.Bool("demo", false, "Inserts a sleep between steps to make demos clearer")

	// Once all flags are declared, call `flag.Parse()`
	// to execute the command-line parsing.
	flag.Parse()

	// save the cli arguments in the config struct
	vaultConfig := VaultConfig{
		Host:         *vaultHost,
		RoleID:       *roleID,
		CubHoleToken: *cubHoleToken,
	}

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
	log.Println("Reading from cubbyhole with token:", vaultConfig.CubHoleToken)
	client.SetToken(vaultConfig.CubHoleToken)
	secret, err := client.Logical().Read("cubbyhole/response")
	if err != nil {
		log.Println(err)
	}

	// Unmarshal and assign the cubby hole reponse in to the struct
	cubbyResponse := CubbyHoleResponse{}
	if err := json.Unmarshal([]byte(secret.Data["response"].(string)), &cubbyResponse); err != nil {
		panic(err)
	}

	// pull out the secret id from the cubby hole reponse
	vaultConfig.SecretID = cubbyResponse.Data["secret_id"]
	if *demo == true {
		time.Sleep(time.Second * 5)
	}
	log.Println("Received secret ID", vaultConfig.SecretID, "from cubbyhole")

	// login to vault with role-id and secret-id
	// the response will contain a token.
	// this token will be used or all further secret requests
	var IDs = make(map[string]interface{})
	IDs["role_id"] = vaultConfig.RoleID
	IDs["secret_id"] = vaultConfig.SecretID

	log.Println("Authenticating with role ID", IDs["role_id"], "and secret ID", IDs["secret_id"])
	secret, err = client.Logical().Write("auth/approle/login", IDs)
	if err != nil {
		log.Println(err)
	}
	vaultConfig.ActiveToken = secret.Auth.ClientToken
	if *demo == true {
		time.Sleep(time.Second * 5)
	}
	log.Println("Received token", vaultConfig.ActiveToken, "for making future credential requests")

	// make request for the dummy hello world credentials
	log.Println("Requesting data from " + *secretPath) // TODO: make this dynamic based on role id maybe?
	client.SetToken(vaultConfig.ActiveToken)
	secret, err = client.Logical().Read(*secretPath)
	if err != nil {
		log.Println(err)
	}

	if *demo == true {
		time.Sleep(time.Second * 5)
	}
	// output secret data
	log.Println("Secret Data:")
	for k, v := range secret.Data {
		log.Println("Key:", k, "Value:", v)
	}
}
