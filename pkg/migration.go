package pkg

import (
	"fmt"
	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"time"
)

func Migrat(host string) error {
	time.Sleep(time.Second)
	m, err := migrate.New(
		"file://migrate", fmt.Sprintf(
			"postgres://root:root@%s:5432/root?sslmode=disable", host))
	if err != nil {
		return err
	}
	return m.Up()
}
