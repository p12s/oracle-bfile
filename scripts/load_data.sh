#!/bin/bash

dd if=/dev/urandom of=/tmp/my_file1.txt bs=1G count=4
dd if=/dev/urandom of=/tmp/my_file2.txt bs=1G count=4

# Подключение к базе данных
sqlplus 'SYSTEM/12345@(DESCRIPTION=(ADDRESS=(PROTOCOL=TCP)(HOST=localhost)(PORT=1521))(CONNECT_DATA=(SERVICE_NAME=ORABFILE)))' <<EOF

CREATE TABLE ORA_BFILE (
    FILE_ID NUMBER(10) NOT NULL,
    FILE_DATA BFILE
);

CREATE OR REPLACE DIRECTORY TMP_DIR AS '/tmp';

DECLARE
  l_bfile1  BFILE := BFILENAME('TMP_DIR', 'my_file1.txt');
  l_bfile2  BFILE := BFILENAME('TMP_DIR', 'my_file2.txt');
BEGIN
  INSERT INTO ORA_BFILE (FILE_ID, FILE_DATA) VALUES (1, l_bfile1);
  INSERT INTO ORA_BFILE (FILE_ID, FILE_DATA) VALUES (2, l_bfile2);
  COMMIT;
END;
/

EOF
