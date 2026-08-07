#!/bin/bash
echo "Setting up database..."

files=$(find sql/*.sql)

for file in $files; do
    echo "Executing $file..."
    docker compose exec -T database mysql -o snippetbox -u web -pgintoki  < "$file" 2> /dev/null
done

echo "Database setup complete!"
