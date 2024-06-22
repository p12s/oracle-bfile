package main

import (
	"database/sql"
	"fmt"
	go_ora "github.com/sijms/go-ora/v2"
	"io"
	"os"
)

func main() {
	driver := "oracle"
	server := "localhost"
	port := 1521
	service := "ORABFILE"
	user := "SYSTEM"
	password := "12345"
	var options map[string]string

	connStr := go_ora.BuildUrl(server, port, service, user, password, options)
	conn, err := sql.Open(driver, connStr)
	if err != nil {
		panic(err)
	}
	defer conn.Close()

	err = conn.Ping()
	if err != nil {
		panic(err)
	}

	rows, err := conn.Query("SELECT FILE_ID, FILE_DATA FROM ORA_BFILE")
	if err != nil {
		panic(err)
	}

	defer rows.Close()
	var id int
	var data go_ora.BFile

	for rows.Next() {
		err = rows.Scan(&id, &data)
		if err != nil {
			panic(err)
		}

		err = data.Open()
		if err != nil {
			panic(err)
		}
		length, err := data.GetLength()
		if err != nil {
			panic(err)
		}
		fmt.Println("id:", id, "🏏name:", data.GetFileName(),
			"length:", length, "bytes", length/1024, "kb",
			length/(1024*1024), "mb", length/(1024*1024*1024), "gb")

		b, err := data.Read()
		if err != nil {
			panic(err)
		}

		fo, err := os.Create(fmt.Sprintf("file-%v.txt", id))
		if err != nil {
			panic(err)
		}
		defer fo.Close()

		n, err := fo.Write(b)
		if err != nil && err != io.EOF {
			panic(err)
		}
		fmt.Println(n)
	}

}
