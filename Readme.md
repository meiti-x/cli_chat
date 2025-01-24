## **Project Documentation**

### **Introduction**
Thank you for the incredible challenge! This project was a rewarding experience that allowed me to learn and explore new concepts. Due to my full-time job, I had limited time to work on this, but I put my best effort into making the project as robust and secure as possible.

While I made significant progress, there are additional features and improvements I couldn’t complete within the time frame. I plan to continue enhancing the project and updating it regularly as a learning exercise. Below are the planned improvements and the features I implemented during this challenge.

---

### **Planned Improvements**
1. **User Tracking Bug Fix**: Address the issue where the total user count shows as 0 when a user leaves the chat after logging in with a single session.
2. **Database Schema Documentation**: Create an **Entity-Relationship Diagram (ERD)** for the database schema.
3. **Commands for Retrieving Messages**:
    - Add a `#all_messages` command to retrieve all user messages.
    - Add a `#my_messages` command to fetch messages from the user's chatrooms.
4. **CLI Improvements**: Use **Viper** and **Cobra** for better command-line interface integration.
5. **Log Aggregation and Visualization**: Integrate **Loki** and **Grafana** to aggregate and visualize logs effectively.
6. **Rate Limiting**: Add rate limiting for publishing messages in **NATS** to improve performance and prevent abuse.
7. **Enhanced Logging**: Add custom fields and trace IDs to the logger for better traceability.
8. **Redis Consistency**: Remove users from **Redis** on interruption or scan all users during login to ensure consistency.
9. **NATS Security**: Implement authentication for **NATS** to enhance security.
10. **Log Management**: Mount a log folder for easier log management and add log rotation for better handling of large log files.
11. **Health Check Endpoint**: Implement a health check endpoint for the API server to monitor its status.
12. **API Framework**: Migrate to **Fiber** for the API server to benefit from its robust middleware support (e.g., CORS, rate limiting).
13. **Pre-Commit Hooks**: Add pre-commit hooks for linting and formatting using **golangci-lint**.
14. **Swagger Documentation**: Write Swagger documentation for the API server to improve usability.
15. **TODOs and FIXMEs**: Address all TODOs and FIXMEs in the codebase to ensure completeness and maintainability.

---

### **Implemented Features**
1. **Chatroom Functionality**:
    - Join and leave chatrooms.
    - Show the total number of users in a chatroom.
    - Display a welcome message when a user joins a chatroom.
2. **Project Containerization**: Dockerized the project for easier deployment and testing.
3. **Logging**: Added logging using **Zap** for efficient log management.
4. **Authentication**: Implemented simple authentication for users.
5. **Message Storage**: Stored user messages in the database for future retrieval.
6. **Testing**: Wrote basic tests for login and registration functionality.
7. **NATS Integration**: Added functionality to publish messages to **NATS**.

---

### **How to Run the Project**

1. Install dependencies:
   ```bash
   go mod tidy
   ```
2. Start the server:
   ```bash
   docker compose up --build
   ```
3. Run the client:
   ```bash
   go run cmd/app/client.go
   ```
4. To join a different chatroom, specify the chatroom name:
   ```bash
   go run cmd/app/client.go --chatroom=b
   ```

---

### **How to Run the Tests**

The basic tests cover login and registration functionality.

1. Navigate to the test directory:
   ```bash
   cd test
   ```  
2. Clean up any existing containers and volumes:
   ```bash
   docker compose down -v
   ```  
3. Start the test environment:
   ```bash
   docker-compose -f test-docker-compose.yml up --build
   ```  
4. Run the tests:
   ```bash
   go test -v
   ```  
