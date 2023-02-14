; Tests timing of accesses made by
; memory write instructions

.include "shell.inc"
.include "tima_64.s"

instructions:
     ; last value is time of write
     .byte $36,$FF,$00,3 ; LD   (HL),n
     .byte $70,$00,$00,2 ; LD   (HL),B
     .byte $71,$00,$00,2 ; LD   (HL),C
     .byte $72,$00,$00,2 ; LD   (HL),D
     .byte $73,$00,$00,2 ; LD   (HL),E
     .byte $74,$00,$00,2 ; LD   (HL),H
     .byte $75,$00,$00,2 ; LD   (HL),L
     