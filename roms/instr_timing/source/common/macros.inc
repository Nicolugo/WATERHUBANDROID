; General macros

; Reads A from addr, from $FF00 to $FFFF
; Preserved: F, BC, DE, HL
; Time: 3 cycles
.macro lda ARGS addr
     ldh  a,(addr - $FF00)
.endm

; Writes A to addr, from $FF00 to $FFFF
; Preserved: AF, BC, DE, HL
; Time: 3 cycles
.macro sta ARGS addr
     ldh  (addr - $FF00),a
.endm

; Writes immediate data to addr, from $FF00 to $FFFF
; Preserved: F, BC, DE, HL
; Time: 5 cycles
.macro wreg ARGS addr, data
     ld   a,data
     sta  addr
.endm

; Calls routine multiple times, with A having the
; value 'start' the first time, 'start+step' the
; second time, up to 'end' for the last time.
; Preserved: BC, DE, HL
.macro for_loop ; routine,start,end,step
     ld   a,\2
     
for_loop\@:
     push af
     call \