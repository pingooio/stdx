.PHONY: fmt
fmt:
	cargo +nightly fmt


.PHONY: clean
clean:
	rm -rf target
