source ~/.zprofile
nohup  ./node --config ./node_local.json run >local.log 2>&1 &
tail -f local.log
