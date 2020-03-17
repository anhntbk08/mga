# A Self-Documenting Makefile: http://marmelab.com/blog/2016/02/29/auto-documented-makefile.html

# Add the ability to override some variables
# Use with care
-include override.mk

# Main targets
include main.mk

# Add custom targets here
-include custom.mk

.PHONY: generate
generate: export PATH := $(abspath ${BUILD_DIR}/):${PATH}
generate: build ## Generate test code
	go generate ./...
	${BUILD_DIR}/mga generate kit endpoint ./...
	${BUILD_DIR}/mga generate event handler ./...
	${BUILD_DIR}/mga generate event handler --output subpkg:suffix=gen ./...
	${BUILD_DIR}/mga generate event dispatcher ./...
	${BUILD_DIR}/mga generate event dispatcher --output subpkg:suffix=gen ./...
	${BUILD_DIR}/mga generate testify mock ./...
	${BUILD_DIR}/mga generate testify mock --output subpkg:suffix=mocks ./...
	${BUILD_DIR}/mga create service --force internal/scaffold/service/test
