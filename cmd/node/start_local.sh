source ~/.zprofile
nohup  ./node --config ./local.json run >local.log 2>&1 &
tail -f local.log
