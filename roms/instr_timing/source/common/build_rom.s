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
   