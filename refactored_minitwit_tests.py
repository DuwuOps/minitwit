# -*- coding: utf-8 -*-
"""
    MiniTwit Tests
    ~~~~~~~~~~~~~~

    Tests a MiniTwit application.

    :refactored: (c) 2024 by HelgeCPH from Armin Ronacher's original unittest version
    :copyright: (c) 2010 by Armin Ronacher.
    :license: BSD, see LICENSE for more details.
"""
import re
import requests


# import schema
# import data
# otherwise use the database that you got previously
BASE_URL = "http://localhost:8000"

def get_csrf_token(session: requests.Session, path: str) -> str:
    """Fetch the form at `path` and extract the CSRF token from the hidden input."""
    response = session.get(f"{BASE_URL}{path}")
    # look for: <input type="hidden" name="_csrf" value="...">
    m = re.search(r'name="_csrf"\s+value="([^"]+)"', response.text)
    if not m:
        raise RuntimeError(f"Unable to find CSRF token on {path}")
    return m.group(1)

def register(username, password, password2=None, email=None):
    """Helper function to register a user"""
    if password2 is None:
        password2 = password
    if email is None:
        email = username + '@example.com'

    session = requests.Session()
    token = get_csrf_token(session, "/register")

    data = {
        "_csrf":    token,
        "username": username,
        "password": password,
        "password2": password2,
        "email":     email,
    }
    response = session.post(f"{BASE_URL}/register", data=data, allow_redirects=True)
    return response

def login(username, password):
    """Helper function to login"""
    session = requests.Session()
    token = get_csrf_token(session, "/login")

    data = {
        "_csrf":    token,
        "username": username,
        "password": password,
    }
    response = session.post(f"{BASE_URL}/login", data=data, allow_redirects=True)
    return response, session

def register_and_login(username, password):
    """Registers and logs in in one go"""
    register(username, password)
    return login(username, password)

def logout(http_session):
    """Helper function to logout"""
    return http_session.get(f'{BASE_URL}/logout', allow_redirects=True)

def add_message(http_session, text):
    """Records a message"""
    token = get_csrf_token(http_session, "/")
    data = {
        "_csrf": token,
        "text":  text,
    }
    response = http_session.post(f"{BASE_URL}/add_message", data=data, allow_redirects=True)
    if text:
        assert "Your message was recorded" in response.text
    return response

# testing functions

def test_register():
    """Make sure registering works"""
    response = register('user1', 'default')
    assert 'You were successfully registered ' \
           'and can login now' in response.text
    response = register('user1', 'default')
    assert 'The username is already taken' in response.text
    response = register('', 'default')
    assert 'You have to enter a username' in response.text
    response = register('meh', '')
    assert 'You have to enter a password' in response.text
    response = register('meh', 'x', 'y')
    assert 'The two passwords do not match' in response.text
    response = register('meh', 'foo', email='broken')
    assert 'You have to enter a valid email address' in response.text

def test_login_logout():
    """Make sure logging in and logging out works"""
    response, http_session = register_and_login('user1', 'default')
    assert 'You were logged in' in response.text
    response = logout(http_session)
    assert 'You were logged out' in response.text
    response, _ = login('user1', 'wrongpassword')
    assert 'Invalid password' in response.text
    response, _ = login('user2', 'wrongpassword')
    assert 'Invalid username' in response.text

def test_message_recording():
    """Check if adding messages works"""
    _, http_session = register_and_login('foo', 'default')
    add_message(http_session, 'test message 1')
    add_message(http_session, '<test message 2>')
    response = requests.get(f'{BASE_URL}/')
    assert 'test message 1' in response.text
    assert '&lt;test message 2&gt;' in response.text

def test_timelines():
    """Make sure that timelines work"""
    _, http_session = register_and_login('foo', 'default')
    add_message(http_session, 'the message by foo')
    logout(http_session)
    _, http_session = register_and_login('bar', 'default')
    add_message(http_session, 'the message by bar')
    response = http_session.get(f'{BASE_URL}/public')
    assert 'the message by foo' in response.text
    assert 'the message by bar' in response.text

    # bar's timeline should just show bar's message
    response = http_session.get(f'{BASE_URL}/')
    assert 'the message by foo' not in response.text
    assert 'the message by bar' in response.text

    # now let's follow foo
    response = http_session.get(f'{BASE_URL}/foo/follow', allow_redirects=True)
    assert 'You are now following &#34;foo&#34;' in response.text

    # we should now see foo's message
    response = http_session.get(f'{BASE_URL}/')
    assert 'the message by foo' in response.text
    assert 'the message by bar' in response.text

    # but on the user's page we only want the user's message
    response = http_session.get(f'{BASE_URL}/bar')
    assert 'the message by foo' not in response.text
    assert 'the message by bar' in response.text
    response = http_session.get(f'{BASE_URL}/foo')
    assert 'the message by foo' in response.text
    assert 'the message by bar' not in response.text

    # now unfollow and check if that worked
    response = http_session.get(f'{BASE_URL}/foo/unfollow', allow_redirects=True)
    assert 'You are no longer following &#34;foo&#34;' in response.text
    response = http_session.get(f'{BASE_URL}/')
    assert 'the message by foo' not in response.text
    assert 'the message by bar' in response.text
