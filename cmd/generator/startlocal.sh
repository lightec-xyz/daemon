source ~/.zprofile
nohup ./generator --config ./local.json run >local.log 2>&1 &
tail -f local.log
