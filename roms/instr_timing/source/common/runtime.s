
; Common routines and runtime

; Must be defined by target-specific runtime:
;
; init_runtime:     ; target-specific inits
; std_print:        ; default routine to print char A
; post_exit:        ; called at end of std_exit
; report_byte:      ; report A to user

.define RUNTIME_INCLUDED 1

.ifndef bss
     ; address of next normal variable
     .define bss    $D800
.endif

.ifndef dp
     ; address of next direct-page ($FFxx) variable
     .define dp     $FF80
.endif

; DMG/CGB hardware identifier
.define gb_id_cgb        $10  ; mask for testing CGB bit
.define gb_id_devcart    $04 ; mask for testing "on devcart" bit

.define gb_id  bss
.redefine bss  bss+1

; Stack is normally here
.define std_stack $DFFF

; Copies $1000 bytes from HL to $C000, then jumps to it.
; A is preserved for jumped-to code.
copy_to_wram_then_run:
     ld   b,a
     
     ld   de,$C000
     ld   c,$10
-    ldi  a,(hl)
     ld   (de),a
     inc  e
     jr   nz,-
     inc  d
     dec  c
     jr   nz,-
     
     ld   a,b
     jp   $C000

.ifndef CUSTOM_RESET
     reset:
          ; Run code from $C000, as is done on devcart. This
          ; ensures minimal difference in how it behaves.
          ld   hl,$4000
          jp   copy_to_wram_then_run
     
     .bank 1 slot 1
     .org $0 ; otherwise wla pads with lots of zeroes
          jp   std_reset
.endif

; Common routines
.include "gb.inc"
.include "macros.inc"
.include "delay.s"
.include "crc.s"
.include "printing.s"
.include "numbers.s"
.include "testing.s"
     
; Sets up hardware and runs main
std_reset:

     ; Init hardware
     di
     ld   sp,std_stack
     
     ; Save DMG/CGB id
     ld   (gb_id),a
     
     ; Init hardware
     .ifndef BUILD_GBS
          wreg TAC,$00
          wreg IF,$00
          wreg IE,$00
     .endif
     