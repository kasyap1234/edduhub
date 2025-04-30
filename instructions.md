# Navigate to the directory containing the docker-compose file if needed
cd /home/tgt/Desktop/edduhub/build

# Stop and remove the container
docker-compose -f docker-compose.test.yml down

# Start it again in detached mode
docker-compose -f docker-compose.test.yml up -d

# Give it a few seconds to become healthy, then try the migrate command again
sleep 5
migrate -path /home/tgt/Desktop/edduhub/server/db/migrations -database postgres://user:password@localhost:5432/testdb?sslmode=disable up
