# File Storage Service with User Management

A RESTful API for secure file uploads and user authentication, featuring JWT-based authorization and metadata tracking.

## Features

- 📁 **File Management**
  - Upload/download files with metadata
  - List user files with sorting/filtering
  - Automatic MIME type detection
  - Human-readable file sizes (KB, MB, GB)

- 🔐 **Authentication**
  - JWT-based user registration/login
  - Password hashing with bcrypt
  - Protected endpoints

- 📊 **Metadata Tracking**
  - File size in bytes and human-readable format
  - Upload timestamps
  - User ownership tracking
  - File paths and unique IDs

## API Endpoints

### Authentication

| Method | Endpoint    | Description               |
|--------|-------------|---------------------------|
| POST   | /register   | Register new user         |
| POST   | /login      | Login and get JWT token   |

### File Operations (Require Auth)

| Method | Endpoint    | Description               |
|--------|-------------|---------------------------|
| POST   | /api/upload | Upload file               |
| GET    | /api/files  | List all user files       |
| GET    | /api/storage/remaining | Get storage usage |

## Example Responses

### File Metadata Example (Upload Response)
{
    "id": "1cb9b828f2b510f6",
    "userId": "riza",
    "filename": "TotalForeignExchangeSummaryReport.xls",
    "path": "riza/1cb9b828f2b510f6_TotalForeignExchangeSummaryReport.xls",
    "size": 5632,
    "sizeHuman": "5.5 KB",
    "uploadedAt": "2025-03-25T22:58:25.063405253+05:30",
    "mimeType": "application/octet-stream"
}
### JWT token (login response)
{
    "token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJleHAiOjE3NDMwMTAyNDYsInVzZXJuYW1lIjoicml6YSJ9.AippWpu8_0RdhrHrvWjeQbP9afqZV8vFKA5EeGWgK8c"
}
### user register 
{
    "message": "User registered successfully"
}

### get storage details
{
    "AllocatedStorage": "10.0 MB",
    "usedStorage": "722.3 KB",
    "remaining": "9.3 MB"
}

### Environment config details
JWT_SECRET=your_strong_secret_key_here
STORAGE_PATH=./storage
DEFAULT_USER_QUOTA=10485760  # 10MB
SERVER_PORT=8080


### Run the server:
  * go run cmd/server/main.go

### flow

.
├── README.md
├── cmd
│   └── server
│       └── main.go
├── config.env
├── data
│   └── users.db
├── go.mod
├── go.sum
├── internal
│   ├── handlers
│   │   ├── create_user.go
│   │   ├── get_usage.go
│   │   ├── list_files.go
│   │   ├── storage_service.go
│   │   ├── upload_files.go
│   │   └── user_service.go
│   ├── middleware
│   │   ├── auth.go
│   │   └── logging.go
│   ├── models
│   │   └── storage.go
│   ├── services
│   │   ├── auth.go
│   │   └── storage.go
│   └── utils
│       └── util.go
├── pkg
│   └── config
│       └── config.go
└── storage
    


### Explanation of the Structure:

- **README.md**: Documentation file.
- **cmd/server/main.go**: Entry point for the server application.
- **config.env**: Environment configuration file.
- **data/users.db**: Database file for storing user-related data.
- **go.mod & go.sum**: Go module and dependency management files.
- **internal**: Contains core business logic of the application.
  - **handlers**: Handlers for different endpoints like creating users, file uploads, etc.
  - **middleware**: Code for middleware such as authentication and logging.
  - **models**: Contains the application's data models.
  - **services**: Contains business logic for authentication and storage.
  - **utils**: Utility functions.
- **pkg/config/config.go**: Configuration code that can be shared across different parts of the application.
- **storage**: Directory where files are stored.


