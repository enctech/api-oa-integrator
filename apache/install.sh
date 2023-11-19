#!/bin/bash

# Update the package list
sudo apt update

# Install Apache
sudo apt install -y apache2

# Start Apache
sudo service apache2 start

# Enable Apache to start at boot
sudo systemctl enable apache2

# Open port 80 in the firewall (if using ufw)
sudo ufw allow 1323

# Enable Necessary Apache Modules:
sudo a2enmod proxy
sudo a2enmod proxy_http

mv ./default.conf /etc/apache2/sites-available/000-default.conf

# Restart Apache
sudo systemctl restart apache2

# Display status message
echo "Apache is now installed and running. You can access it at http://localhost"
