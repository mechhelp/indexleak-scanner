# Test Server Setup

This document explains how to set up and run a local Nginx test server that serves mock sensitive data for testing the `indexleak` scanner.

The server runs as a Docker container, configured via the `Dockerfile` and `nginx.conf` files in this directory. The test data is located in the `zPeter` directory.

### Requirements

- [Docker](https://www.docker.com/get-started) must be installed.

### Running the Server

1.  **Build the Docker Image:**
    From this directory (`test-data`), build the Docker image using the following command:
    ```bash
    docker build -t indexleak-scanner-test-server .
    ```

2.  **Start the Container:**
    Run this command to start the server:
    ```bash
    docker run -d -p 8080:80 --name indexleak-server indexleak-scanner-test-server
    ```

3.  **Verify:**
    To confirm the setup was successful, visit [http://localhost:8080/zPeter/](http://localhost:8080/zPeter/) in your web browser. You should see the contents of the `zPeter` directory.

### Stopping the Server

You can stop and remove the test server container with the following command:
```bash
docker stop indexleak-server && docker rm indexleak-server
```