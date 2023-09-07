.PHONY: all
all: bin/pikchr

bin/pikchr: pikchr.c
	mkdir -p $(dir $@)
	gcc -DPIKCHR_SHELL -o $@ $< -lm
