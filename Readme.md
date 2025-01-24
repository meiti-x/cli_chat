Thank you for the great challenge! It was a fun experience, and I learned a lot of new things. However, as I also have a full-time job, I could only dedicate limited time to this project. Despite that, I did my best to implement improvements that would make the project more robust and secure.

There are several additional features and improvements that I couldn’t complete within the given time frame, but I will continue working on them and periodically update the project for learning purposes. Here’s a list of enhancements I’ve planned:

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
* and all TODOs in the code.

## How to run the project

## How run the tests
