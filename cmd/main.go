package main

import (
	"fmt"
	"github.com/Futturi/ozst/inmemory"
	"github.com/Futturi/ozst/pkg"
	"github.com/Futturi/ozst/postgre"
	"github.com/spf13/viper"
	"log/slog"
	"net/http"
	"os"

	"github.com/99designs/gqlgen/graphql/handler"
	"github.com/99designs/gqlgen/graphql/playground"
	"github.com/Futturi/ozst"
)

const defaultPort = "8080"

func main() {
	logg := slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelInfo}))
	slog.SetDefault(logg)
	InitCfg()
	var srv *handler.Server
	c := 0
l:
	for {
		var tp string
		fmt.Print("Выбери тип сохранения данных(inmemory или bd): ")
		fmt.Scan(&tp)
		switch tp {
		case "inmemory":
			srv = handler.NewDefaultServer(ozst.NewExecutableSchema(ozst.Config{Resolvers: &inmemory.Resolver{Store: inmemory.NewStore()}}))
			break l
		case "bd":
			cfg := pkg.Config{
				Host:     viper.GetString("host"),
				Port:     viper.GetString("port"),
				Username: viper.GetString("username"),
				Password: viper.GetString("password"),
				Dbname:   viper.GetString("dbname"),
				Sslmode:  viper.GetString("sslmode"),
			}
			db, err := pkg.InitPostgres(cfg)
			if err != nil {
				slog.Error("error", err)
			}
			slog.Info("bd started on port: ", cfg.Port)
			err = pkg.Migrat(viper.GetString("host"))
			if err != nil {
				slog.Info("error with migrat", "error", err)
			}
			srv = handler.NewDefaultServer(ozst.NewExecutableSchema(ozst.Config{Resolvers: &postgre.Resolver{Db: db}}))
			break l
		default:
			fmt.Println("incorrect тип сохранения данных")
			c++
			if c == 5 {
				cfg := pkg.Config{
					Host:     viper.GetString("host"),
					Port:     viper.GetString("port"),
					Username: viper.GetString("username"),
					Password: viper.GetString("password"),
					Dbname:   viper.GetString("dbname"),
					Sslmode:  viper.GetString("sslmode"),
				}
				db, err := pkg.InitPostgres(cfg)
				if err != nil {
					slog.Error("error", err)
				}
				slog.Info("bd started on port: ", cfg.Port)
				err = pkg.Migrat(viper.GetString("host"))
				if err != nil {
					slog.Info("error with migrat", "error", err)
				}
				srv = handler.NewDefaultServer(ozst.NewExecutableSchema(ozst.Config{Resolvers: &postgre.Resolver{Db: db}}))
				break l
			}
			break
		}
	}
	port := os.Getenv("PORT")
	if port == "" {
		port = defaultPort
	}
	http.Handle("/", playground.Handler("GraphQL playground", "/query"))
	http.Handle("/query", srv)

	slog.Info("connect to http://localhost:%s/ for GraphQL playground", "port", port)
	slog.Error("error", http.ListenAndServe(":"+port, nil))
}

func InitCfg() error {
	viper.SetConfigType("yml")
	viper.SetConfigName("config")
	viper.AddConfigPath("config")
	return viper.ReadInConfig()
}
