build:
	nasm -f bin ./boot_sect_simple.asm -o boot_sect_simple.bin

print_hello:
	nasm -f bin $@.asm -o $@.bin

run:
	qemu-system-x86_64 boot_sect_simple.bin --nographic 
