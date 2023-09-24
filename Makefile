SHELL := /bin/bash

SPEAKER_NOTES ?=
CLOUDFLARE_WA_TOKEN ?=
ASCIIDOCTOR = bundle exec asciidoctor-revealjs
ASCIIDOCTOR_ARGS = \
	   -r asciidoctor-diagram \
	   -r ./lib/analytics.rb \
	   -a imagesdir=images \
	   -a docinfo=shared \
	   -a pikchr="$(shell pwd)/scripts/pikchr.sh" \
	   -t

ifneq ($(SPEAKER_NOTES),)
ASCIIDOCTOR_ARGS += -a revealjs_showNotes=separate-page
endif

ifneq ($(CLOUDFLARE_WA_TOKEN),)
ASCIIDOCTOR_ARGS += -a cloudflare-wa-token=$(CLOUDFLARE_WA_TOKEN)
endif

.PHONY: site
site: third_party
	@rm -rf _site && mkdir -p _site/reveal.js/plugin
	$(ASCIIDOCTOR) -D _site $(ASCIIDOCTOR_ARGS) index.adoc
	rm -rf _site/.asciidoctor
	cp -R css _site/css
	cp -R fonts _site/fonts
	cp -R webfonts _site/webfonts
	cp -r images/* _site/images
	cp -r highlight _site/highlight
	cp -R reveal.js/dist _site/reveal.js/dist
	cp -R reveal.js/plugin/{highlight,notes} _site/reveal.js/plugin

.PHONY: build-dev
build-dev: third_party
	$(ASCIIDOCTOR) $(ASCIIDOCTOR_ARGS) index.adoc

.PHONY: serve
serve: third_party
	@go run ./scripts/serve

.PHONY: third_party
third_party:
	make -C third_party
