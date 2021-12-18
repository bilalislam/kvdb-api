# Key-Value Rest Api




# Table of contents

- [Table of contents](#table-of-contents)
- [Installation](#installation)
- [Usage](#usage)
- [Known Issues](#known-issues)
- [Architecure](#architecure)
- [References](#references)

# Installation
[(Back to top)](#table-of-contents)

To use this project, first clone the repo on your device using the command below:

```git init```

```git clone https://github.com/bilalislam/kvdb-api```


# Usage
[(Back to top)](#table-of-contents)

### `build`

```sh
$ cd kvdb
$ go build main.go -o kvdb-api
```

### `api start`

```sh
$ ./kvdb-api
```

Runs the app in the development mode.<br />
Open [http://localhost:8080](http://localhost:8080) to view it in the browser.

The page will reload if you make edits.<br />
You will also see any lint errors in the console.

### `unit test`

```sh
$ go test
```

# Deployment

[(Back to top)](#table-of-contents)

### `docker compose for local build`

Compose is a tool for defining and running multi-container Docker applications. With Compose, you use a YAML file to configure your application’s services. Then, with a single command, you create and start all the services from your configuration. To learn more about all the features of Compose, [see the list of features ](https://docs.docker.com/compose/#features)

A docker-compose.yml looks like:

```docker

version: '3'

services:
  api:
    build:
      context: .
      dockerfile: tools/Dockerfile
    ports:
      - "8080:8080"

```

```sh
$ docker-compose up
```

Open [http://localhost:8080](http://localhost:8080) to view it in the browser.


# Known Issues

[(Back to top)](#table-of-contents)


# Architecure

[(Back to top)](#table-of-contents)


# References

[(Back to top)](#table-of-contents)


