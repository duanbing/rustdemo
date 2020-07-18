[org 0x7c00]

mov bx, HELLO_MSG;
call print_string
call print_nl

mov bx , GOODBYE_MSG
call print_string

jmp $

%include "print_string.asm"

HELLO_MSG:
	db 'hello, world!', 0

GOODBYE_MSG:
	db 'Good bye!', 0

times 510 - ($ - $$) db 0
dw 0xaa55
