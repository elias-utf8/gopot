ssh-keygen -f "/home/$USER/.ssh/known_hosts" -R "[$1]:2222"
ssh "$1" -p 2222
