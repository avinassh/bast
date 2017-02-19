import os

import requests
import requests.auth

app_key = os.getenv("BAST_APP_KEY")
app_secret = os.getenv("BAST_APP_SECRET")
ua_string = os.getenv("BAST_USER_AGENT_STRING")
username = os.getenv("REDDIT_USERNAME")
password = os.getenv("REDDIT_PASSWORD")


def get_access_token():
    client_auth = requests.auth.HTTPBasicAuth(app_key, app_secret)
    post_data = {"grant_type": "password", "username": username,
                 "password": password}
    headers = {"User-Agent": ua_string}
    response = requests.post("https://www.reddit.com/api/v1/access_token",
                             auth=client_auth, data=post_data, headers=headers)
    return response.json().get("access_token")


if __name__ == '__main__':
    print(get_access_token())
