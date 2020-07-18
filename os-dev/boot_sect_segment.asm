mov ah, 0x0e

mov al, [the_secret]
int 0x10

call print_nl

mov bx, 0x7c0
mov ds, bx;
mov al, [the_secret]
int 0x10
call print_nl

mov al, [es:the_secret]
int 0x10 ; doesn't look right... isn't 'es' currently 0x000?
call print_nl

mov bx, 0x7c0
mov es, bx
mov al, [es:the_secret] ;  else segment
int 0x10
call print_nl


jmp $

%include "print_string.asm"

the_secret:
    db "X"

times 510 - ($-$$) db 0
dw 0xaa55



