; Infonite loop (e9 fd ff)
loop:
	jmp loop

; Fill iwth 510 zeros minus the size of the previous code 
times 510 - ($-$$) db 0
; Magic number
dw 0xaa55
