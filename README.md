# Build and run

## 1. Customize [template(s)](./templates)

Read about the template syntax [here](https://github.com/vologzhan/maker)

Use the existing [go](./templates/go) template or replace it with a new

## 2. Build `maker_app` binary

```sh
make build
```

## 3. Move `maker_app` to dir with your projects

Projects are searched in executable path (in directory with file `maker_app`)

## 4. Run and use

Run binary

```sh
./maker_app
```

Open in browser http://localhost:1551

# Run dev

Clone this repository to directory with your projects and run:

```sh
make run-dev
```

Projects will be searched in path `../`
