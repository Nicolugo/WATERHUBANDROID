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
     .byte $77,$00,$00,2 ; LD   (HL),A
     .byte $02,$00,$00,2 ; LD   (BC),A
     .byte $12,$00,$00,2 ; LD   (DE),A
     .byte $22,$00,$00,2 ; LD   (HL+),A
     .byte $32,$00,$00,2 ; LD   (HL-),A
     .byte $E2,$00,$00,2 ; LDH  (C),A
     .byte $E0,<tima_64,$00,3 ; LDH  (n),A
     .byte $EA,<tima_64,>tima_64,4 ; LD   (nn),A
instructions_end:

main:
     call init_tima_64
     set_test 0
     
     ; Test instructions
     ld   hl,instructions
-    call @test_instr
     cp   (hl)
     call nz,@print_failed
     inc  hl
     ld   a,l
     cp