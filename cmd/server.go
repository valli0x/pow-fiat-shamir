package cmd

import (
	"crypto/rand"
	"encoding/base64"
	"errors"
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
	server := ServerCmd()

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
		switch r.Method {
		case "GET":
		default:
			sdk.RespondError(w, http.StatusMethodNotAllowed, nil)
			return
		}

		// compute crypto data
		m, err := sdk.GenerateKey(rand.Reader, 32)
		if err != nil {
			sdk.RespondError(w, http.StatusInternalServerError, err)
			return
		}
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

		// base64
		dataBase64 := base64.StdEncoding.EncodeToString(dataCbor)

		// send result
		sdk.RespondOk(w, map[string]string{
			"round1": dataBase64,
		})
	})
}

func fiatShamirResultHandler(suite *edwards25519.SuiteEd25519) http.Handler {
	type Body struct {
		Round2 string `json:"round2"`
	}
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case "POST":
		default:
			sdk.RespondError(w, http.StatusMethodNotAllowed, nil)
			return
		}

		// parse body
		body := &Body{}
		_, err := sdk.ParseJSONRequest(r, w, body)
		if err != nil {
			sdk.RespondError(w, http.StatusBadRequest, err)
			return
		}

		// get cbor format
		dataCbor, err := base64.StdEncoding.DecodeString(body.Round2)
		if err != nil {
			sdk.RespondError(w, http.StatusInternalServerError, err)
			return
		}

		// get round2
		round2 := &sdk.Round2{
			C: suite.Scalar(),

			XH: suite.Point(),
			XG: suite.Point(),

			RG: suite.Point(),
			RH: suite.Point(),

			VG: suite.Point(),
			VH: suite.Point(),
		}
		if err := round2.UnmarshalBinary(dataCbor); err != nil {
			sdk.RespondError(w, http.StatusInternalServerError, err)
			return
		}

		// compute
		a, b := sdk.ComputeAB(suite, round2.C, round2.XH, round2.XG, round2.RG, round2.RH)
		valid := sdk.Valid(round2.VG, round2.VH, a, b)

		// validation
		if valid {
			sdk.RespondOk(w, map[string]string{
				"quotes": "Wisdom is not a product of schooling but of the lifelong attempt to acquire it. Albert Einstein",
			})
		} else {
			sdk.RespondError(w, http.StatusMethodNotAllowed, errors.New("Incorrect proof!"))
		}
	})
}
