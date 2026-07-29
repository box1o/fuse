.PHONY: link unlink install

RUNE_LINK_DIR ?= $(HOME)/.local/bin


link:
	mkdir -p "$(RUNE_LINK_DIR)"
	ln -sfn "$(CURDIR)/rune" "$(RUNE_LINK_DIR)/rune"

unlink:
	@if [ -L "$(RUNE_LINK_DIR)/rune" ]; then unlink "$(RUNE_LINK_DIR)/rune"; fi

install: link
	@printf '%s\n' 'Development link installed. Ensure $(RUNE_LINK_DIR) is on PATH.'
