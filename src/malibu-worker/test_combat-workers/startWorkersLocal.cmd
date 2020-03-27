@echo off
set /p UserInput=Enter count: 
set /a Test=UserInput

echo %UserInput%

FOR /L %%A IN (1,1,%UserInput%) DO (
  del /F /S /Q %%A
  timeout 1
  mkdir %%A
  copy ..\combat-worker.exe %%A\
  cd %%A && start combat-worker.exe http://localhost:9090 && cd ..
  ECHO %%A
)