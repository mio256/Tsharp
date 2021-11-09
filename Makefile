ifeq ($(OS), Windows_NT)
	exec = tsh.exe
else
	exec = tsh.out
endif

sources = $(wildcard src/*.c)
objects = $(sources:.c=.o)
flags = -g

$(exec): $(objects)
	gcc $(objects) $(flags) -o $(exec)

%.o: %.c include/%.h
	gcc -c $(flags) $< -o $@

install:
	make
ifeq ($(OS), Windows_NT)
	echo "Not Available in Windows"
else
	cp ./tsh.out /usr/local/bin/tsh
endif

clean:
ifeq ($(OS), Windows_NT)
	-rm *.exe
else
	-rm *.out
endif
	-rm *.o
	-rm src/*.o
