
; Scrolling text console

; Console is 20x18 characters. Buffers lines, so
; output doesn't appear until a newline or flush.
; If scrolling isn't supported (i.e. SCY is treated
; as if always zero), the first 18 lines will
; still print properly). Also works properly if
; LY isn't supported (always reads back as the same
; value).

.define console_width 20

.define console_buf      bss+0
.define console_pos      bss+console_width
.define console_mode     bss+console_width+1
.define console_scroll   bss+console_width+2
.redefine bss            bss+console_width+3


; Waits for start of LCD blanking period
; Preserved: BC, DE, HL
console_wait_vbl:
     push bc
     
     ; Wait for start of vblank, with
     ; timeout in case LY doesn't work
     ; or LCD is disabled.
     ld   bc,-1250
-    inc  bc
     ld   a,b
     or   c
     jr   z,@timeout
     lda  LY
     cp   144
     jr   nz,-
@timeout:
     
     pop  bc
     ret


; Initializes text console
console_init:
     call console_hide
     
     ; CGB-specific inits
     ld   a,(gb_id)
     and  gb_id_cgb
     call nz,@init_cgb
     
     ; Clear nametable
     ld   a,' '
     call @fill_nametable
     
     ; Load tiles
     ld   hl,TILES+$200
     ld   c,0
     call @load_tiles
     ld   hl,TILES+$A00
     ld   c,$FF
     call @load_tiles
     
     ; Init state
     ld   a,console_width
     ld   (console_pos),a
     ld   a,0
     ld   (console_mode),a
     ld   a,-8
     ld   (console_scroll),a
     call console_scroll_up_
     jr   console_show

@fill_nametable:
     ld   hl,BGMAP0
     ld   b,4
-    ld   (hl),a
     inc  l
     jr   nz,-
     inc  h
     dec  b
     jr   nz,-