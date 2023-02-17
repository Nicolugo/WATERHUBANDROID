
; Delays in cycles, milliseconds, etc.

; All routines are re-entrant (no global data). Routines never
; touch BC, DE, or HL registers. These ASSUME CPU is at normal
; speed. If running at double speed, msec/usec delays are half advertised.

; Delays n cycles, from 0 to 16777215
; Preserved: AF, BC, DE, HL
.macro delay ARGS n
     .if n < 0
          .printt "Delay must be >= 0"
          .fail
     .endif
     .if n > 16777215
          .printt "Delay must be < 16777216"
          .fail
     .endif
     delay_ n&$FFFF, n>>16
.endm

; Delays n clocks, from 0 to 16777216*4. Must be multiple of 4.
; Preserved: AF, BC, DE, HL
.macro delay_clocks ARGS n
     .if n # 4 != 0
          .printt "Delay must be a multiple of 4"
          .fail
     .endif
     delay_ (n/4)&$FFFF,(n/4)>>16
.endm

; Delays n microseconds (1/1000000 second)
; n can range from 0 to 4000 usec.
; Preserved: AF, BC, DE, HL
.macro delay_usec ARGS n
     .if n < 0
          .printt "Delay must be >= 0"
          .fail
     .endif
     .if n > 4000
          .printt "Delay must be <= 4000 usec"
          .fail
     .endif
     delay_ ((n * 1048576 + 500000) / 1000000)&$FFFF,((n * 1048576 + 500000) / 1000000)>>16
.endm

; Delays n milliseconds (1/1000 second)
; n can range from 0 to 10000 msec.
; Preserved: AF, BC, DE, HL
.macro delay_msec ARGS n
     .if n < 0
          .printt "Delay must be >= 0"
          .fail
     .endif
     .if n > 10000
          .printt "Delay must be <= 10000 msec"
          .fail
     .endif
     delay_ ((n * 1048576 + 500) / 1000)&$FFFF, ((n * 1048576 + 500) / 1000)>>16
.endm

     ; All the low/high quantities are to deal wla-dx's asinine
     ; restriction full expressions must evaluate to a 16-bit
     ; value. If the author ever rectifies this, all "high"
     ; arguments can be treated as zero and removed. Better yet,
     ; I'll just find an assembler that didn't crawl out of
     ; the sewer (this is one of too many bugs I've wasted
     ; hours working around).

     .define max_short_delay 28
     
     .macro delay_long_ ARGS n, high
          ; 0+ to avoid assembler treating as memory read
          ld   a,0+(((high<<16)+n) - 11) >> 16
          call delay_65536a_9_cycles_
          delay_nosave_ (((high<<16)+n) - 11)&$FFFF, 0
     .endm
     
     ; Doesn't save AF, allowing minimization of AF save/restore
     .macro delay_nosave_ ARGS n, high
          ; 65536+11     = maximum delay using delay_256a_9_cycles_
          ; 255+22  = maximum delay using delay_a_20_cycles
          ; 22      = minimum delay using delay_a_20_cycles
          .if high > 1
               delay_long_ n, high
          .else
               .if high*n > 11
                    delay_long_ n, high
               .else