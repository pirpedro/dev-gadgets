PREFIX ?= /usr/local
BINPREFIX ?= "$(PREFIX)/bin"

SYSCONFDIR ?= $(PREFIX)/etc
BINS = $(wildcard bin/git-*)
CODE_DIR := $(dir $(abspath $(lastword $(MAKEFILE_LIST))))
INSTALL_VIA ?= source
# Libraries used by all commands
LIB = "helper/reset-env" "helper/dev-gadgets-utils"

COMMANDS = $(subst bin/, , $(BINS))

default: build-go

check:
	@echo "Check dependencies before installation"
	@./installation/check_dependencies.sh
	@echo

install: check
	@mkdir -p $(DESTDIR)$(BINPREFIX)
	@echo "... installing bins to $(DESTDIR)$(BINPREFIX)"
	$(eval TEMPFILE := $(shell mktemp -q $${TMPDIR:-/tmp}/dev-gadgets.XXXXXX 2>/dev/null || mktemp -q))
	@# chmod from rw-------(default) to rwxrwxr-x, so that users can exec the scripts
	@chmod 775 $(TEMPFILE)
	$(eval EXISTED_ALIASES := $(shell \
		git config --get-regexp 'alias.*' | awk '{print "git-" substr($$1, 7)}'))
	@$(foreach COMMAND, $(COMMANDS), \
		disable=''; \
		if test ! -z "$(filter $(COMMAND), $(EXISTED_ALIASES))"; then \
			read -p "$(COMMAND) conflicts with an alias, still install it and disable the alias? [y/n]" answer; \
			test "$$answer" = 'n' -o "$$answer" = 'N' && disable="true"; \
		fi; \
		if test -z "$$disable"; then \
			echo "... installing $(COMMAND)"; \
			head -1 bin/$(COMMAND) > $(TEMPFILE); \
			cat $(LIB) >> $(TEMPFILE); \
			if grep "$(COMMAND)" need_git_repo >/dev/null; then \
				cat ./helper/is-git-repo >> $(TEMPFILE); \
			fi; \
			if grep "$(COMMAND)" need_git_commit >/dev/null; then \
				cat ./helper/has-git-commit >> $(TEMPFILE); \
			fi; \
			tail -n +2 bin/$(COMMAND) >> $(TEMPFILE); \
			cp -f $(TEMPFILE) $(DESTDIR)$(BINPREFIX)/$(COMMAND); \
		fi; \
	)
	@mkdir -p $(DESTDIR)$(SYSCONFDIR)/bash_completion.d
	cp -f contrib/completion/bash_completion.sh $(DESTDIR)$(SYSCONFDIR)/bash_completion.d/dev-gadgets
	@echo ""
	@echo "If you are a zsh user, you may want to 'source $(CODE_DIR)contrib/completion/dev-gadgets-completion.zsh'" \
		"and put this line into ~/.zshrc to enable zsh completion"

uninstall:
	@$(foreach BIN, $(BINS), \
		echo "... uninstalling $(DESTDIR)$(BINPREFIX)/$(notdir $(BIN))"; \
		rm -f $(DESTDIR)$(BINPREFIX)/$(notdir $(BIN)); \
	)
	rm -f $(DESTDIR)$(SYSCONFDIR)/bash_completion.d/dev-gadgets

.PHONY: default check install uninstall

# --- Go build (new path) ---
GO_CMD?=go
BIN_DIR?=dist

.PHONY: build-go test-go run-go
build-go:
	$(GO_CMD) build -o $(BIN_DIR)/dev-gadgets ./cmd/dev-gadgets

test-go:
	$(GO_CMD) test ./...

run-go:
	$(GO_CMD) run ./cmd/dev-gadgets --help
