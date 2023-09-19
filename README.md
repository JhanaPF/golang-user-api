# golang-user-api

[![MIT License](https://img.shields.io/badge/License-MIT-green.svg)](https://choosealicense.com/licenses/mit/)

User api for my learning language app

## Table of Contents

- [Installation](#installation)
- [Contributing](#contributing)
- [License](#license)


## Installation

Install Golang on your machine ( https://golang.org/dl/ for Windows )
```bash
    sudo apt update
    sudo apt install golang

    git clone https://github.com/JhanaPF/golang-user-api.git
    cd golang-user-api
    go mod download
```

#### Postgresql

Official documentation : [https://www.postgresql.org/download/](https://www.postgresql.org/download/).

Note username and password !

On linux:
```bash
    sudo apt update
    sudo apt install postgresql-15 
    sudo systemctl status postgresql
    sudo systemctl start postgresql
```

Execute sql file:
```bash
    psql -U username -d myDataBase -a -f info.sql
```

Or connect to sql shell to create database:
```bash

    sudo -u postgres psql // postgres is superuser
    now you are in sql shell :
    CREATE USER nouvel_utilisateur WITH PASSWORD 'mot_de_passe';
    GRANT ALL PRIVILEGES ON DATABASE ma_base_de_donnees TO nouvel_utilisateur;
    SELECT usename FROM pg_user;
    CREATE DATABASE myDatabase WITH ENCODING 'UTF8';
    \du to check databases list
```

## Run Locally

```bash
    go run main.go
```

## Build

```bash
    go build
```

## Tech Stack

Golang
GIN
Postgresql
Redis

## Environment Variables

To run this project, you will need to add the following environment variables to your .env file

`DICTIONNARY_API_TOKEN`
`DICTIONNARY_API_USERID`
`DICTIONNARY_API_USERPASSWORD`
`PORT`


## API Reference

#### Get item

```http
  GET /items/${id}
```

| Parameter | Type     | Description                       |
| :-------- | :------- | :-------------------------------- |
| `id`      | `string` | **Required**. Id of item to fetch |
