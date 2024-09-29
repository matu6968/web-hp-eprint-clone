# S3 Client Web

This is a WebGUI for [s3-client](https://git.fluffy.pw/leafus/s3-client)

## Prerequisites

- Go (1.23.1 or later)

## Installation

1. Clone the repository:
   ```
   git clone https://git.fluffy.pw/leafus/s3-client-web
   ```

2. Go to the project directory:
   ```
   cd s3-client-web
   ```

3. Build the binary:
   ```
   go build -o s3-client-web
   ```

## Configuration

In the .env file this is the only thing you can set

```
PORT=8080
```

### For this to even work

You need to download a latest linux binary release of [s3-client](https://git.fluffy.pw/leafus/s3-client) from the "Releases" tab for your architecture

and put the binary in the "bin" folder !IMPORTANT! you need to rename the file to just s3-client and create s3config.toml file, the configuration of that is in the [s3-client](https://git.fluffy.pw/leafus/s3-client) repository