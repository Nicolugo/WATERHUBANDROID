; Timer that's incremented every 64 cycles

; Initializes timer for use by sync_tima_64
init_tima_64:
     wreg TMA,0
     wreg TAC,$07
     ret

; Synchronizes to timer
sync_tima_64:
     push af
     push hl
     
     ld   a,0
     ld   hl,TIMA
     ld   (hl),a
-    or   (hl)
     jr   z,-
     
-    delay 65-12
     xor  a
     ld   (hl),a
     or   (hl)
     dela