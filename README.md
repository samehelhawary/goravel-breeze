# goravel-breeze
This projects change goravel project default routing from gin to fiber and use jet templating engine instead of default html templating engine 

## Follow the following command to create goravel project and install goravel-breeze

```bash
git clone --depth=1 https://github.com/goravel/goravel.git goravel-test && rm -rf goravel-test/.git*
```


```bash
cd goravel-test && go mod tidy && cp .env.example .env && go run . artisan key:generate
```


Add the following line to go.mod
```
replace github.com/samehelhawary/goravel-breeze => ../goravel-breeze
```

```bash
go get github.com/samehelhawary/goravel-breeze
```

**in config/app.go**
```go
import breeze "github.com/samehelhawary/goravel-breeze"

&breeze.ServiceProvider{}, // In providers map
```

```bash
go mod tidy 
```

Then add mysql connection in .env file (Optional)

```bash
go run . artisan breeze:install
```

go run .