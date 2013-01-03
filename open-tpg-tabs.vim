
let baseDir = "~/Workspace/go-tp/src/ble/tpg/"
let subDirs = ["ephemeral", "handler", "model", "persistence", "room", 
              \"switchboard", "test"]
for subDir in subDirs
  tabnew
  TName subDir
  let theDir = baseDir.subDir
  execute "lcd ".theDir
  edit .
endfor
