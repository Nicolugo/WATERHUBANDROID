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
     jr   prin