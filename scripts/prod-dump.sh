#!/bin/bash
# README !!!
# This requires ssh key authentication

# scp root@mashdrop.com:'/opt/backups/rethink/'rethinkdb_dump_2022-01-22T03:00:03.tar.gz dump.tar.gz
# date +%Y-%m-%d.%H:%M:%S

target_db="ratatxt_dev"

db_host="root@mashdrop.com"
out_file="dump.tar.gz"

dump_path="/opt/backups/rethink/"
cur_file=`date +%Y-%m-%d`
file_name="rethinkdb_dump_${cur_file}T03:00:03"

# compose file path
target="${dump_path}${file_name}.tar.gz"

# download via scp
#scp $db_host:$target $out_file

# extract files
tar -xvzf $out_file

# prepare database dump files to import
mkdir dump
mv $file_name/$target_db dump/

# import to db
docker run -it --rm \
  -v "$PWD/dump":/dump \
  --link dotagiftx-rethinkdb \
  kudarap/rethinkdb:2.4 \
  rethinkdb import -d /dump -c dotagiftx-rethinkdb --force

# clean up
rm -r $file_name $out_file dump