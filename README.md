<div align="center">
  <a href="https://raw.githubusercontent.com/yogawahyudi7/hash-tag/develop/doc/hashtag.webp">
    <img src="https://raw.githubusercontent.com/yogawahyudi7/hash-tag/develop/doc/hashtag.webp">
  </a>
</div>
<div>
#  HASH-TAG

[![Go Reference](https://pkg.go.dev/badge/golang.org/x/example.svg)](https://pkg.go.dev/golang.org/x/example)

## Table of Content

- [HASHTAG:CLUB](#hash-tag)
  - [Table of Content](#table-of-content)
  - [Features](#features)
    - [Endpoints](#endpoints)
    - [API Documentation](#api-documentation)
  - [System Design](#system-design)
    - [ERD](#erd)
    - [Layered Architecture](#layered-architecture)
  - [Getting Started](#getting-started)
    - [Installing](#installing)
  - [Authors](#authors)

## Features

- JWT Authentication
- Layered Architecture
- Dependency Injection
- Multi Role Middleware (user, admin)

### Endpoints

- [x] Register
- [x] Login
- [x] Add Post
- [x] Get All Post
- [x] Update Post
- [x] Get Post By Id
- [x] Get Post By Tag
- [x] Delete Post
- [x] Pulbish Post (admin only)

### API Documentation

![API Documentation](https://raw.githubusercontent.com/yogawahyudi7/hash-tag/develop/doc/postman.png)
Application Programming Interface is available at [POSTMAN.](https://documenter.getpostman.com/view/16411992/2sA3QqfsDc)

## System Design

### ERD

![HashTag- ERD](https://raw.githubusercontent.com/yogawahyudi7/hash-tag/develop/doc/erd.png)

### Layered Architecture

![HashTag - Layered Architecture](https://raw.githubusercontent.com/yogawahyudi7/hash-tag/develop/doc/tree.png)

## Getting Started

Below we describe how to start this project

### Installing

You must download and install `Golang`, follow [this instruction](https://golang.org/doc/install) to install.
And `Postgresql` as database [this instruction](https://www.postgresql.org/download/).

After Golang installed, Follow this instructions

```bash
$ git clone https://github.com/yogawahyudi7/hash-tag
$ cd hash-tag
$ .cmd/db-generator.sh
$ go run main.go
```

Go to `http://localhost:9000/` to [start this application.](http://localhost:9000/)

## Authors

- [@yogawahyudi7](https://github.com/yogawahyudi7) - Developer
