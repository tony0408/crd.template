package db

import (
	"database/sql"

	_ "github.com/lib/pq"
)

type Postgres struct {
	db *sql.DB
}

func (p *Postgres) InsertId(sql string) (int64, error) {

	// log.Printf("sql: %s", sql)
	var id int64
	err := p.db.QueryRow(sql).Scan(&id)

	return id, err
}

func (p *Postgres) Query(sql string) (*sql.Rows, error) {

	// log.Printf("sql: %s", sql)
	rows, err := p.db.Query(sql)
	return rows, err
}

func (p *Postgres) Exec(sql string) (sql.Result, error) {

	// log.Printf("sql: %s", sql)
	result, err := p.db.Exec(sql)

	// lastInsertId := 0
	// err = db.QueryRow(sql).Scan(&lastInsertId)
	if err != nil {
		return nil, err
	}

	return result, err
}

func (p *Postgres) Close() {
	p.db.Close()
}

func CreatePostgres(uri string) (*Postgres, error) {
	p := &Postgres{}
	var err error
	p.db, err = sql.Open("postgres", uri)
	if err != nil {
		return nil, err
	}

	err = p.db.Ping()
	return p, err
}
