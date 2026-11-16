# bloom

Small Go utilities shared across meerkat projects.

```sh
go get github.com/meerkat-io/bloom
```

Requires Go 1.25+.

| Package | Description |
|---|---|
| `flag` | Parse command line args or environment variables into a struct via tags |
| `folder` | `Exists` and `Create` (recursive) for directories |
| `logger` | Leveled, colored logger with caller position |
| `strcase` | Convert strings to snake_case, chain-case, camelCase, PascalCase |
| `time` | Unix-second helpers: `Timestamp`, `Unix`, `Elapse`, `Expired` |

## Examples

```go
// flag
type Flags struct {
	Help   bool   `flag:"h" usage:"help"`
	Input  string `flag:"i" usage:"input folder" tip:"input" default:"."`
	Lang   string `flag:"t" usage:"language" required:"true"`
}
var f Flags
err := flag.ParseCommandLine(&f) // or flag.ParseEnvironment(&f)

// strcase
strcase.SnakeCase("HelloWorld") // hello_world
strcase.CamelCase("hello_world") // helloWorld

// logger
logger.SetLevel(logger.INFO)
logger.Infof("started %s", name)
```

Supported `flag` field types: `bool`, `string`, `int`, `float64`.

## Test

```sh
go test -race -cover ./...
```

## License

[MIT](LICENSE)
