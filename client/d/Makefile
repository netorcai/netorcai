default: lib

library:
	dub build

library-cov:
	dub build -b cov

test:
	dub test

test-cov:
	dub test -b unittest-cov

clean:
	rm -f -- *.lst
	rm -f -- libnetorcai-client.*
	rm -f -- netorcai-client-*
