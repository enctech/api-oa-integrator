#!/bin/bash

# Set the current date and time to include in the backup filename
timestamp=$(date +%Y%m%d%H%M%S)

# Set your PostgreSQL container name, username, and database name
container_name="postgres_db"
username="postgres"
database="postgres"

# Set the path where you want to store the backups on your host machine
backup_path="${HOME}/dev"

# Rename all current backups to .bak. This will be used to query old backups deletion.
for f in $backup_path/backup_*.sql.tar.gz;
do
    if [ -f "$f" ]; then
        echo "Moving old backup $f"
        mv $f $f.bak
    fi
done

# Use pg_dump to create a new backup
docker exec $container_name pg_dump -U $username -d $database > $backup_path/backup_$timestamp.sql

# Compress the backup file to save space
tar -czvf backup_$timestamp.sql.tar.gz $backup_path/backup_$timestamp.sql

# Remove old temp backups.
rm -rf $backup_path/*.bak
