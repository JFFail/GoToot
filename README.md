# GoToot

CLI Mastodon client. Currently in extremely early stages. Expecting there to be a `client.json` file in the same directory. The JSON file should have two entries: the access token from Mastodon and the URL of your instance. You can generate an access token at: `https://yourmastoinstance.com/settings/applications`. The instance listed in the JSON file should start with "https" and not include the trailing slash. Obviously this will all be part of the application later if we get that far. A sample JSON file:

    {
        "access_token": "tokenGoesHere",
        "instance": "https://mastodon.social"
    }

## Current

Currently implemented:

- Toots
- Toots w/ CW
- Viewing Home timeline
- Viewing Local timeline

## To Do

Still need to add:

- Viewing Notifications
- Boosts
- Favorites
- Better authentication
- Viewing Favorites?
- CLI toot for non-interactive posting?
