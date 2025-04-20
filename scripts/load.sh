#!/bin/bash

IFS=','

while read -r line; do
    read -ra stocks <<< "$line"
    
    for stock in "${stocks[@]}"; do
        echo "Loading Stock: $stock"
        ./quantgo-cli sync_stock_data -c $stock -t xxx
    done
    echo "------"
done < symbols
