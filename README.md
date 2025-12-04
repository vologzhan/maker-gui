# Build and run

## 1. Customize template(s)

Use the existing [go template](./templates/go) or use it as an example and replace it with your own template

Read about the [template syntax](https://github.com/vologzhan/maker)

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

Projects will be searched in path `../`

Clone this repository to directory with your projects and run:

```sh
make run-dev
```
