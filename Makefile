FRONTEND_DIR=frontend
BUILD_BACKEND=go build -o maker_app

build:
	${BUILD_BACKEND}

build-dev:
	${BUILD_BACKEND} -tags dev

build-front:
	cd $(FRONTEND_DIR) && npm run build

run-front:
	cd $(FRONTEND_DIR) && npm run dev

npm-install:
	cd $(FRONTEND_DIR) && rm -rf node_modules && npm install
