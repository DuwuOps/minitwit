#!/bin/bash
##
# Website status check script
# Author: David Martin Sørensen
# Date: 13/03/2025
##

# Stolen from source: https://www.digitalocean.com/community/tutorials/how-to-use-discord-webhooks-to-get-notifications-for-your-website-status-on-ubuntu-18-04

##
# Discord webhook
# Change the 'your_discord_webhook_name' with your actual Discord Webhook
##
url="discord_webhook_url"

##
# List of websites to check
# To add more websites just use space as a separator, for example:
# websites_list="your_domain1 your_domain2 your_domain3"
##
websites_list="http://134.209.137.191/"

for website in ${websites_list} ; do
        status_code=$(curl --write-out %{http_code} --silent --output /dev/null -L ${website})

        if [[ "$status_code" -ne 200 ]] ; then
            # POST request to Discord Webhook with the domain name and the HTTP status code
            content="@everyone ${website} is down! It returned status code ${status_code}!!"
            curl -H "Content-Type: application/json" -X POST -d '{"content":"'"${content}"'", "allowed_mentions": { "parse": ["everyone"] }}'  $url
        else
            echo "${website} is running!"
        fi
done