PREFIX=/
OUTFILE=./ghost

all:
	go build -o ${OUTFILE}

clean:
	@rm -vf ${OUTFILE}
	@echo Cleaned!

install:
	cp ${OUTFILE} ${PREFIX}/usr/bin/

uninstall:
	rm ${PREFIX}/usr/bin/${OUTFILE}

.PHONY: all clean install
