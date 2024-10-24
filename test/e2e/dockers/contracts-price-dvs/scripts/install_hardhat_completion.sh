#!/usr/bin/expect
set time 10
spawn hardhat-completion install
expect "Which Shell do you use ?"
send -- "bash\r"
expect "We will install completion to ~/.bashrc, is it ok ? (y/N)"
send -- "y\r"
expect eof