package repository

import "github.com/jackc/pgx/v5/pgxpool"

type Storage struct {
	conn pgxpool.Conn
}

