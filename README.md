# Dropit

DropIt is a lightweight file-sharing service written in Go. It allows
quick uploads and sharing of small files without the need for sign-up or
login. Users can also optionally create accounts to manage and delete
their files.

## Features

-   **Anonymous file upload**
    Upload files without creating an account.
-   **Direct file sharing**
    Get a shareable token to let anyone access your file.
-   **Optional user accounts**
    Register with a simple username and password (no email confirmation
    required).
-   **User file management**
    Authenticated users can delete their files or list their uploaded
    files.

## Endpoints

### Files

-   `POST /files`
    Upload a file anonymously or as an authenticated user.
    Returns a unique file token.

-   `GET /files/:token`
    Download a file by its token.
    Anyone with the token can access the file.

-   `DELETE /files/:token` (requires authentication)
    Delete a file owned by the authenticated user.

### Users

-   `POST /users`
    Register a new user with a username and password.

-   `GET /users/:id/files`
    Get a list of all files uploaded by a user.
    Since file lists are public, you can share your user ID with others
    to act as a "shared group folder".

### Authentication

-   `POST /auth/login`
    Log in with username and password.
    Returns a token that can be used to:
    -   Delete your files
    -   Upload files under your account

## Usage Examples

### Upload a file anonymously

``` bash
curl -X POST http://localhost:8000/files   -F "file=@example.txt"
```

### Get a file by token

``` bash
curl -X GET http://localhost:8000/files/{token} -O
```

### Register a new user

``` bash
curl -X POST http://localhost:8000/users   -H "Content-Type: application/json"   -d '{"username":"alice","password":"mypassword"}'
```

### Login

``` bash
curl -X POST http://localhost:8000/auth/login   -H "Content-Type: application/json"   -d '{"username":"alice","password":"mypassword"}'
```

Response includes a token to use in `Authorization` headers.

### Upload a file as a user

``` bash
curl -X POST http://localhost:8000/files   -H "Authorization: Bearer <token>"   -F "file=@example.txt"
```

### Delete a file

``` bash
curl -X DELETE http://localhost:8000/files/{token}   -H "Authorization: Bearer <token>"
```

### Get all files of a user

``` bash
curl -X GET http://localhost:8000/users/{id}/files
```

## Notes

-   Anonymous files are public --- anyone with the token can access
    them.
-   User accounts allow more control, but file lists are still public by
    design.
-   DropIt is designed to be **fast, simple, and minimal**.

## License

MIT
