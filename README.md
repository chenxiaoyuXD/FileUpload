
# File Upload API with Chunked Uploads and Encryption

This project implements a RESTful API for uploading and downloading files using MinIO as the object storage backend. It supports chunked uploads with AES encryption to ensure secure file storage. This document explains the choices made during implementation, how to set up and run the project, and how to use the API.

---

## Features
1. **Chunked File Uploads**:
   - Files are uploaded in configurable chunk sizes (default: 5 MB).
   - Useful for handling large files efficiently without overloading memory.

2. **AES Encryption**:
   - Each file chunk is encrypted using AES-CFB mode before upload.
   - Encryption ensures the confidentiality of data at rest in MinIO.

3. **File Downloads**:
   - Encrypted files are downloaded from MinIO and decrypted before being returned to the client.

---

## Design Choices

1. **MinIO for Storage**:
   - MinIO is chosen for its lightweight, high-performance object storage capabilities.
   - Fully compatible with the S3 API.

2. **AES Encryption**:
   - AES-CFB mode is used for its simplicity and streaming support.
   - A static 32-byte key is used in this example for encryption/decryption, but in a production environment, use secure key management.

3. **Chunked Uploads**:
   - The file is split into configurable chunks for better scalability.
   - This prevents memory exhaustion when uploading large files.

4. **RESTful API**:
   - Designed with a simple and intuitive RESTful interface using Gin as the web framework.

5. **Configuration**:
   - The chunk size, MinIO credentials, and server settings are configurable using environment variables.

---

## Dependencies

- **Go**: Version 1.16 or later.
- **MinIO Go SDK**: v7.x.x (`github.com/minio/minio-go/v7`).
- **Gin Web Framework**: (`github.com/gin-gonic/gin`).

---

## Setup

### 1. Install Dependencies
Ensure you have Go installed. Then, clone the project and install dependencies:

```bash
git clone <repository-url>
cd <repository-directory>
go mod tidy
```

### 2. Set Up MinIO
Start MinIO using Docker Compose:

```yaml
version: '3'
services:
  minio:
    image: quay.io/minio/minio
    command:
      - server
      - /data
      - --console-address
      - ":9001"
    ports:
      - "9000:9000"
      - "9001:9001"
    environment:
      MINIO_ROOT_USER: minioadmin
      MINIO_ROOT_PASSWORD: minioadmin
```

Run the following command to start the MinIO service:

```bash
docker-compose up -d
```

Access the MinIO web console at [http://127.0.0.1:9001](http://127.0.0.1:9001).


### 3. Run the Project
Run the Go server:

```bash
go run main.go
```

---

## Usage

### Upload a File
Use the `/upload` endpoint to upload a file. The file will be uploaded in chunks and encrypted.

```bash
curl -X POST -F "file=@path/to/your/file.txt" http://localhost:8080/upload
```

**Response**:
```json
{
  "message": "File uploaded successfully in chunks"
}
```

### Download a File
Use the `/download/:filename` endpoint to download and decrypt a file.

```bash
curl -X GET http://localhost:8080/download/file.txt -o downloaded_file.txt
```

**Response**: The file will be saved as `downloaded_file.txt`.

---

## API Endpoints

### 1. Upload File
- **Endpoint**: `POST /upload`
- **Description**: Uploads a file to MinIO in chunks with AES encryption.
- **Form Data**:
  - `file`: The file to upload.

### 2. Download File
- **Endpoint**: `GET /download/:filename`
- **Description**: Downloads and decrypts a file from MinIO.
- **Path Parameter**:
  - `:filename`: The name of the file to download.

---

## Configuration

### Environment Variables
| Variable             | Default Value   | Description                                         |
|----------------------|-----------------|-----------------------------------------------------|
| `MINIO_ACCESS_KEY`   | `minioadmin`    | MinIO access key.                                  |
| `MINIO_SECRET_KEY`   | `minioadmin`    | MinIO secret key.                                  |
| `CHUNK_SIZE`         | `5242880`       | Chunk size in bytes (default: 5 MB).               |

---

## Security Considerations
1. **Encryption Key Management**:
   - Replace the static AES encryption key with a key management system (e.g., AWS KMS, HashiCorp Vault).
   - Never hardcode sensitive credentials in production.

2. **Environment Variables**:
   - Use `.env` files or secrets management tools to store environment variables securely.

3. **HTTPS**:
   - Configure the server to use HTTPS to protect data in transit.

---

## Limitations
- The current implementation assumes files are small enough to fit into memory when encrypting/decrypting.
- Error handling could be further enhanced for retrying failed chunk uploads.

---

## Future Enhancements
1. **Parallel Chunk Uploads**:
   - Optimize performance by uploading chunks in parallel.
2. **Retry Logic**:
   - Implement retry logic for failed chunk uploads.
