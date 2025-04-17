# s3cli

A CLI tool to manage S3-compatible stores, supporting multiple stores and credentials.

## Features
- Manage multiple S3 stores (add, list, delete, use, logout)
- Upload, download, list, and remove files and folders
- Uses `mc` (MinIO Client) under the hood

## Installation

1. Clone the repository:
   ```sh
   git clone https://github.com/kabbesgit/s3cli.git
   cd s3cli
   ```
2. Build the CLI:
   ```sh
   go build -o s3cli
   ```

## Usage

```sh
./s3cli --help
```

### Example Commands
- Add a store:
  ```sh
  ./s3cli store add --name mystore --endpoint https://s3.example.com --access-key AKIA... --secret-key ...
  ```
- Use a store:
  ```sh
  ./s3cli store use mystore
  ```
- Upload a file:
  ```sh
  ./s3cli put mystore mybucket/path/to/file.txt ./localfile.txt
  ```
- List contents:
  ```sh
  ./s3cli ls mystore mybucket/
  ```

## Contributing
Pull requests are welcome. For major changes, please open an issue first to discuss what you would like to change.

## License
[MIT](LICENSE)
