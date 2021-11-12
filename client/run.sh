# ===================
# ==== ARGUMENTS ====
# ===================
flags=(cp au)
vars=(CLIENT_PORT API_URL)
for arg in "$@"; do
    if [[ "$arg" == "-h" ]]; then
        echo "-cp : Client port [default 3000]"
        echo "-au : Api full url [default http://localhost:8080]"
        echo "-h  : Help"
        exit
    fi
done
i=0
for flag in "${flags[@]}"; do
    var=${vars[$i]}
    found=0
    for arg in "$@"; do
        if [[ $found -eq 0 && "$arg" == "-$flag" ]]; then
            found=1
        elif [[ $found -eq 1 ]]; then
            declare $var=$arg
            break
        fi
    done
    i=${i+1}
done
# ===================
# ===================
# ===================

command_build="GOOS=js GOARCH=wasm go build"
if [[ "$API_URL" != "" ]]; then
    command_build+=" -ldflags \"-X main.ApiUrl=$API_URL\""
fi
command_build+=" -o public/main.wasm wasm/main.go"
echo Building wasm/main.go
eval $command_build

echo Building client.go
GOOS= GOARCH= go build -o build/ client.go
if [[ "$CLIENT_PORT" != "" ]]; then
    args_client=" -p $CLIENT_PORT"
fi
./build/client$args_client
