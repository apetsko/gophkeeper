# GophKeeper

GophKeeper is a secure backend service written in Go for storing and managing sensitive user data, such as passwords, notes, and files. It is designed with a focus on security, scalability, and modularity.

## Features

- **User Authentication:**  
  Stateless authentication using JWT (JSON Web Tokens). Users sign up and log in to receive a token for API access.

- **gRPC API:**  
  All operations are exposed via a gRPC server, including:
  - `Ping` (health check)
  - `Login` and `Signup`
  - `DataList`, `DataSave`, `DataDelete`, `DataView` for managing user data

- **Secure Data Storage:**  
  - User data is encrypted before storage.
  - Supports integration with S3-compatible storage (e.g., MinIO) for files and large objects.

- **Password Security:**  
  - Passwords are hashed using bcrypt before storage.
  - Password verification is performed securely.

- **ID Generation:**  
  - Unique IDs for stored objects are generated using SHA-256 and base64 encoding.

- **Validation:**  
  - Struct validation is performed using `go-playground/validator` to ensure data integrity.

## Architecture

- **Modular Design:**  
  The codebase is organized into packages for configuration, storage, server logic, authentication, utilities, and mocks for testing.

- **Testing:**  
  Comprehensive unit tests cover all core logic, including authentication, password handling, gRPC handlers, and utility functions.

- **Configuration:**  
  Supports configuration for JWT secrets, S3 credentials, and gRPC server settings.

## Security

- **JWT Middleware:**  
  All protected gRPC endpoints require a valid JWT, enforced by a middleware interceptor.

- **Password Hashing:**  
  Passwords are never stored in plain text; bcrypt is used for hashing.

- **Data Encryption:**  
  Sensitive data is encrypted before storage, using envelope encryption.

## Tech Stack

- **Language:** Go
- **API:** gRPC
- **Storage:** MinIO/S3-compatible
- **Password Hashing:** bcrypt
- **Authentication:** JWT

## Intended Use

GophKeeper is suitable for applications where users need to securely store, retrieve, and manage sensitive data with strong authentication and encryption.