package cmd

import (
	"crypto/rand"
	"encoding/base64"
	"net/http"
	"os"
	"pow-fiat-shamir/sdk"

	hclog "github.com/hashicorp/go-hclog"
	"github.com/spf13/cobra"
	"go.dedis.ch/kyber/v3/group/edwards25519"
)

type ServerFlags struct {
	LogsFormat bool
}

var (
	serverFlags = &ServerFlags{}
)

func init() {
	server := ClientCmd()

	server.PersistentFlags().BoolVar(&serverFlags.LogsFormat, "log_json", false, "logs format")

	RootCmd.AddCommand(server)
}

func ServerCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:          "server",
		Short:        "Fiat-Shamir server",
		Args:         cobra.ExactArgs(0),
		SilenceUsage: true,
		RunE: func(cmd *cobra.Command, args []string) (err error) {
			// create logger
			logger := hclog.NewInterceptLogger(&hclog.LoggerOptions{
				Output:     os.Stdout,
				Level:      hclog.DefaultLevel,
				JSONFormat: serverFlags.LogsFormat,
			})
			logger.Info("starting fiat-shamir server")

			suite := edwards25519.NewBlakeSHA256Ed25519()

			mux := http.NewServeMux()
			mux.Handle("/fiat-shamir/start", fiatShamirStartHandler(suite))
			mux.Handle("/fiat-shamir/result", fiatShamirResultHandler(suite))

			logger.Error("Server error:", http.ListenAndServe(config.Address, mux))

			return nil
		},
	}
	return cmd
}

func fiatShamirStartHandler(suite *edwards25519.SuiteEd25519) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// compute crypto data
		m, _ := sdk.GenerateKey(rand.Reader, 32)
		G, H := sdk.ComputeGH(suite)

		// marshal sbor
		round1 := &sdk.Round1{
			Message: m,
			G:       G,
			H:       H,
		}
		dataCbor, err := round1.MarshalBinary()
		if err != nil {
			sdk.RespondError(w, http.StatusInternalServerError, err)
			return
		}

		dataBase64 := base64.StdEncoding.EncodeToString(dataCbor)

		sdk.RespondOk(w, map[string]string{
			"round1": dataBase64,
		})
	})
}

func fiatShamirResultHandler(suite *edwards25519.SuiteEd25519) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		sdk.RespondOk(w, nil)
	})
}
