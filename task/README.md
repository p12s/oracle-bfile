# Oracle BFile

Суть задания:
Имеется БД `oracle` в докере, в ней содержатся несколько [BFILE](https://docs.oracle.com/en/database/oracle/oracle-database/21/adlob/BFILEs.html)
Необходимо прочитать содержимое этих файлов и записать их в файл на локальный диск.
Будут также оцениваться дополнительные оптимизации

Цель задания: научится работать с интерфейсами и получить навыки работы по оптимизации приложения

Ограничения:
1. Необходимо использовать только библиотеку [go-ora v2.8.19](https://github.com/sijms/go-ora/releases/tag/v2.8.19)
2. [Oracle container](https://hub.docker.com/r/gvenzl/oracle-xe). (Если Apple Silicon то через [colima](https://github.com/abiosoft/colima))
3. Программа генерирует 8gb данных

Дополнительные задания:
1. Сделать патч для библиотеки go-ora, чтобы методы библиотеки для работы с BFile реализовывали интерфейс [io.ReadCloser](https://pkg.go.dev/io#ReadCloser)
2. Сравнить производительности решений с интерфейсом и без

Первый запуск:
```bash
make init-data
```
