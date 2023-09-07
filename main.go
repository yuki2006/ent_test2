package main

import (
	"context"
	"fmt"
	"go_test/ent"
	"go_test/ent/user"
	"log"

	"entgo.io/ent/dialect/sql"
	"entgo.io/ent/dialect/sql/schema"
	"github.com/go-sql-driver/mysql"
)

var client *ent.Client

func main() {
	mc := mysql.Config{
		User:                 "root",
		Passwd:               "",
		Net:                  "tcp",
		Addr:                 "localhost" + ":" + "3306",
		DBName:               "ent",
		AllowNativePasswords: true,
		ParseTime:            true,
	}
	drv, _ := sql.Open("mysql", mc.FormatDSN())
	client = ent.NewClient(ent.Driver(drv), ent.Debug())
	ctx := context.Background()
	err := client.Schema.Create(ctx, schema.WithDropIndex(true), schema.WithDropColumn(true), schema.WithForeignKeys(true))
	if err != nil {
		log.Fatal(err)
		return
	}
	client.User.Delete().ExecX(ctx)
	ctx = context.Background()
	tx, err := client.Tx(ctx)
	if err != nil {
		log.Fatalf("failed creating transaction: %v", err)
	}

	users := []*ent.User{
		{Name: "Alice", Age: nil},
		{Name: "Bob", Age: Ptr(25)},
	}

	for _, u := range users {
		_, err = tx.User.Create().SetName(u.Name).SetNillableAge(u.Age).Save(ctx)
		if err != nil {
			log.Fatalf("failed creating user: %v", err)
		}
	}

	if err := tx.Commit(); err != nil {
		log.Fatalf("failed committing transaction: %v", err)
	}

	tx, err = client.Tx(ctx)
	if err != nil {
		panic(err)
		return
	}
	all, err := tx.User.Query().Order(user.ByAge(sql.OrderNullsFirst())).All(ctx)
	if err != nil {
		fmt.Printf("err: %v\n", err)
		return
	}
	fmt.Printf("all: %v\n", all)
}

func Ptr[T any](v T) *T {
	return &v
}
