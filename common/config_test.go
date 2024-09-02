package common

import (
	"os"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestConfig(t *testing.T) {
	t.Run("Can set env locale", func(t *testing.T) {
		for _, test := range []struct {
			env  string
			want Envs
		}{
			{
				env: "Prod",
				want: Envs{
					isLocal: false,
					isStage: false,
					isProd:  true,
				},
			},
			{
				env: "prod",
				want: Envs{
					isLocal: false,
					isStage: false,
					isProd:  true,
				},
			},
			{
				env: "Production",
				want: Envs{
					isLocal: false,
					isStage: false,
					isProd:  true,
				},
			},
			{
				env: "Stage",
				want: Envs{
					isLocal: false,
					isStage: true,
					isProd:  false,
				},
			},
			{
				env: "stage",
				want: Envs{
					isLocal: false,
					isStage: true,
					isProd:  false,
				},
			},
			{
				env: "Staging",
				want: Envs{
					isLocal: false,
					isStage: true,
					isProd:  false,
				},
			},
			{
				env: "Local",
				want: Envs{
					isLocal: true,
					isStage: false,
					isProd:  false,
				},
			},
			{
				env: "Dev",
				want: Envs{
					isLocal: true,
					isStage: false,
					isProd:  false,
				},
			},
			{
				env: "",
				want: Envs{
					isLocal: true,
					isStage: false,
					isProd:  false,
				},
			},
		} {
			t.Run(test.env, func(t *testing.T) {
				os.Setenv("ENV", test.env)

				cfg, err := GetConfig()
				require.Nil(t, err)

				require.Equal(t, test.want.isLocal, cfg.isLocal)
				require.Equal(t, test.want.isStage, cfg.isStage)
				require.Equal(t, test.want.isProd, cfg.isProd)

				os.Unsetenv("ENV")
			})
		}

	})
}
