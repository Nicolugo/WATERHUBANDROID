; CRC-32 checksum calculation

.define checksum    dp+0 ; little-endian, complemented
.redefine dp        dp+4


; Initializes checksum module. Might initialize tables
; in the future.
init_crc:
     jr   reset_crc


; Clears CRC
; Preserved: BC, DE, HL
reset_crc:
     ld   a,$FF
     sta  checksu