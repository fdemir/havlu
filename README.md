# havlu

Havlu lets you focus on the frontend by making it effortless to have a custom REST. You give the JSON, it gives you the REST API. Fast, easy, and reliable. You can self-host it, easily deploy it to somewhere, or use it as a Docker container.

<!-- One of the amazing things about havlu is that it can create a custom REST API with a model that you already have. -->

## Motivation

When you are developing a frontend application, you need a backend to serve your data. You can use a mock server, but it is not always enough. You may need to have a custom REST API. Havlu is a tool that helps you to create a custom REST API with a JSON file.

However, havlu is not offering a new approach for existing tools like `json-server`. It does the same thing but faster. Why not?

## Installation

<!-- MacOS

```bash
brew install havlu
``` -->

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
| --file    | Path to the JSON file          | data.json |
| --port    | Port number                    | 3000      |
| --host    | Host name                      | localhost |
| --cors    | Enable CORS                    | false     |
| --delay   | Response delay in milliseconds | 0         |
| --help    | Show help                      |           |
| --version | Show version number            |           |

## JSON File

The JSON file should be an array of objects. Each object represents a resource. The key of the object is the resource name. The value of the object is an array of objects. Each object represents a resource item. The key of the object is the resource item id. The value of the object is the resource item.

`db.json`

```json
[
  {
    "users": [
      {
        "1": {
          "id": 1,
          "name": "John Doe"
        },
        "2": {
          "id": 2,
          "name": "Jane Doe"
        }
      }
    ]
  }
]
```

Run havlu with the following command.

```bash
havlu db.json
```

## Using as Module

You can adapt havlu to your own server by using it as a module.

```go
todo!
```

## Examples

```json
{
  "locations": [
    {
      "city": "New York"
    },
    {
      "city": "Izmir"
    },
    {
      "city": "Oslo"
    },
    {
      "city": "Beijing"
    },
    {
      "city": "Tokyo"
    }
  ]
}
```

```bash
havlu --file data.json --port 3000
```

It will create a REST API with the following endpoints.

```bash
GET /locations
GET /locations?city=Tokyo

```

<!-- GET /locations?order=city&sort=desc -->

## License

MIT
