#!/bin/bash -eu

rm -rf item_db

sqlite3 item_db 'CREATE TABLE item(id INTEGER PRIMARY KEY ASC, name TEXT, price INTEGER);'

rm -rf score_db

sqlite3 score_db 'CREATE TABLE score(id INTEGER PRIMARY KEY ASC, name TEXT, score INTEGER);'


# sqlite3 item_db 'INSERT INTO item(name, price) VALUES("SmartPhone", 32);'

# # 以下だと10000件入れるのに１分22秒かかった
# for i in $(seq 1 10000)
# do
#   name="Item$i"
#   price=$i
#   sqlite3 item_db "INSERT INTO item(name, price) VALUES(\"$name\", $price);"
# done


# 以下だと10000件は以下の問題が生じた
# https://stackoverflow.com/questions/51710864/usr-bin-sqlite3-argument-list-too-long
# N=10000
# for i in $(seq 1 1000)
# do
#   name="Item$i"
#   price=$i

#   if [ $i -eq $N ]
#   then
#     insertData=${insertData}"(\"$name\", $price);"
#     break
#   fi

#   insertData=${insertData}"(\"$name\", $price),"
# done
# sqlite3 item_db "INSERT INTO item(name, price) VALUES $insertData"

insertItemFn () {
  insertData=""

  for i in $(seq $1 $2)
  do
    name="Item$i"
    price=$i

    if [ $i -eq $2 ]
    then
      insertData=${insertData}"(\"$name\", $price);"
      break
    fi

    insertData=${insertData}"(\"$name\", $price),"
  done

  sqlite3 item_db "INSERT INTO item(name, price) VALUES $insertData"
}

insertScoreFn () {
  insertData=""

  for i in $(seq $1 $2)
  do
    name="Item$i"
    score=$(($i%10)) # 10の余りをスコアにする

    if [ $i -eq $2 ]
    then
      insertData=${insertData}"(\"$name\", $score);"
      break
    fi

    insertData=${insertData}"(\"$name\", $score),"
  done

  sqlite3 score_db "INSERT INTO score(name, score) VALUES $insertData"
}


for i in $(seq 0 9)
do
  start=$(($(($i*1000))+1))
  end=$(($(($i+1))*1000))

  mod=$((end%10000))
  if [ $mod -eq 0 ]
  then
    echo $start $end
  fi

  insertScoreFn $start $end
done
# => これで10000件のItemとScoreができる DBのサイズは200kずつ