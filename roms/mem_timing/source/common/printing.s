; Main printing routine that checksums and
; prints to output device

; Character that does equivalent of print_newline
.define newline 10

; Prints char without updating checksum
; Preserved: BC, DE, HL
.define print_char_nocrc bss
.redefine bss            bss+3


; Initializes printing. HL = print routine
init_printing:
     ld   a,l
     ld   (print_char_nocrc+1),a
     ld   a,h
     ld   (print