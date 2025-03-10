# grpc-go
 gRPC microservice in Go


Core Requirements
1. CRUD Operations
   Accounts: Manage user credentials (username, email, hashed password), status lifecycle (ACTIVE/SUSPENDED), and audit timestamps.
   Preferences: Handle user settings (theme, notifications, locale) with cascading referential integrity.
   Roles: Implement RBAC (Role-Based Access Control) with permission scopes (e.g., orders:write) and default role assignments.
   Orders: Support transactional data, including an itemized product list, automated total validation, and status workflows.
2. File Management
   Implement chunked file uploads (max 5MB per chunk) for:
   Account avatars
   Order attachments (PNG, JPEG, PDF)
   Store metadata in PostgreSQL or MySQL.
3. Validation & Security
   Ensure proper email validation and password security.
   Follow best practices for data validation, authentication, and authorization.
4. Unit Testing
   Implement unit tests only for some core logic.
5. Error Handling & Logging