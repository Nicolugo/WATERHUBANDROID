
; Tests timing of accesses made by
; memory read instructions

.include "shell.inc"
.include "tima_64.s"

instructions:
     ; last value is time of read
     .byte $B6,$00,$00,2 ; OR   (HL)
     .byte $BE,$00,$00,2 ; CP   (HL)
     .byte $86,$00,$00,2 ; ADD  (HL)
     .byte $8E,$00,$00,2 ; ADC  (HL)
     .byte $96,$00,$00,2 ; SUB  (HL)
     .byte $9E,$00,$00,2 ; SBC  (HL)
     .byte $A6,$00,$00,2 ; AND  (HL)
     .byte $AE,$00,$00,2 ; XOR  (HL)
     .byte $46,$00,$00,2 ; LD   B,(HL)
     .byte $4E,$00,$00,2 ; LD   C,(HL)
     .byte $56,$00,$00,2 ; LD   D,(HL)
     .byte $5E,$00,$00,2 ; LD   E,(HL)
     .byte $66,$00,$00,2 ; LD   H,(HL)
     .byte $6E,$00,$00,2 ; LD   L,(HL)
     .byte $7E,$00,$00,2 ; LD   A,(HL)
     .byte $F2,$00,$00,2 ; LDH  A,(C)
     .byte $0A,$00,$00,2 ; LD   A,(BC)
     .byte $1A,$00,$00,2 ; LD   A,(DE)
     .byte $2A,$00,$00,2 ; LD   A,(HL+)
     .byte $3A,$00,$00,2 ; LD   A,(HL-)
     .byte $F0,<tima_64,$00,3 ; LDH  A,($00)
     .byte $FA,<tima_64,>tima_64,4 ; LD   A,($0000)
     
     .byte $CB,$46,$00,3 ; BIT  0,(HL)
     .byte $CB,$4E,$00,3 ; BIT  1,(HL)
     .byte $CB,$56,$00,3 ; BIT  2,(HL)
     .byte $CB,$5E,$00,3 ; BIT  3,(HL)
     .byte $CB,$66,$00,3 ; BIT  4,(HL)
     .byte $CB,$6E,$00,3 ; BIT  5,(HL)
     .byte $CB,$76,$00,3 ; BIT  6,(HL)
     .byte $CB,$7E,$00,3 ; BIT  7,(HL)
instructions_end:

main:
     call init_tima_64
     set_test 0
     
     ; Test instructions
     ld   hl,instructions