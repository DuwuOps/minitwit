import json
import sqlite3
import requests
from pathlib import Path
from contextlib import closing


BASE_URL = 'http://127.0.0.1:8000'
DATABASE = "/tmp/minitwit.db"
USERNAME = 'simulator'
PWD = 'super_safe!'
HEADERS = {
    'Connection': 'close',
    'Content-Type': 'application/json',
}


def init_db():
    """Creates the database tables."""
    with closing(sqlite3.connect(DATABASE)) as db:
        with open("queries/schema.users.sql") as fp:
            db.cursor().executescript(fp.read())
        db.commit()
        with open("queries/schema.message.sql") as fp:
            db.cursor().executescript(fp.read())
        db.commit()
        with open("queries/schema.follower.sql") as fp:
            db.cursor().executescript(fp.read())
        db.commit()
        with open("queries/schema.latest_processed.sql") as fp:
            db.cursor().executescript(fp.read())
        db.commit()


def get_csrf_token(session: requests.Session, path: str) -> str:
    """
    Issue a GET to `path` to pick up the `csrf_token` cookie,
    then return its value.
    """
    _ = session.get(BASE_URL + path)
    # csrf_token is set as an HttpOnly cookie by Echo's middleware
    token = session.cookies.get('csrf_token')
    if not token:
        raise RuntimeError(f"No csrf_token cookie after GET {path}")
    return token


def make_auth_session() -> requests.Session:
    """
    Returns a Session pre-configured with BasicAuth and
    a valid CSRF cookie (by GETting /register).
    """
    s = requests.Session()
    s.auth = (USERNAME, PWD)
    # seed the CSRF cookie (can be any GET endpoint that uses the middleware)
    get_csrf_token(s, '/register')
    return s


def test_latest():
    session = make_auth_session()
    csrf = session.cookies['csrf_token']

    # post something to update LATEST
    url = f"{BASE_URL}/register"
    data = {'username': 'test', 'email': 'test@test', 'pwd': 'foo'}
    params = {'latest': 1337}
    response = session.post(
        url,
        params=params,
        data=json.dumps(data),
        headers={**HEADERS, 'X-CSRF-Token': csrf}
    )
    assert response.ok

    # verify that latest was updated
    response = session.get(f"{BASE_URL}/latest", headers=HEADERS)
    assert response.ok
    assert response.json()['latest'] == 1337


def test_register():
    session = make_auth_session()
    csrf = session.cookies['csrf_token']

    username = 'a'
    email = 'a@a.a'
    pwd = 'a'
    data = {'username': username, 'email': email, 'pwd': pwd}
    params = {'latest': 1}
    response = session.post(
        f'{BASE_URL}/register',
        params=params,
        data=json.dumps(data),
        headers={**HEADERS, 'X-CSRF-Token': csrf}
    )
    assert response.ok

    # verify that latest was updated
    response = session.get(f'{BASE_URL}/latest', headers=HEADERS)
    assert response.json()['latest'] == 1


def test_create_msg():
    session = make_auth_session()
    csrf = session.cookies['csrf_token']

    username = 'a'
    data = {'content': 'Blub!'}
    url = f'{BASE_URL}/msgs/{username}'
    params = {'latest': 2}
    response = session.post(
        url,
        params=params,
        data=json.dumps(data),
        headers={**HEADERS, 'X-CSRF-Token': csrf}
    )
    assert response.ok

    # verify that latest was updated
    response = session.get(f'{BASE_URL}/latest', headers=HEADERS)
    assert response.json()['latest'] == 2


def test_get_latest_user_msgs():
    session = make_auth_session()
    username = 'a'

    query = {'no': 20, 'latest': 3}
    url = f'{BASE_URL}/msgs/{username}'
    response = session.get(url, headers=HEADERS, params=query)
    assert response.status_code == 200

    found = any(
        msg['content'] == 'Blub!' and msg['user'] == username
        for msg in response.json()
    )
    assert found

    # verify that latest was updated
    response = session.get(f'{BASE_URL}/latest', headers=HEADERS)
    assert response.json()['latest'] == 3


def test_get_latest_msgs():
    session = make_auth_session()
    username = 'a'
    query = {'no': 20, 'latest': 4}
    url = f'{BASE_URL}/msgs'
    response = session.get(url, headers=HEADERS, params=query)
    assert response.status_code == 200

    found = any(
        msg['content'] == 'Blub!' and msg['user'] == username
        for msg in response.json()
    )
    assert found

    # verify that latest was updated
    response = session.get(f'{BASE_URL}/latest', headers=HEADERS)
    assert response.json()['latest'] == 4


def test_register_b():
    session = make_auth_session()
    csrf = session.cookies['csrf_token']

    username = 'b'
    email = 'b@b.b'
    pwd = 'b'
    data = {'username': username, 'email': email, 'pwd': pwd}
    params = {'latest': 5}
    response = session.post(
        f'{BASE_URL}/register',
        params=params,
        data=json.dumps(data),
        headers={**HEADERS, 'X-CSRF-Token': csrf}
    )
    assert response.ok

    # verify that latest was updated
    response = session.get(f'{BASE_URL}/latest', headers=HEADERS)
    assert response.json()['latest'] == 5


def test_register_c():
    session = make_auth_session()
    csrf = session.cookies['csrf_token']

    username = 'c'
    email = 'c@c.c'
    pwd = 'c'
    data = {'username': username, 'email': email, 'pwd': pwd}
    params = {'latest': 6}
    response = session.post(
        f'{BASE_URL}/register',
        params=params,
        data=json.dumps(data),
        headers={**HEADERS, 'X-CSRF-Token': csrf}
    )
    assert response.ok

    # verify that latest was updated
    response = session.get(f'{BASE_URL}/latest', headers=HEADERS)
    assert response.json()['latest'] == 6


def test_follow_user():
    session = make_auth_session()
    csrf = session.cookies['csrf_token']

    username = 'a'
    url = f'{BASE_URL}/fllws/{username}'
    data = {'follow': 'b'}
    params = {'latest': 7}
    response = session.post(
        url,
        params=params,
        data=json.dumps(data),
        headers={**HEADERS, 'X-CSRF-Token': csrf}
    )
    assert response.ok

    data = {'follow': 'c'}
    params = {'latest': 8}
    response = session.post(
        url,
        params=params,
        data=json.dumps(data),
        headers={**HEADERS, 'X-CSRF-Token': csrf}
    )
    assert response.ok

    query = {'no': 20, 'latest': 9}
    response = session.get(url, headers=HEADERS, params=query)
    assert response.ok

    json_data = response.json()
    assert "b" in json_data["follows"]
    assert "c" in json_data["follows"]

    # verify that latest was updated
    response = session.get(f'{BASE_URL}/latest', headers=HEADERS)
    assert response.json()['latest'] == 9


def test_a_unfollows_b():
    session = make_auth_session()
    csrf = session.cookies['csrf_token']

    username = 'a'
    url = f'{BASE_URL}/fllws/{username}'

    #  first send unfollow command
    data = {'unfollow': 'b'}
    params = {'latest': 10}
    response = session.post(
        url,
        params=params,
        data=json.dumps(data),
        headers={**HEADERS, 'X-CSRF-Token': csrf}
    )
    assert response.ok

    # then verify that b is no longer in follows list
    query = {'no': 20, 'latest': 11}
    response = session.get(url, params=query, headers=HEADERS)
    assert response.ok
    assert 'b' not in response.json()['follows']

    # verify that latest was updated
    response = session.get(f'{BASE_URL}/latest', headers=HEADERS)
    assert response.json()['latest'] == 11