ASY= $(wildcard *.asy)
ASYPNG= $(ASY:.asy=.png)
ASYPDF= $(ASY:.asy=.pdf)
ASYSVG= $(ASY:.asy=.svg)

default: svg
all: png pdf svg

png: $(ASYPNG)
pdf: $(ASYPDF)
svg: $(ASYSVG)

%.png: %.asy
	asy -o $@ -f png $^

%.pdf: %.asy
	asy -o $@ -f pdf $^

%.svg: %.asy
	asy -o $@ -f svg $^

distclean:
	rm -f $(ASYPDF) $(ASYPNG) $(ASYSVG)

mrproper: distclean
