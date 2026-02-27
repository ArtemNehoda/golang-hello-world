# Go GraphQL hello world

A GraphQL hello world service. The service allows you to create and retrieve messages, which are stored in a MySQL database. The service is built using Go, GraphQL, and Docker. Architecture is layered, with a clear separation of concerns between the different layers.

### Quick Start

```bash
cd golang
docker-compose up -d
```

Then server will be available at:

```
http://localhost:8080/
```

---

### GraphQL Endpoint

**`POST /graphql`**

I suggest to use postman or similar to any one who will be testing this service.

### Query example:

```graphql
query {
  messages {
    id
    content
    author
    createdAt
  }
}
```
