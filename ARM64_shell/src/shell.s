.section .rodata
prompt: .asciz "$: "
env_path: .asciz "PATH=/bin"
envp: .quad env_path, 0

.section .bss
buffer: .skip 256
argv:   .skip 80
cmd:	.skip 128

.section .text
.global _start

// load pointers into memory for use across the program
_start:
    ldr x19, =prompt
    ldr x20, =envp
    ldr x21, =buffer
    ldr x22, =argv
    ldr x23, =cmd

write_prompt:
    mov x0, #1  // stdout
    mov x1, x19 // prompt
    mov x2, #3
    mov x8, #64 // write syscall
    svc #0
    b clear_input
read_input:
    mov x0, #0  // stdin
    mov x1, x21 // read into buffer
    mov x2, #256
    mov x8, #63 // read syscall
    svc #0
run_command:
    b strip_input
_parse:
    bl parse_command
    ldrb w0, [x21]
    cmp w0, #0 // skip execution on empty command
    beq write_prompt
_fork_process:
    mov x0, #0x11
    mov x1, #0
    mov x2, #0
    mov x3, #0
    mov x4, #0
    mov x8, #220 // clone syscall
    svc #0
    cmp x0, #0
    beq child_exec
_wait_exec:
    mov x0, #-1 // any child
    mov x1, #0
    mov x2, #0
    mov x3, #0
    mov x8, #260 // wait4 syscall
    svc #0
    b write_prompt // loop for new command


clear_input:
    mov x0, #0
_clear_buffer:
    strb wzr, [x21, x0]
    add x0, x0, 1
    cmp x0, #256
    bne _clear_buffer
_clear_argv:
    stp xzr, xzr, [x22]
    stp xzr, xzr, [x22, #0x10]
    stp xzr, xzr, [x22, #0x20]
    stp xzr, xzr, [x22, #0x30]
    b read_input


strip_input:
    mov x2, #0
_strip_loop:
    ldrb w3, [x21, x2]
    cmp w3, #'\n'
    beq _remove_newline
    cmp w3, #0 // null char = end of string
    beq _parse
    add x2, x2, #1
    b _strip_loop
_remove_newline:
    strb wzr, [x21, x2]
    b _parse


parse_command:
    stp fp, lr, [sp, #-0x10]!
    mov x0, #0 // index for loop
    mov x3, #0 // arg count
    mov x5, #8
    str x21, [x22]
_find_args:
    ldrb w2, [x21, x0]
    cmp w2, #0 // null means no args
    beq _end_parse
    cmp w2, #' ' // space means args
    bne _continue
    strb wzr, [x21, x0] // change space to null
    // store address for args
    add x0, x0, #1
    add x3, x3, #1
    add x4, x21, x0
    mul x6, x5, x3
    str x4, [x22, x6]
    b _find_args
_continue:
    add x0, x0, #1
    b _find_args
_end_parse:
    ldp fp, lr, [sp], #0x10
    ret


child_exec:
    mov x0, x21 // buffer pointer
    mov x1, x22 // argv pointer
    mov x2, x20 // envp pointer
    mov x8, #221 // execve syscall
    svc #0
    b exit // exit if exec fails


exit:
    mov x0, #0
    mov x8, 93
    svc #0
