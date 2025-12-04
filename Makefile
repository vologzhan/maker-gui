FRONTEND_DIR=frontend

build:
	go build -o maker_app

run-dev:
	go run -tags dev .

build-front:
	cd $(FRONTEND_DIR) && npm run build

run-front:
	cd $(FRONTEND_DIR) && npm run dev

npm-install:
	cd $(FRONTEND_DIR) && rm -rf node_modules && npm install
