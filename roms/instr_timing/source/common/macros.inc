; General macros

; Reads A from addr, from $FF00 to $FFFF
; Preserved: F, BC, DE, HL
; Time: 3 cycles
.macro lda ARGS addr
     ldh  a,(addr - $FF00