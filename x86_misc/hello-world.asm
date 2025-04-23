;; DEFINE DATA SECTION OF THE CODE!!!!!
section .data
    ;;
    msg db      "Hello, World!", 0x0A

;; Definition of the text section!!
section .text
    ;;
    global _start

;;Entry point
_start:
    ;; syscall number (1 is sys_write)
    mov     rax, 1
    ;; sets first sys_write arg for syscall to 1 (stdout)
    mov     rdi, 1
    ;; sets second sys_write arg to the reference of the msg variable
    mov     rsi, msg
    ;; sets third sys_write arg to the length of the msg variable (13 chars or more specifically 13 bytes)
    mov     rdx, 14

    ;; call sys_write
    syscall
    ;; specify the number of the syscall (60 is sys_exit)
    mov     rax, 60
    ;; set first syscall arg to 0 (status successful code)
    mov     rdi, 0
    ;; call sys_exit
    syscall