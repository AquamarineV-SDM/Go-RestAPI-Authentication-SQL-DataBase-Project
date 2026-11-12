# Go REST API Event Registration System

A complete RESTful API built with Go (Golang) and the Gin framework. This project serves as a backend system for managing events and user registrations.

## Features
- **User Authentication:** Signup and Login system with password hashing (bcrypt).
- **JWT Authorization:** Secured endpoints using JSON Web Tokens (JWT).
- **Event Management:** Create, Read, Update, and Delete (CRUD) events.
- **Role-Based Access:** Only event creators can update or delete their own events.
- **Registration System:** Users can register for or cancel their attendance at events (Many-to-Many relationship).

## Tech Stack
- **Language:** Go (Golang)
- **Web Framework:** [Gin](https://github.com/gin-gonic/gin)
- **Database:** SQLite
- **Security:** bcrypt (Password Hashing), golang-jwt (Authentication Token)

## How to Run
1. Clone the repository:
   ```bash
   git clone [\[https://github.com/your-username/your-repo-name.git\](https://github.com/your-username/your-repo-name.git)](https://github.com/AquamarineV-SDM/Go-RestAPI-Authentication-SQL-DataBase-Project.git)

    Install dependencies:
Bash
go mod tidy

    Run the server:
Bash
go run .
The server will start on http://localhost:8080.