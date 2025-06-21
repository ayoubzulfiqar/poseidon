import platform
import os
import sys

WEBSITE_LIST = [
    "www.facebook.com",
    "facebook.com",
    "www.twitter.com",
    "twitter.com",
    "www.instagram.com",
    "instagram.com",
    "www.youtube.com",
    "youtube.com",
    "www.netflix.com",
    "netflix.com",
    "www.reddit.com",
    "reddit.com",
    "www.tiktok.com",
    "tiktok.com",
    "www.pinterest.com",
    "pinterest.com",
    "www.linkedin.com",
    "linkedin.com"
]

REDIRECT_IP = "127.0.0.1"

def get_hosts_path():
    if platform.system() == "Windows":
        return r"C:\Windows\System32\drivers\etc\hosts"
    else:
        return "/etc/hosts"

def block_websites():
    hosts_path = get_hosts_path()
    try:
        with open(hosts_path, "r+") as file:
            content = file.readlines()
            file.seek(0)
            
            filtered_content = []
            for line in content:
                is_our_blocked_line = False
                for site in WEBSITE_LIST:
                    if f"{REDIRECT_IP} {site}" in line or f"{REDIRECT_IP}\t{site}" in line:
                        is_our_blocked_line = True
                        break
                if not is_our_blocked_line:
                    filtered_content.append(line)
            
            file.writelines(filtered_content)
            
            for website in WEBSITE_LIST:
                block_line_space = f"{REDIRECT_IP} {website}\n"
                block_line_tab = f"{REDIRECT_IP}\t{website}\n"
                
                if block_line_space not in "".join(filtered_content) and \
                   block_line_tab not in "".join(filtered_content):
                    file.write(block_line_space)
            
            file.truncate()
        print("Distracting websites blocked successfully!")
        print("Please flush your DNS cache for changes to take effect.")
        print("  Windows: ipconfig /flushdns")
        print("  Linux/macOS: sudo dscacheutil -flushcache; sudo killall -HUP mDNSResponder")

    except PermissionError:
        print("Error: Permission denied. Please run this script as an administrator/root.")
    except Exception as e:
        print(f"An error occurred: {e}")

def unblock_websites():
    hosts_path = get_hosts_path()
    try:
        with open(hosts_path, "r+") as file:
            content = file.readlines()
            file.seek(0)
            
            new_content = []
            for line in content:
                is_our_blocked_line = False
                for site in WEBSITE_LIST:
                    if f"{REDIRECT_IP} {site}" in line or f"{REDIRECT_IP}\t{site}" in line:
                        is_our_blocked_line = True
                        break
                if not is_our_blocked_line:
                    new_content.append(line)
            
            file.writelines(new_content)
            file.truncate()
        print("Distracting websites unblocked successfully!")
        print("Please flush your DNS cache for changes to take effect.")
        print("  Windows: ipconfig /flushdns")
        print("  Linux/macOS: sudo dscacheutil -flushcache; sudo killall -HUP mDNSResponder")

    except PermissionError:
        print("Error: Permission denied. Please run this script as an administrator/root.")
    except Exception as e:
        print(f"An error occurred: {e}")

def main():
    print("Website Blocker")
    print("-----------------")
    print("NOTE: This script requires administrator/root privileges to modify the hosts file.")

    while True:
        print("\nChoose an option:")
        print("1. Block distracting websites")
        print("2. Unblock websites")
        print("3. Exit")

        choice = input("Enter your choice (1/2/3): ")

        if choice == '1':
            block_websites()
        elif choice == '2':
            unblock_websites()
        elif choice == '3':
            print("Exiting program. Goodbye!")
            break
        else:
            print("Invalid choice. Please enter 1, 2, or 3.")

if __name__ == "__main__":
    main()

# Additional implementation at 2025-06-21 01:18:10
import datetime
import time
import os
import sys

# --- Configuration ---
# Determine the hosts file path based on the operating system
if sys.platform.startswith('win'):
    HOSTS_PATH = r"C:\Windows\System32\drivers\etc\hosts"
else:
    HOSTS_PATH = "/etc/hosts"

REDIRECT_IP = "127.0.0.1"
WEBSITE_LIST = [
    "www.facebook.com", "facebook.com",
    "www.twitter.com", "twitter.com",
    "www.instagram.com", "instagram.com",
    "www.youtube.com", "youtube.com",
    "www.reddit.com", "reddit.com",
    "www.tiktok.com", "tiktok.com"
]

