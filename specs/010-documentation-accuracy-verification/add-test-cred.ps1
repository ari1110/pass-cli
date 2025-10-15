# Add test credential with proper input sequence
$masterPwd = "TestMaster123!"
$username = "test@example.com"
$password = "TestPass123!"

# Pipe inputs: master password
echo $masterPwd | & "R:\Test-Projects\pass-cli\pass-cli.exe" add testservice -u $username -p $password
