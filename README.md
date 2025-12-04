# Build and run

## 1. Use exists, add new or replace template(-s) in directory [templates](./templates)

Read about the template syntax [here](https://github.com/vologzhan/maker)

See examples [go](./templates/go) and [maker](./templates/maker) templates

**Note:** `.dotfiles` must be specified in `go:embed` in [main.go](./main.go)

## 2. Build binary `maker_app`

```sh
make build
```

## 3. Move `maker_app` to dir with your projects

Projects are searched in the executable path

## 4. Run and use

Run binary

```sh
./maker_app
```

Open browser http://localhost:1551

# Build dev

Clone this repository to dir with your projects

```sh
make build-dev
```

Projects will be searched in path `../`
