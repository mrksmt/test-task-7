# Test task for Server engineer (Go)

## another header

### Usage

Run all-in-one by

```bash
make run-all-in-one
```

### Build

Build docker containers by

``` bash
make docker-build-client
make docker-build-server
```

### Diagram

yes, pictures are fun

```mermaid
sequenceDiagram

    actor Client
    Client->>Server: establish tcp connection 
    Server->>Client: challenge code
    loop work
        Client->>Client: try to find valid md5 hash
    end
    Client->>Server: challenge response
    Server->>Server: response validating
    Server->>Client: sentence
```
