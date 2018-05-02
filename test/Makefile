default:
	@echo "No default target."

clean:
	rm -f *.covout *-cover.bats

cover: clean
	./cover.bash

test: clean
	bats *.bats
