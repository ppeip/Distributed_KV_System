# 定义要执行的服务器进程  
<#
$serverProcesses = @(  
    { go build -o server; ./server -port=8001 },  
    { go build -o server; ./server -port=8002 },  
    { go build -o server; ./server -port=8003 -api=1 }  
)  
  
# 启动服务器进程  
foreach ($process in $serverProcesses) {  
    Start-Process powershell -ArgumentList "-NoProfile -WindowStyle Hidden -Command `"$process`""  
}  
#>

go build -o server.exe

#Start-Sleep -Seconds 2 
<#
./server -port=8001 
./server -port=8002 
./server -port=8003 -api=1 
#>

$server8001 = Start-Process -FilePath ".\server.exe" -ArgumentList "-port=8001" -PassThru  
$server8002 = Start-Process -FilePath ".\server.exe" -ArgumentList "-port=8002" -PassThru  
$server8003 = Start-Process -FilePath ".\server.exe" -ArgumentList "-port=8003 -api=1" -PassThru  
# 等待服务器进程启动  

Start-Sleep -Seconds 2  
  
# 执行 HTTP 请求  

  
curl "http://localhost:9999/api?key=Tom" 
curl "http://localhost:9999/api?key=Tom" 
curl "http://localhost:9999/api?key=Tom" 
curl "http://localhost:9999/api?key=Tom" 
curl "http://localhost:9999/api?key=Tom" 
curl "http://localhost:9999/api?key=Tom" 
curl "http://localhost:9999/api?key=Tom" 

$server8001.WaitForExit()  
$server8002.WaitForExit()  
$server8003.WaitForExit()  
if ($server8001.HasExited) {  
    $server8001.Close()  
    Remove-Item ".\server.exe"  
}  


#Remove-item server -Force
# 清理（如果需要的话）  
# PowerShell 没有类似于 Bash 中 EXIT trap 的功能，但可以在脚本结束时执行清理操作  
# 例如，如果需要删除某个文件，可以在脚本的末尾添加相应的命令  
# Remove-Item server -Force  
  
# 注意：以上脚本假设 go 和 curl 已经在 PowerShell 环境中可用  
# 如果不是，你可能需要安装它们或调整脚本以使用 PowerShell 的原生命令