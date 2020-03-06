package app

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/99designs/gqlgen/graphql/handler"
	"github.com/99designs/gqlgen/graphql/playground"
	"github.com/go-chi/chi"
	database "github.com/moutend/gqlgen-todoapp/internal/db/mysql"
	"github.com/moutend/gqlgen-todoapp/internal/graph"
	"github.com/moutend/gqlgen-todoapp/internal/graph/generated"
	"github.com/moutend/gqlgen-todoapp/internal/middleware/auth"
	"github.com/moutend/gqlgen-todoapp/internal/middleware/db"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var RootCommand = &cobra.Command{
	Use:  "todoapp",
	RunE: rootRunE,
}

func rootRunE(cmd *cobra.Command, args []string) error {
	if path, _ := cmd.Flags().GetString("config"); path != "" {
		viper.SetConfigFile(path)
	}

	viper.SetEnvKeyReplacer(strings.NewReplacer(".", "_", "-", "_"))
	viper.AutomaticEnv()

	if err := viper.ReadInConfig(); err != nil {
		return err
	}

	opt := database.Option{
		User:     viper.GetString("mysql.user"),
		Password: viper.GetString("mysql.pass"),
		Host:     viper.GetString("mysql.host"),
		Port:     viper.GetInt("mysql.port"),
		Database: viper.GetString("mysql.dbname"),
	}

	if err := database.Initialize(opt); err != nil {
		return err
	}
	if err := database.Migrate(viper.GetString("migration.path")); err != nil {
		return err
	}

	server := handler.NewDefaultServer(generated.NewExecutableSchema(generated.Config{Resolvers: &graph.Resolver{}}))

	router := chi.NewRouter()
	router.Use(db.Middleware())
	router.Use(auth.Middleware())
	router.Handle("/", playground.Handler("GraphQL playground", "/query"))
	router.Handle("/query", server)

	address := fmt.Sprintf(
		"%s:%s",
		viper.GetString("web.host"),
		viper.GetString("web.port"),
	)

	cmd.Printf("Listening on %s\n", address)

	return http.ListenAndServe(address, router)
}

func init() {
	RootCommand.PersistentFlags().StringP("config", "c", "", "Path to configuration file")
}
