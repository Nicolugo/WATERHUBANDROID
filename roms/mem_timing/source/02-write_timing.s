; Tests timing of accesses made by
; memory write instructions

.include "shell.inc"
.include "tima_64.s"

instructions:
     ; last value is time of write
     .byte $36,$FF,$0