# Blocking hours (24-hour format)
# Example: Block from 9 AM (09) to 5 PM (17)
START_HOUR = 9
END_HOUR = 17

LOG_FILE = "website_blocker.log"

# --- Helper Functions ---

def log_event(message):
    """Logs a timestamped message to the specified log file."""
    timestamp = datetime.datetime.now().strftime("%Y-%m-%d %H:%M:%S")
    log_entry = f"[{timestamp}] {message}\n"
    try:
        with open(LOG_FILE, "a") as f:
            f.write(log_entry)
    except IOError as e:
        print(f"Error writing to log file: {e}")

def is_blocking_time():
    """Checks if the current time falls within the defined blocking hours."""
    now = datetime.datetime.now().hour
    return START_HOUR <= now < END_HOUR

def block_websites(hosts_path, redirect_ip, website_list):
    """Adds website entries to the hosts file to block them."""
    try:
        with open(hosts_path, "r+") as file:
            content = file.readlines()
            file.seek(0) # Go to the beginning of the file
            blocked_count = 0
            
            # Write back lines that are not our block entries
            for line in content:
                # Keep existing lines unless they are our block entries
                if not any(f"{redirect_ip} {site}" in line for site in website_list):
                    file.write(line)
            
            # Add new block entries for websites not already present or correctly formatted
            for website in website_list:
                block_entry = f"{redirect_ip} {website}\n"
                # Check if the exact block entry already exists in the original content
                if block_entry not in content:
                    file.write(block_entry)
                    blocked_count += 1
            file.truncate() # Remove any remaining old content if file became shorter
        
        if blocked_count > 0:
            log_event(f"Blocked {blocked_count} new website entries.")
        else:
            log_event("Websites already blocked or no new websites to block.")
    except IOError as e:
        log_event(f"Error blocking websites: {e}. Ensure script has administrative/root privileges.")
    except Exception as e:
        log_event(f"An unexpected error occurred during blocking: {e}")

def unblock_websites(hosts_path, website_list):
    """Removes website entries from the hosts file to unblock them."""
    try:
        with open(hosts_path, "r+") as file:
            content = file.readlines()
            file.seek(0) # Go to the beginning of the file
            unblocked_count = 0
            for line in content:
                # Write back lines that are not our block entries
                if not any(f"{REDIRECT_IP} {site}" in line for site in website_list):
                    file.write(line)
                else:
                    unblocked_count += 1 # Count lines that were removed
            file.truncate() # Remove any remaining old content if file became shorter
        
        if unblocked_count > 0:
            log_event(f"Unblocked {unblocked_count} website entries.")
        else:
            log_event("Websites already unblocked or no entries to remove.")
    except IOError as e:
        log_event(f"Error unblocking websites: {e}. Ensure script has administrative/root privileges.")
    except Exception as e:
        log_event(f"An unexpected error occurred during unblocking: {e}")

# --- Main Logic ---

if __name__ == "__main__":
    log_event("Website Blocker started.")
    log_event(f"Blocking active from {START_HOUR}:00 to {END_HOUR}:00.")

    # Initialize block_active based on current time
    block_active = False
    if is_blocking_time():
        log_event("Initial check: Currently within blocking hours. Applying block.")
        block_websites(HOSTS_PATH, REDIRECT_IP, WEBSITE_LIST)
        block_active = True
    else:
        log_event("Initial check: Currently outside blocking hours. Ensuring unblock.")
        unblock_websites(HOSTS_PATH, WEBSITE_LIST)
        block_active = False

    while True:
        try:
            if is_blocking_time():
                if not block_active:
                    log_event("Entering blocking hours. Activating block.")
                    block_websites(HOSTS_PATH, REDIRECT_IP, WEBSITE_LIST)
                    block_active = True
            else:
                if block_active:
                    log_event("Exiting blocking hours. Deactivating block.")
                    unblock_websites(HOSTS_PATH, WEBSITE_LIST)
                    block_active = False
            
            # Check every 5 minutes (300 seconds)
            time.sleep(300) 

        except KeyboardInterrupt:
            log_event("Website Blocker stopped by user (KeyboardInterrupt).")
            # Optional: Unblock websites on exit if desired, otherwise they remain blocked until next time check
            # unblock_websites(HOSTS_PATH, WEBSITE_LIST) 
            break
        except Exception as e:
            log_event(f"An unexpected error occurred in the main loop: {e}")
            time.sleep(60) # Wait a bit before retrying after an error