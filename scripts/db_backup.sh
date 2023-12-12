#!/bin/bash

# Set the current date and time to include in the backup filename
timestamp=$(date +%Y%m%d%H%M%S)

# Set your PostgreSQL container name, username, and database name
container_name="postgres_db"
username="postgres"
database="postgres"

# Set the path where you want to store the backups on your host machine
backup_path="${HOME}/dev"

rm -rf $backup_path/backup_*.sql.tar.gz && echo "Removed old backups"

for f in $backup_path/backup_*.sql;
do
    if [ -f "$f" ]; then
        echo "Moving old backup $f"
        mv $f $f.bak
    fi
done

# Use pg_dump to create a backup
docker exec $container_name pg_dump -U $username -d $database > $backup_path/backup_$timestamp.sql

# Optionally, you may want to compress the backup file to save space
tar -czvf backup_$timestamp.sql.tar.gz $backup_path/backup_$timestamp.sql

rm -rf $backup_path/backup_$timestamp.sql
rm -rf $backup_path/backup_*.sql
