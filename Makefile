# Main Makefile - Includes all sub-makefiles
.PHONY: help

# Include all makefiles
include Makefile.dev
include Makefile.test
include Makefile.docker
include Makefile.utils

# Default target
help:
	@echo "Go HTTP Server - Makefile Commands"
	@echo ""
	@echo "Use category-specific help commands:"
	@echo "  make dev.help      - Development commands (dev.*)"
	@echo "  make tests.help    - Testing commands (tests.*)"
	@echo "  make app.help      - Docker commands (app.*)"
	@echo "  make utils.help    - Utility commands (utils.*)"
	@echo ""
	@echo "Or run 'make' with any target from the included makefiles."
	@echo ""
	@echo "Examples:"
	@echo "  make dev.run       - Run the application locally"
	@echo "  make app.build     - Build Docker containers"
	@echo "  make tests.test    - Run all tests"
	@echo "  make utils.clean   - Clean build artifacts"
