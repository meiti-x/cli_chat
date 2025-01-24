Thank you for the great challenge! It was a fun experience, and I learned a lot of new things. However, as I also have a full-time job, I could only dedicate limited time to this project. Despite that, I did my best to implement improvements that would make the project more robust and secure.

There are several additional features and improvements that I couldn’t complete within the given time frame, but I will continue working on them and periodically update the project just for learning purposes. Here’s a list of enhancements I’ve planned:

* Fix the issue where the total number of users shows as 0 when a user leaves the chat after logging in with a single session.
* Create an Entity-Relationship Diagram (ERD) for the database schema.
* Implement a command (#all_messages) to retrieve all user messages.
* Implement a command (#my_messages) to retrieve all messages from a user's chatrooms.
* Use Viper and Cobra for better integration with the CLI.
* Add Loki and Grafana for aggregating and visualizing logs.
* Introduce rate limiting for publishing messages in NATS.
* Enhance the logger with custom fields and trace IDs for better traceability.
* Remove users from Redis on interruption or scan all users during login to ensure consistency.
* Implement authentication for NATS to improve security.
* Mount the log folder for easier log management.
* Implement a health check endpoint for the API server.
* Use Fiber for the API server for better middleware support(add CORS, rate limiting, etc.).
* add pre-commit hooks for linting and formatting with golangci-lint.
* add log rotation for better log management.
* and all TODOs and FIXME in the code base.

## features
1. join and leave chatroom
2. show total users in the chatroom
3. show welcome message when user join the chatroom
4. dockerized the project
5. logging with zap
6. simple auth
7. store user message in db
8. simple test for login and register
9. publish message to nats

## How to run the project
1. go mod tidy
2. run the server docker compose up --build
3. go run cmd/app/client.go
4. or if you want join  the different chatroom you can run the client with the chatroom name go run cmd/app/client.go --chatroom=b

## How run the tests
the basic tests only wrote for login and register.

cd test
docker compose down -v
docker-compose -f test-docker-compose.yml up --build

go test -v