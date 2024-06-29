# Oracle-bfile

## Prepare the database and files
```make init-data```

## Run slow solution
```go run main.go -option slow```

## Run faster solution
```go run main.go -option faster```

## Check pprof
```bash
go tool pprof cpu.prof

go tool pprof mem.prof

go tool trace trace.out
```
