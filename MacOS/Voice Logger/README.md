# Phantom Audio Capturer

## Description
Phantom Audio Capturer, developed by Narsty and updated by YonKu0, is an advanced macOS payload designed for discreet audio recording. It automatically records audio via the microphone and stores it as an MP3 file in a concealed folder. This script leverages `sox` for audio recording and handles microphone permissions automatically using either a Golang or Python script (featuring `pyautogui` and `cv2`).

## Features
- Creates a hidden folder `.phantom_audio` for recordings.
- Automated Permission Management: Handles microphone permissions seamlessly in the background.
- Long-duration Recording: Capable of continuous audio recording until manually interrupted.
- Simple File Access: Facilitates easy access to the recorded audio files.

## Requirements
- MacOS
- `sox` (Install via Homebrew: `brew install sox`)

## Usage
1. Execute the script on the target MacOS device.
2. The script automatically starts recording audio in the background.
3. To access the recordings, navigate to `~/.phantom_audio`.

## Stopping the Recording
1. Find the PID of the `sox` process.
2. Execute `kill pid#` or `killall -9 sox` to stop the recording process.

## Cleaning Up
- Delete the `.phantom_audio` folder with `rm -r ~/.phantom_audio`.

## Go Script (Currently used) 
- Function: Handles automated microphone permission via a Golang script.
- For re-compilation: Download the Voice Logger directory and run: `// For Compiling : packr2 && go build -o allow_permission_mac_go && packr2 clean`

## Python Script Alternative (Slower and heavier)
- Function: An alternative Python script for handling microphone permissions.
- For re-compilation: Download the Voice Logger directory and run:
 `pyinstaller --onefile --add-data "allow_button_images/allow_button_image.jpg:allow_button_images" --add-data "allow_button_images/allow_button_image1.jpg:allow_button_images" --add-data "allow_button_images/allow_button_image2.jpg:allow_button_images" --add-data "allow_button_images/allow_button_image3.jpg:allow_button_images" --add-data "allow_button_images/allow_button_image4.jpg:allow_button_images" --add-data "allow_button_images/allow_button_image5.jpg:allow_button_images" allow_permission_mac.py`
- The executable file will be inside the dist directory
  
## Disclaimer
- For educational purposes only.
- Not intended for malicious use.
- The author is not responsible for misuse or damage caused by this script.
