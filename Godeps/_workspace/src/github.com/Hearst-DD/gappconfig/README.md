## GappConfig: Simple environment config loader

GappConfig is a simple way to load config parameters from environment variables. It is designed to be invoked by a service at startup. 

## Install

```bash
go get github.com/Hearst-DD/gappconfig
```

## Use

```go
import "github.com/Hearst-DD/gappconfig"

func main() {
	// prefix causes gappconfig to load all vars with the given prefix, e.g. MYSERVICE_ENV
	config = gappconfig.New("MYSERVICE_", gappconfig.Map{
		{"ENV", "dev"},
		{"PRETTY", true},
		{"LOG_LEVEL", "DEBUG"},
		{"HOST", "localhost"},
		{"PORT", 4001},
		{"SERVER_TIMEOUT", time.Second*30},
		// ...
	})

	host := config.String("HOST")
	port := config.Int("PORT")
	timeout := config.Duration("SERVER_TIMEOUT")
	prettyPrint := config.Bool("PRETTY")
}
```