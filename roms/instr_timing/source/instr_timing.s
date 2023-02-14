
; Tests number of cycles taken by instructions
; except STOP, HALT, and illegals.

.include "shell.inc"
.include "timer.s"

.define saved_sp    bss+0
.define instr       bss+2 ; 3-byte instr + JP instr_end
.define instr_addr  bss+8 ; JP instr_end
.redefine bss       bss+11

main:
     call init_timer
     call test_timer
     set_test 0
     call test_main_ops
     call test_cb_ops
     jp   tests_done


; Ensures timer works
test_timer:
     call start_timer
     call stop_timer
     or   a
     ret  z
     set_test 2,"Timer doesn't  work properly"
     jp   test_failed


; Tests main opcodes
test_main_ops:
     ld   l,0
-    ld   h,>op_times
     ld   a,(hl)
     cp   0
     call nz,@test_op
     inc  l
     jr   nz,-
     ret

@test_op:
     ; Can't test the 8 RST instructions on devcart
     ld   a,l
     cpl
     and  $C7
     jr   nz,+
     ld   a,(gb_id)
     and  gb_id_devcart
     ret  nz
+
     ; Test with flags set so that branches are
     ; not taken
     ld   a,l  ; e = (l & 0x08 ? 0 : 0xFF)
     and  $08
     add  $F8
     ld   e,a
     call @copy_and_exec
     ld   d,0
     cp   (hl)
     jr   z,+
     ld   d,a
     call print_failed_opcode
+

     ; Time with branches not taken
     ld   a,e
     cpl
     ld   e,a
     call @copy_and_exec
     ld   h,>op_times_taken
     cp   (hl)
     ret  z
     
     ; If opcode already failed and timed the
     ; same again, avoid re-reporting.
     cp   d
     ret  z
     
     call print_failed_opcode
     ret

@copy_and_exec:
     push de
     push hl
     
     ld   h,>op_lens
     ld   c,(hl)
     ld   a,l
     ld   hl,instr
     ld   (hl+),a
     dec  c
     jr   z,@one_byte
     ld   a,0
     dec  c
     jr   z,@two_bytes
     ld   a,<instr_addr
     ld   (hl+),a
     ld   a,>instr_addr
@two_bytes:
     ld   (hl+),a
@one_byte:     
     ld   a,e
     call time_instruction
     
     pop  hl
     pop  de
     ret


; Tests CB opcodes
test_cb_ops:
     ld   hl,cb_op_times
-    ld   a,(hl)
     cp   0
     call nz,@test_op_cb
     inc  l
     jr   nz,-
     ret

@test_op_cb:
     ; Test with flags clear
     ld   e,$00
     call @copy_and_exec_cb
     cp   (hl)
     jr   nz,+
     
     ; Test with flags set
     ld   e,$FF
     call @copy_and_exec_cb
     cp   (hl)
     jr   nz,+
     
     ret
+    print_str "CB "
     call print_failed_opcode
     ret

@copy_and_exec_cb:
     push hl
     
     ; Copy instr to exec space
     ld   a,l
     ld   hl,instr+1
     ld   (hl+),a
     ld   a,$CB
     ld   (instr),a
     call time_instruction
     
     pop  hl
     ret


; Reports failed opcode
; L    -> opcode
; A    -> cycles it took
; (HL) -> cycles it should have taken
; Preserved: HL
print_failed_opcode:
     ; Print opcode
     push af
     ld   a,l
     call print_hex
     ld   a,':'
     call print_char
     pop  af
     
     ; Print actual and correct times
     call print_dec
     ld   a,'-'
     call print_char
     ld   a,(hl)
     call print_dec
     ld   a,' '
     call print_char
     
     ; Remember that failure occurred
     set_test 1
     
     ret


; Times instruction.
; HL -> address of byte just after instruction
; A  -> flags when executing instruction
; A  <- number of cycles instruction took
time_instruction:
     ld   c,a
     
     ; Write JP instr_end to HL and instr_addr
     ld   a,$C3     ; JP
     ld   (hl+),a
     ld   (instr_addr),a
     
     ld   a,<instr_end
     ld   (instr_addr+1),a
     ld   (hl+),a
     
     ld   a,>instr_end
     ld   (instr_addr+2),a
     ld   (hl),a
     
     ; Save sp
     ld   (saved_sp),sp
     
     ; Set regs and stack contents
     push bc
     ld   bc,instr_addr
     ld   de,instr_addr
     ld   hl,instr_addr
     call start_timer
     pop  af
     push hl
     
     ; Environment instruction executes in:
     ; 1 byte: OP
     ; 2 byte: OP 00
     ; 3 byte: OP instr_addr
     ; BC,DE,HL = instr_addr
     ; Stack has instr_addr pushed on it.
     ; Stack pointer can be trashed by instr.
     ; instr_addr contains JP instr_end, that
     ; can be trashed. Instructions which trash
     ; this don't execute it.
     
     jp   instr
instr_end:     ; instruction jumps here when done
     di
     