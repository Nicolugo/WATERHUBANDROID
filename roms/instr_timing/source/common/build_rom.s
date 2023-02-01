; Build as GB ROM

.memoryMap
     defaultSlot 0
     slot 0 $0000 size $4000
     slot 1 $C000 size $4000
.endMe

.romBankSize   $4000 ; generates $8000 byte ROM
.romBanks      2

.cartridgeType 1 ; MBC1
.computeChecksum
.computeComplementCheck


;;;; GB ROM header

; GB header read by bootrom
.org $100
     nop
     jp   reset

; Nintendo logo required for proper boot
.byte $CE,$ED,$66,$66,$CC,$0D,$00,$0B
.byte $03,$73,$00,$83,$00,$0C,$00,$0D
.byte $00,$08,$11,$1F,$88,$89,$00,$0E
.byte $DC,$CC,$6E,$E6,$DD,$DD,$D9,$99
.byte $BB,$BB,$67,$63,$6E,$0E,$EC,$CC
.byte $DD,$DC,$99,$9F,$BB,$B