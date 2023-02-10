; Printing of numeric values

; Prints value of indicated register/pair
; as 2/4 hex digits, followed by a space.
; Updates checksum with printed values.
; Preserved: AF, BC, DE, HL

print_regs:
     call print_af
     call print_bc
     call print_de
     call print_hl
     call print_newline
     ret

print_a:
     push af
print_a_:
     call print_hex
     ld   a,' '
     call print_char_nocrc
     pop  af
     ret

print_af:
     push af
     call print_hex
     pop  af
print_f:
     push bc
     push af
     pop  bc
     call print_c
     pop  bc
     ret

print_b:
     push af
     ld   a,b
     jr   print_a_

print_c:
     push af
     ld   a,c
     jr   print_a_

print_d:
     push af
     ld   a,d
     jr   print_a_

print_e:
     push af
     ld   a,e
     jr   print_a_

print_h:
     push af
     ld   a,h
     jr   print_a_

print_l:
     push af
     ld   a,l
     jr   print_a_

print_bc:
     push af
     push bc
print_bc_:
     ld   a,b
     call print_hex
     ld   a,c
     pop  bc
     jr   print_a_
     
print_de:
     push af
     push bc
     ld   b,d
     ld   c,e
     jr   print_bc_
     
print_hl:
     push af
     push bc
     ld   b,h
     ld   c,l
     jr   print_bc_
     

; Prints A as two hex chars and updates checksum
; Preserved: BC, DE, HL
print_hex:
     call update_crc
print_hex_nocrc:
     push af
     swap a
     call +
     pop  af
     
+    and  $0F
     cp   10
     jr   c,+
     add  7
+    add  '0'
     jp   print_char_nocrc


; Prints char_nz if Z flag is clear,
; char_z if Z flag is set.
; Preserved: AF, BC, DE, HL
.macro print_nz ARGS char_nz, char_z
     push af
     ld   a,char_nz
     jr   nz,print_nz\@
     ld   a,char_z
print_nz\@:
     call print_char
     pop  af
.e