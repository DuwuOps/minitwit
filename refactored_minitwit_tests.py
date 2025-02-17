# -*- coding: utf-8 -*-
"""
    MiniTwit Tests
    ~~~~~~~~~~~~~~

    Tests a MiniTwit application.

    :refactored: (c) 2024 by HelgeCPH from Armin Ronacher's original unittest version
    :copyright: (c) 2010 by Armin Ronacher.
    :license: BSD, see LICENSE for more details.
"""
import random
import requests
import requests.cookies


# import schema
# import data
# otherwise use the database that you got previously
BASE_URL = "http://localhost:8000"

def register(http_session : requests.Session, username, password, password2=None, email=None):
    """Helper function to register a user"""
    if password2 is None:
        password2 = password
    if email is None:
        email = username + '@example.com'
    r = http_session.post(f'{BASE_URL}/register', data={
        'username':     username,
        'password':     password,
        'password2':    password2,
        'email':        email,
    }, allow_redirects=True, cookies=http_session.cookies, stream=True)
    return r

def login(http_session : requests.Session, username, password):
    """Helper function to login"""
    r = http_session.post(f'{BASE_URL}/login', data={
        'username': username,
        'password': password
    }, allow_redirects=True, cookies=http_session.cookies)
    return r

def register_and_login(http_session : requests.Session, username, password):
    """Registers and logs in in one go"""
    _ = register(http_session, username, password)
    return login(http_session, username, password)

def logout(http_session):
    """Helper function to logout"""
    r = http_session.get(f'{BASE_URL}/logout', allow_redirects=True, cookies=http_session.cookies)
    return r

def add_message(http_session, text):
    """Records a message"""
    r = http_session.post(f'{BASE_URL}/add_message', data={'text': text},
                                allow_redirects=True, cookies=http_session.cookies)
    if text:
        assert 'Your message was recorded' in r.text
    return r

# testing functions

def test_register():
    http_session = requests.Session()
    """Make sure registering works"""
    usern = 'user1'  + random.choice("abcdefghijklmonpqrstuvxyz") + random.choice("abcdefghijklmonpqrstuvxyz")
    r = register(http_session, usern, 'default')
    print("\n\n ======= COOKIES ======= \n\n")
    print("r:", r.cookies)
    print()
    print("http_session:", http_session.cookies)
    print("\n\n ======= ACTUAL RESPONSE ======= \n\n")
    print(r.text)
    assert 'You were successfully registered ' \
           'and can login now' in r.text
    r = register(http_session, 'user1', 'default')
    assert 'The username is already taken' in r.text
    r = register(http_session, '', 'default')
    assert 'You have to enter a username' in r.text
    r = register(http_session, 'meh', '')
    assert 'You have to enter a password' in r.text
    r = register(http_session, 'meh', 'x', 'y')
    assert 'The two passwords do not match' in r.text
    r = register(http_session, 'meh', 'foo', email='broken')
    assert 'You have to enter a valid email address' in r.text

# def test_login_logout():
#     http_session = requests.Session()
#     """Make sure logging in and logging out works"""
#     r = register_and_login(http_session, 'user1', 'default')
#     assert 'You were logged in' in r.text
#     r = logout(http_session)
#     assert 'You were logged out' in r.text
#     r = login(http_session, 'user1', 'wrongpassword')
#     assert 'Invalid password' in r.text
#     r = login(http_session, 'user2', 'wrongpassword')
#     assert 'Invalid username' in r.text

# def test_message_recording():
#     http_session = requests.Session()
#     """Check if adding messages works"""
#     _ = register_and_login(http_session, 'foo', 'default')
#     add_message(http_session, 'test message 1')
#     add_message(http_session, '<test message 2>')
#     r = requests.get(f'{BASE_URL}/')
#     assert 'test message 1' in r.text
#     assert '&lt;test message 2&gt;' in r.text

# def test_timelines():
#     http_session = requests.Session()
#     """Make sure that timelines work"""
#     _ = register_and_login(http_session, 'foo', 'default')
#     _ = add_message(http_session, 'the message by foo')
#     _ = logout(http_session)
#     _ = register_and_login(http_session, 'bar', 'default')
#     _ = add_message(http_session, 'the message by bar')
#     r = http_session.get(f'{BASE_URL}/public', cookies=http_session.cookies)
#     assert 'the message by foo' in r.text
#     assert 'the message by bar' in r.text

#     # bar's timeline should just show bar's message
#     r = http_session.get(f'{BASE_URL}/', cookies=http_session.cookies)
#     assert 'the message by foo' not in r.text
#     assert 'the message by bar' in r.text

#     # now let's follow foo
#     r = http_session.get(f'{BASE_URL}/foo/follow', allow_redirects=True, cookies=http_session.cookies)
#     assert 'You are now following &#34;foo&#34;' in r.text

#     # we should now see foo's message
#     r = http_session.get(f'{BASE_URL}/', cookies=http_session.cookies)
#     assert 'the message by foo' in r.text
#     assert 'the message by bar' in r.text

#     # but on the user's page we only want the user's message
#     r = http_session.get(f'{BASE_URL}/bar', cookies=http_session.cookies)
#     assert 'the message by foo' not in r.text
#     assert 'the message by bar' in r.text
#     r = http_session.get(f'{BASE_URL}/foo', cookies=http_session.cookies)
#     assert 'the message by foo' in r.text
#     assert 'the message by bar' not in r.text

#     # now unfollow and check if that worked
#     r = http_session.get(f'{BASE_URL}/foo/unfollow', allow_redirects=True, cookies=http_session.cookies)
#     assert 'You are no longer following &#34;foo&#34;' in r.text
#     r = http_session.get(f'{BASE_URL}/', cookies=http_session.cookies)
#     assert 'the message by foo' not in r.text
#     assert 'the message by bar' in r.text
