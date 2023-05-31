package cmd

import (
	"bytes"
	"encoding/base64"
	"encoding/hex"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"pow-fiat-shamir/sdk"
	"time"

	hclog "github.com/hashicorp/go-hclog"
	"github.com/hashicorp/vault/sdk/helper/jsonutil"
	"github.com/spf13/cobra"
	"go.dedis.ch/kyber/v3/group/edwards25519"
)

type ClientFlags struct {
	LogsFormat bool
}

var (
	clientFlags = &ClientFlags{}
)

func init() {
	client := ClientCmd()

	client.PersistentFlags().BoolVar(&clientFlags.LogsFormat, "log_json", false, "logs format")

	RootCmd.AddCommand(client)
}

func ClientCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:          "client",
		Short:        "Fiat-Shamir client",
		Args:         cobra.ExactArgs(0),
		SilenceUsage: true,
		RunE: func(cmd *cobra.Command, args []string) (err error) {
			// create logger
			logger := hclog.NewInterceptLogger(&hclog.LoggerOptions{
				Output:     os.Stdout,
				Level:      hclog.DefaultLevel,
				JSONFormat: clientFlags.LogsFormat,
			})
			logger.Info("starting fiat-shamir client")

			suite := edwards25519.NewBlakeSHA256Ed25519()

			client := HttpClient()

			round1, err := StartRequest(client, config.Address)
			if err != nil {
				return err
			}

			x := sdk.ComputeX(round1.Message, suite)
			xG, xH := sdk.ComputexGxH(suite, round1.G, round1.H, x)
			c := sdk.ComputeC(suite, round1.G, round1.H, xG, xH)
			v, vG, vH := sdk.ComputeVvGvH(suite, round1.G, round1.H)
			_, rG, rH := sdk.ComputeRrGrH(suite, c, x, v, round1.G, round1.H)

			round2 := &sdk.Round2{
				C: c,

				XH: xH,
				XG: xG,

				RG: rG,
				RH: rH,

				VG: vG,
				VH: vH,
			}

			quotes, err := ResultRequest(client, config.Address, round2)
			if err != nil {
				return err
			}

			fmt.Println("Quotes from “word of wisdom” book:", quotes)

			return nil
		},
	}
	return cmd
}

func StartRequest(client *http.Client, address string) (*sdk.Round1, error) {
	// request
	body, err := getRequestFiatShamir(client, address)
	if err != nil {
		return nil, err
	}

	// get json from body
	type dataJson struct {
		Round1 string
	}
	data := &dataJson{}
	if err := jsonutil.DecodeJSON(body, &data); err != nil {
		return nil, err
	}

	// get base64
	dataBase64 := data.Round1

	// get cbor format
	dataCbor, err := base64.StdEncoding.DecodeString(dataBase64)
	if err != nil {
		return nil, err
	}
	fmt.Println(hex.EncodeToString(dataCbor))

	// get round1
	round1 := &sdk.Round1{}
	if err := round1.UnmarshalBinary(dataCbor); err != nil {
		return nil, err
	}

	return round1, nil
}

func getRequestFiatShamir(client *http.Client, address string) ([]byte, error) {
	req, err := http.NewRequest(http.MethodGet, address+"/fiat-shamir/start", nil)
	if err != nil {
		return nil, err
	}

	res, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return nil, err
	}
	return body, nil
}

func ResultRequest(client *http.Client, address string, round2 *sdk.Round2) (string, error) {
	// decode to cbor format
	dataCbor, err := round2.MarshalBinary()
	if err != nil {
		return "", err
	}

	// decode to base64 format
	dataBase64 := base64.StdEncoding.EncodeToString(dataCbor)

	// decode to json format
	dataJson, err := jsonutil.EncodeJSON(map[string]string{
		"round2": dataBase64,
	})
	if err != nil {
		return "", err
	}

	// send result server
	body, err := resultRequestFiatShamir(client, address, dataJson)
	if err != nil {
		return "", nil
	}

	return string(body), nil
}

func resultRequestFiatShamir(client *http.Client, address string, data []byte) ([]byte, error) {
	payload := bytes.NewReader(data)

	req, err := http.NewRequest(http.MethodPost, address+"/fiat-shamir/result", payload)
	if err != nil {
		return nil, err
	}
	req.Header.Add("Content-Type", "application/json")

	res, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return nil, err
	}

	return body, nil
}

func HttpClient() *http.Client {
	client := &http.Client{Timeout: 10 * time.Second}
	return client
}
