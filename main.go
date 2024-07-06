package main

import (
	"database/sql"
	"flag"
	"fmt"
	go_ora "github.com/sijms/go-ora/v2"
	"io"
	"os"
	"runtime/pprof"
	"runtime/trace"
	"time"
)

const (
	driver   = "oracle"
	server   = "localhost"
	port     = 1521
	service  = "ORABFILE"
	user     = "SYSTEM"
	password = "12345"
)

var options map[string]string

func main() {
	cpuProfileFile, err := os.Create("cpu.prof")
	if err != nil {
		panic(err)
	}
	defer cpuProfileFile.Close()
	if err := pprof.StartCPUProfile(cpuProfileFile); err != nil {
		panic(err)
	}
	defer pprof.StopCPUProfile()

	memProfileFile, err := os.Create("mem.prof")
	if err != nil {
		panic(err)
	}
	defer memProfileFile.Close()
	if err := pprof.WriteHeapProfile(memProfileFile); err != nil {
		panic(err)
	}

	traceFile, err := os.Create("trace.out")
	if err != nil {
		panic(err)
	}
	defer traceFile.Close()
	if err := trace.Start(traceFile); err != nil {
		panic(err)
	}
	defer trace.Stop()

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

	option := flag.String("option", "", "Launch option")
	flag.Parse()

	start := time.Now()

	switch *option {
	case "slow": // 6-8 min, 24 Gb RAM
		fmt.Println("run slow")

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
			defer data.Close()
			length, err := data.GetLength()
			if err != nil {
				panic(err)
			}
			fmt.Println("id:", id, "🏏name:", data.GetFileName(),
				"length:", length, "bytes,", length/1024, "kb,",
				length/(1024*1024), "mb,", length/(1024*1024*1024), "gb.")

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
			fmt.Printf("%d bytes copied\n", n)
		}

	case "faster": // 6-8 sec, 5 Mb RAM
		fmt.Println("run faster")

		rows, err := conn.Query("SELECT FILE_ID, FILE_DATA FROM ORA_BFILE")
		if err != nil {
			panic(err)
		}
		defer rows.Close()

		for rows.Next() {
			item := MyStruct{}
			err = rows.Scan(&item.FileID, &item)
			if err != nil {
				panic(err)
			}

			err = item.Open()
			if err != nil {
				panic(err)
			}
			defer item.Close()
			length, err := item.GetLength()
			if err != nil {
				panic(err)
			}
			fmt.Println("🦆 id:", item.FileID, "name:", item.GetFileName(),
				"length:", length, "bytes,", length/1024, "kb,",
				length/(1024*1024), "mb,", length/(1024*1024*1024), "gb.")

			fo, err := os.Create(fmt.Sprintf("file-%v.txt", item.FileID))
			if err != nil {
				panic(err)
			}
			defer fo.Close()

			err = item.CopyDataToFile(fo)
			if err != nil {
				panic(err)
			}
		}
	}

	fmt.Println("exec time:", time.Since(start))
}

var _ io.ReadCloser = (*MyStruct)(nil)

type MyStruct struct {
	FileID int
	go_ora.BFile
}

func (m *MyStruct) Read(p []byte) (n int, err error) {
	return m.Read(p)
}

func (m *MyStruct) CopyDataToFile(dst *os.File) error {
	src, err := os.Open(fmt.Sprintf("tmp/%s", m.GetFileName()))
	if err != nil {
		return err
	}
	defer src.Close()

	n, err := io.Copy(dst, src)
	if err != nil {
		return err
	}
	fmt.Printf("copy data to file: %d bytes copied\n", n)

	return nil
}
