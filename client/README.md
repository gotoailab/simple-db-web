# SimpleDBWeb Client

This is a complete SimpleDBWeb client version with user authentication and management features.

## Features

- ✅ User login and authentication
- ✅ Session management (stored in SQLite)
- ✅ User management (admins can create, edit, and delete users)
- ✅ Password change functionality
- ✅ Logout functionality
- ✅ User management UI injected via `SetCustomScript` (non-intrusive to core code)

## Tech Stack

- **Web Framework**: Gin
- **Database**: modernc.org/sqlite (for storing users and sessions)
- **Authentication**: Session-based authentication
- **Password Encryption**: bcrypt

## Installation and Running

### 1. Install Dependencies

```bash
cd client
go mod tidy
```

### 2. Run Application

The application supports the following command-line arguments:

#### Basic Parameters

- `-port` (default: `:8080`): Server port
  - Example: `-port :8080` or `-port 8080`
  
- `-debug` (default: `false`): Enable debug mode
  - When enabled, detailed debug logs will be printed
  - Example: `-debug`
  
- `-log` (default: empty): Log file path
  - If specified, logs will be written to the file
  - Example: `-log ./logs/app.log`
  
- `-auth` (default: `false`): Enable authentication and user management
  - When enabled, login interface and user management features will be available
  - Example: `-auth`
  
- `-prefix` (default: empty): Route prefix
  - All routes will have this prefix added
  - Example: `-prefix /v1` or `-prefix /api`
  
- `-open` (default: `false`): Automatically open browser after startup
  - Example: `-open`
  
- `-db` (default: `client.db`): Database file path
  - Only used when authentication is enabled
  - Example: `-db ./data/client.db`

#### Usage Examples

```bash
# Basic usage (no authentication)
go run main.go

# Enable authentication and user management
go run main.go -auth -db ./data/client.db

# Custom port and route prefix
go run main.go -port :9000 -prefix /v1

# Enable debug mode and log file
go run main.go -debug -log ./logs/app.log

# Full configuration example
go run main.go \
  -port :8080 \
  -debug \
  -log ./logs/app.log \
  -auth \
  -prefix /v1 \
  -open \
  -db ./data/client.db
```

### 3. Access Application

Open your browser and visit: `http://localhost:8080`

Default admin account:
- Username: `admin`
- Password: `admin123`

**Important**: Please change the default admin password immediately after first run!

## Project Structure

```
client/
├── main.go              # Main program entry
├── db.go                # Database initialization and table structure
├── auth.go              # User authentication and session management
├── middleware.go        # Authentication middleware
├── handlers.go          # HTTP request handlers
├── go.mod               # Go module definition
├── templates/           # HTML templates
│   └── login.html       # Login page
└── static/              # Static files
    └── user-management.js  # User management UI (injected via embed)
```

## API Endpoints

### Authentication

- `POST /api/auth/login` - User login
- `POST /api/auth/logout` - Logout
- `GET /api/auth/current` - Get current user information
- `POST /api/auth/password` - Change password

### User Management (Admin Only)

- `GET /api/users` - Get all users list
- `POST /api/users` - Create new user
- `PUT /api/users/:id` - Update user information
- `DELETE /api/users/:id` - Delete user

## Database Structure

### users Table

| Field | Type | Description |
|-------|------|-------------|
| id | INTEGER | Primary key, auto-increment |
| username | TEXT | Username (unique) |
| password_hash | TEXT | Password hash value |
| is_admin | INTEGER | Is administrator (0/1) |
| created_at | DATETIME | Creation time |
| updated_at | DATETIME | Update time |

### sessions Table

| Field | Type | Description |
|-------|------|-------------|
| id | INTEGER | Primary key, auto-increment |
| session_id | TEXT | Session ID (unique) |
| user_id | INTEGER | User ID (foreign key) |
| username | TEXT | Username |
| created_at | DATETIME | Creation time |
| expires_at | DATETIME | Expiration time |

## Security Notes

1. **Password Encryption**: Uses bcrypt for password hashing
2. **Session Expiration**: Sessions expire after 24 hours by default
3. **Cookie Security**: Session ID is stored in HttpOnly Cookie
4. **Permission Control**: User management features are only accessible to administrators

## Notes

1. Default admin account (admin/admin123) will be automatically created on first run
2. Database file will be automatically created if it doesn't exist
3. Expired sessions are automatically cleaned up periodically (every hour)
4. User management UI is injected via `SetCustomScript`, without modifying core project code

