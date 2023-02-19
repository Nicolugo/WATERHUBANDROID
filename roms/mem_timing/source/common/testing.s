; Diagnostic and testing utilities

.define result      bss+0
.define test_name   bss+1
.redefine bss       bss+3


; Sets test code and optional error text
; Preserved: AF, BC, DE, HL
.macro set_test ; code[,text[,text2]]
     push hl
     call set_test_
     jr   @set_test\@
     .byte \1
     .if NARGS > 1
     