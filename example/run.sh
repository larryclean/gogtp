#!/bin/bash

cd /Users/larry/Documents/01-Project/os/KataGo/cpp;
./katago gtp  -config configs/gtp_example.cfg -model models/g170-b6c96-s175395328-d26788732.bin.gz
#./katago gtp  -config configs/gtp_example.cfg -model models/g170-b6c96-s175395328-d26788732.txt.gz
#ssh -p8622 root@192.168.3.244 '/root/ai_service/run.sh'

#ssh root@47.98.142.37 '/mnt/ai_server/run.sh'