@echo off
set loopcount=100
:loop

main -algo=4 -msqrt -debug

timeout /t 5 /NOBREAK

set /a loopcount=loopcount-1
if %loopcount%==0 goto exitloop
goto loop
:exitloop
