@echo off
set args=%*
set arguments=%args:\=\\%

c:\vcap\bosh\agent\boshtar.exe --force-local %arguments%