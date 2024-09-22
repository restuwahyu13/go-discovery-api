#################################
# Application Territory
#################################
NODEMON := @nodemon

.PHONY: dev
dev:
	${NODEMON} -x go run ./cmd/api . --signal SIGKILL -e go --verbose