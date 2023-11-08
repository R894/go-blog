# Blog API in Go

This is a Go-based RESTful API for managing a blog. It allows you to create, read, update, and delete blog posts.

## Table of Contents

- [Features](#features)
- [Installation](#installation)
- [Usage](#usage)
- [API Endpoints](#api-endpoints)

## Features

- Create, Read, Update, and Delete blog posts.
- Authentication and authorization.
- Data validation and error handling.
- Focus on a modular and well-documented codebase.

## Installation

1. Clone the repository:

```
git clone https://github.com/R894/go-blog
cd  go-blog
```
2. Install the required dependencies
```
go mod tidy
```
3. Configure the environment variables. Create a `.env` file and specify the `PORT` and `JWT_SECRET`
4. Build the API
```
go build
```
5. Run the API
```
./go-blog
```

## Usage
The API is designed to be straightforward to use. Here's a brief overview of how to interact with it:

* Make HTTP requests to the provided endpoints (documented below) using a tool like `curl` or `Postman`.
* Ensure that you include the necessary headers and authentication tokens for protected routes.

## API Endpoints

### Posts
* `GET /posts`: Get a list of all blog posts.
* `GET /posts/:id`: Get a specific blog post by ID.
* `POST /posts`:  Create a new blog post.
* `PUT /posts/:id`: Update an existing blog post.
* `DELETE /posts/:id`: Delete a blog post.

### Comments
* `GET /comments/:id`: Get comments by post ID
* `POST /comments/:id`: Create a new comment and associate it with post ID.

### Authentication
* `POST /register`: Register new user
* `POST /login`: User login
