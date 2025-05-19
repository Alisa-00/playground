.section .rodata
prefix: .asciz "$: "
env_path: .asciz "PATH=/bin"
envp: .quad env_path, 0

.section .data
buffer: .space 256
argv: .quad 0, 0, 0

.section .text
.global _start


_start:
    ldr x19, =prefix // x19 = prefix pointer
    ldr x20, =envp   // x20 = envp pointer
    ldr x21, =buffer // x21 = buffer pointer
    ldr x22, =argv   // x22 = argv pointer
main:
    bl write_prefix
    bl read_input
    bl run_command
    b main

write_prefix:
    stp fp, lr, [sp, #-0x10]!
    mov x0, #1
    mov x1, x19 
    mov x2, #3
    mov x8, #64 // write syscall
    svc #0
    ldp fp, lr, [sp], #0x10
    ret

read_input:
    stp fp, lr, [sp, #-0x10]!
    bl clear_buffer
    bl clear_argv
    mov x0, #0
    mov x1, x21 
    mov x2, #256 // buffer length
    mov x8, #63 // read syscall
    svc #0
    bl strip_input
    ldp fp, lr, [sp], #0x10
    ret

clear_buffer:
    stp fp, lr, [sp, #-0x10]!
    mov x0, #0
_clear_loop:
    strb wzr, [x21, x0]
    add x0, x0, 1
    cmp x0, #256
    bne _clear_loop
    ldp fp, lr, [sp], #0x10
    ret

clear_argv:
    stp fp, lr, [sp, #-0x10]!
    stp xzr, xzr, [x22]
    str xzr, [x22, #16]
    ldp fp, lr, [sp], #0x10
    ret

strip_input:
    stp fp, lr, [sp, #-0x10]!
    mov x2, #0
_strip_loop:
    ldrb w3, [x21, x2]
    cmp w3, #'\n'
    beq _remove_newline
    cmp w3, #0 // null char
    beq _end_loop
    add x2, x2, #1
    b _strip_loop
_remove_newline:
    strb wzr, [x21, x2]
_end_loop:
    ldp fp, lr, [sp], #0x10
    ret    

run_command:
    stp fp, lr, [sp, #-0x10]!
    bl parse_command
    ldrb w0, [x21]
    cmp w0, #0 // empty command
    beq _end_execution
    bl fork_process // fork to create child
    cmp x0, #0
    beq child_exec
    bl parent_wait
_end_execution:
    ldp fp, lr, [sp], #0x10
    ret

parse_command:
    stp fp, lr, [sp, #-0x10]!
    mov x0, #0
    str x21, [x22]
_find_args:
    ldrb w2, [x21, x0]
    add x0, x0, #1
    cmp w2, #0 // null means no args
    beq _end_parse
    cmp w2, #' ' // space means args
    bne _find_args
    sub x0, x0, #1
    strb wzr, [x21, x0] // change space to null
    // store address for args
    add x0, x0, #1
    add x4, x21, x0
    str x4, [x22, #8]
_end_parse:
    ldp fp, lr, [sp], #0x10
    ret

fork_process:
    stp fp, lr, [sp, #-0x10]!
    mov x0, #0x11
    mov x1, #0
    mov x2, #0
    mov x3, #0
    mov x4, #0
    mov x8, #220 // clone syscall
    svc #0
    ldp fp, lr, [sp], #0x10
    ret

child_exec:
    stp fp, lr, [sp, #-0x10]!
    mov x0, x21 // buffer pointer
    mov x1, x22 // argv pointer
    mov x2, x20 // envp pointer
    mov x8, #221 // execve syscall
    svc #0
    b exit // exit if exec fails

parent_wait:
    stp fp, lr, [sp, #-0x10]!
    mov x0, #-1
    mov x1, #0
    mov x2, #0
    mov x3, #0
    mov x8, #260 // wait4 syscall
    svc #0
    ldp fp, lr, [sp], #0x10
    ret

exit:
    mov x0, #0
    mov x8, 93
    svc #0

