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
          .byte \2
     .endif
     .if NARGS > 2
          .byte \3
     .endif
     .byte 0
@set_test\@:
     pop  hl
.endm

set_test_:
     pop  hl
     push hl
     push af
     inc  hl
     inc  hl
     ldi  a,(hl)
     ld   (result),a
     ld   a,l
     ld   (test_name),a
     ld   a,h
     ld   (test_name+1),a
     pop  af
     ret


; Initializes testing module
init_testing:
     set_test $FF
     call init_crc
     ret


; Reports "Passed", then exits with code 0
tests_passed:
     call print_newline
     print_str "Passed"
     ld   a,0
     jp   exit


; Reports "Done" if set_test has never been used,
; "Passed" if set_test 0 was last used, or
; failure if set_test n was last used.
tests_done:
     ld   a,(result)
     inc  a
     jr   z,+
     dec  a
     jr   z,tests_passed
     jr   test_failed
+    print_str 