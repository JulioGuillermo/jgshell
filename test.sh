#!/bin/bash

delay=0.05
# Los caracteres Unicode funcionan mejor si los manejamos como un array
spinstr=("" "" "" "" "" "" "" "" "" "" "" "" "" "" "" "" "" "" "" "" "" "" "" "" "" "" "" "")

# Ocultar el cursor para que se vea más profesional
tput civis
# Asegurar que el cursor vuelva al salir (Ctrl+C)
trap "tput cnorm; exit" SIGINT SIGTERM

echo -n "Loading "

for (( c=1; c<=10; c++ )); do
    for i in "${spinstr[@]}"; do
        # \r vuelve al inicio de la línea, luego Loading y el ícono
        # Usamos %s para strings Unicode
        printf "\rLoading [%s]  " "$i"
        sleep $delay
    done
done

echo -e "\nDone"
