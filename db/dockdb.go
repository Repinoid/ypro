package main

import (
	"context"
	"fmt"

	"github.com/jackc/pgx/v5"
)

type DBstruct struct {
	DB *pgx.Conn
}

const dbEndPoint = "postgres://naeel:n@localhost:5434/db"

func main1() {

	ctx := context.Background()

	_, err := InitDBStorage(ctx, dbEndPoint)
	fmt.Println(err)

}

func InitDBStorage(ctx context.Context, dbEndPoint string) (*DBstruct, error) {
	dbStorage := &DBstruct{}
	baza, err := pgx.Connect(ctx, dbEndPoint)
	if err != nil {
		return nil, fmt.Errorf("can't connect to DB %s err %w", dbEndPoint, err)
	}
	// err = TableCreation(ctx, baza)
	// if err != nil {
	// 	return nil, fmt.Errorf("can't create tables in DB %s err %w", dbEndPoint, err)
	// }
	dbStorage.DB = baza
	return dbStorage, nil
}
