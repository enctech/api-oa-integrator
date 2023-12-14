#!/bin/bash

# Set your PostgreSQL container name, username, and database name
container_name="postgres_db"
username="postgres"
database="postgres"

# Set the path where you want to store the backups on your host machine
backup_path="${HOME}/dev"

# Set the maximum number of backups to keep
max_backups=100

# Create a backup timestamp
timestamp=$(date +%Y%m%d%H%M%S)

# Use pg_dump to create a backup
docker exec $container_name pg_dump -U $username -d $database > $backup_path/backup_$timestamp.sql

# Compress the backup file to save space
tar -czvf backup_$timestamp.sql.tar.gz $backup_path/backup_$timestamp.sql

# Count the number of existing backups
backup_count=$(ls -1 $backup_path | grep -c "^backup_.*\.sql$")

# Delete the oldest backup if the number exceeds the maximum
if [ $backup_count -gt $max_backups ]; then
    oldest_backup=$(ls -1t $backup_path | grep "^backup_.*\.sql$" | tail -n 1)
    rm $backup_path/$oldest_backup
fi

