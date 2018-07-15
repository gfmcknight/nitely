$MYDIR="nitely-test-repo"
$START=$PWD

New-Item -Path $env:HOME -Name $MYDIR -ItemType Directory
Set-Location ~/$MYDIR
Write-Output -InputObject "TEST A" | Out-File ./file-a.txt -Encoding ASCII

git init
git add file-a.txt
git commit -m "Change 1"

git checkout -b other
Write-Output -InputObject "TEST B" | Out-File ./file-b.txt -Encoding ASCII
Remove-Item file-a.txt

git add .
git commit -m "Change 2"
git checkout master

git checkout -b other2
Write-Output -InputObject "TEST C" | Out-File ./file-c.txt -Encoding ASCII
Write-Output -InputObject "ADDITION" | Out-File ./file-a.txt -Append -Encoding ASCII
git add .
git commit -m "Change 3"
git checkout master

Set-Location $START