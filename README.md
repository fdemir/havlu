# havlu

Havlu lets you focus on the frontend by making it effortless to have a custom REST. You give the JSON, it gives you the REST API. Fast, easy, and reliable. You can self-host it, easily deploy it to somewhere, or use it as a Docker container. Everything happening on the memory, so you don't need to worry about the database.

## Motivation

When you are developing a frontend application, you need a backend to serve your data. You can use a mock server, but it is not always enough. You may need to have a custom REST API. Havlu is a tool that helps you to create a custom REST API with a JSON file.

However, havlu is not offering a new approach for existing tools like `json-server`. It does the same thing but faster. Why not?

## Installation

Mac OS

```bash
brew install fdemir/tap/havlu
```

Go

```bash
go install github.com/fdemir/havlu@latest
```

Source

```bash
git clone
cd havlu
go install
```

## Usage

```bash
havlu data.json --port 3000
```

## Options

| Option    | Description                    | Default   |
| --------- | ------------------------------ | --------- |
| --port    | Port number                    | 3000      |
| --host    | Host name                      | localhost |
| --cors    | Enable CORS                    | false     |
| --delay   | Response delay in milliseconds | 0         |
| --help    | Show help                      |           |
| --version | Show version number            |           |

## Hav File

Hav is a simple schema file that defines the entities(resources) and their fake data types. It helps you create mock APIs faster. Here is an example:
```
entity user {
  name Person.Name
  email Internet.Email
}

entity address {
  lat Address.Latitude
  lon Address.Longitude
}
```

## JSON File

The JSON file should be an array of objects. Each object represents a resource. The key of the object is the resource name. The value of the object is an array of objects. Each object represents a resource item. The key of the object is the resource item id. The value of the object is the resource item.

`db.json`

```json
[
  {
    "users": [
      {
        "id": 1,
        "name": "Furkan Demir",
        "gender": "G",
      },
      {
        "id": 2,
        "name": "John Doe",
        "gender": "B"
      }
    ]
  }
]
```

Run havlu with the following command.

```bash
havlu db.json
```

It will create a REST API for the `users` resource.

```bash
GET    /users
POST /users
DELETE /users/2
GET /users?gender=B
```

## Using as Module

You can adapt havlu to your own server by using it as a module.

```go
todo!
```

<!-- GET /locations?order=city&sort=desc -->

## License

MIT
