@echo off
@REM ===================
@REM ==== ARGUMENTS ====
@REM ===================
setlocal enabledelayedexpansion
set indexes=0 1
set flags[0]=cp
set vars[0]=CLIENT_PORT
set flags[1]=au
set vars[1]=API_URL
for %%a in (%*) do (
    if %%a==-h (
        echo -cp : Client port [default 3000]
        echo -au : Api full url [default http://localhost:8080]
        echo -h  : Help
        goto :eof
    )
)
for %%i in (%indexes%) do (
    set found=0
    set flag=!flags[%%i]!
    set var=!vars[%%i]!
    for %%a in (%*) do (
        if "!found!"=="0" (
            if %%a==-!flag! (
                set found=1
            )
        ) else (
            if "!found!"=="1" (
                set !var!=%%a
                set found=2
            )
        )
    )
)
@REM ===================
@REM ===================
@REM ===================

set GOOS=js
set GOARCH=wasm
if not "%API_URL%"=="" (
    set args_build= -ldflags "-X main.ApiUrl=%API_URL%"
)
echo Building wasm/main.go
go build%args_build% -o public\main.wasm wasm\main.go

set GOOS=
set GOARCH=
echo Building client.go
go build -o build\ client.go
if not "%CLIENT_PORT%"=="" (
    set args_client= -p %CLIENT_PORT%
)
build\client%args_client%
