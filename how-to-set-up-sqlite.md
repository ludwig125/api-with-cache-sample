## sqlite のインストール

以下、Ubuntu の場合

```
sudo apt -y install sqlite3

$sqlite3 -version
3.31.1 2020-01-27 19:55:54 3bfa9cc97da10598521b342961df8f5f68c7388fa117345eeb516eaa837balt1
```

## sqlite table 作成まで

#### db 作成

注：cleanArchitectureWebAPI ディレクトリ直下で行う

```
sqlite3 item_db
```

#### db 確認と table 作成

```
$sqlite3 item_db
SQLite version 3.31.1 2020-01-27 19:55:54
Enter ".help" for usage hints.
sqlite> .databases
main: /home/ludwig125/go/src/github.com/ludwig125/architecture/cleanArchitectureWebAPI/item_db
sqlite>
```

以下で table 作成

```
CREATE TABLE item(id INTEGER PRIMARY KEY ASC, name TEXT, age INTEGER);
```

確認

```
sqlite> .tables
item

sqlite> .schema item
CREATE TABLE item(id INTEGER PRIMARY KEY ASC, name TEXT, age INTEGER);
```

#### insert test data

data の例

```
INSERT INTO item(name, age) VALUES("Portman", 32);
INSERT INTO item(name, age) values("Knightley", 35);
INSERT INTO item(name, age) values("Hopkins", 56);
```

確認

```
sqlite> select * from item;
1|Portman|32
2|Knightley|35
3|Hopkins|56
sqlite>
```

#### おまけ

db と table の作成とデータの Insert は以下のように一気にすることもできる

```
sqlite3 item_db 'CREATE TABLE item(id INTEGER PRIMARY KEY ASC, name TEXT, age INTEGER);'
sqlite3 item_db 'INSERT INTO item(name, age) VALUES("Portman", 32);'
sqlite3 item_db 'INSERT INTO item(name, age) values("Knightley", 35);'
sqlite3 item_db 'INSERT INTO item(name, age) values("Hopkins", 56);'
sqlite3 item_db 'INSERT INTO item(name, age) values("Depp", 54);'
sqlite3 item_db 'INSERT INTO item(name, age) values("Watson", 24);'
```
