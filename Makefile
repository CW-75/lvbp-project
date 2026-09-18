.PHONY: tools-install graphify

tools-install:
	python3 -m venv .venv
	.venv/bin/pip install -r requirements-dev.txt

graphify:
	@echo "Generating knowledge graph..."
	if [ -f ".venv/bin/graphify" ]; then .venv/bin/graphify update .; else echo "graphify not found, run 'make tools-install'"; exit 1; fi
