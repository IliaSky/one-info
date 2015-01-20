#!/usr/bin/env python

import sys
import email
import imaplib
import getpass
import re
import random
from json import dumps as to_json

IMAP_SERVER = 'imap.gmail.com'
EMAIL_ACCOUNT = "iliasky.test@gmail.com"
PASSWORD = getpass.getpass()
# PASSWORD = ''
# EMAIL_FOLDER = "INBOX"
# https://www.google.com/settings/security/lesssecureapps
OUTPUT_FILE = 'D:/code/go/one-info/experiments/python-version/tmp.json'


def expected_group(content, groups=['cb', '0 work', 'Accounts', 'old-sw']):
    return random.choice(groups)


def parse_mail(mail):
    mail = str(mail)
    msg = email.message_from_string(mail)
    return {"headers": dict(msg), "body": msg.get_payload()}

def process_mailbox(M):
    """
    Dump all emails in the folder to files in output file.
    """
    mails = []

    rv, data = M.search(None, "ALL")
    if rv != 'OK':
        print("No messages found!")
        return
    
    for num in data[0].split():
        rv, data = M.fetch(num, '(RFC822)')
        if rv != 'OK':
            return print("ERROR getting message", num)
            
        # print("Writing message ", num)

        # with open(OUTPUT_DIRECTORY +'-'+str(num)+ '.stupid', 'wb') as f2:
        #     f2.write(data[0][1])
        # with open(OUTPUT_DIRECTORY +'-'+str(num)+ '.stupid', 'r') as f3:
        #     mails += [parse_mail(f3.read())]

        mails += [parse_mail(data[0][1])]


    with open(OUTPUT_FILE, 'wb') as f:
        f.write(bytes(to_json(mails), 'UTF-8'))
    

def folder(M, name):
    mails = []
    print("Open mailbox: ", name)
    if " " in name or "-" in name:
        return print("skipping")

    rv, data = M.select('"' + name + '"')
    if rv == 'OK':
        print("Processing mailbox: ", name)

    rv, data = M.search(None, "ALL")
    if rv != 'OK':
        print("No messages found!")
        return
    
    for num in data[0].split():
        rv, data = M.fetch(num, '(RFC822)')
        if rv != 'OK':
            return print("ERROR getting message", num)

        mails += [parse_mail(data[0][1])]

    M.close()
    return mails

    
def save_grouped_mails(M):
    mails = [folder(M, directory) for directory in directories(M)]

    with open(OUTPUT_DIRECTORY + '.json', 'wb') as f:
        f.write(bytes(to_json(mails), 'UTF-8'))


def directories(M):
    ok, data = M.list()
    if ok != 'OK':
        return print("Could not get directories")
    return [e for e in [str(e).split(' "/" ')[1][1:-2] for e in data] if '[GMAIL]' not in e and e != 'st-fb']

def main():
    M = imaplib.IMAP4_SSL(IMAP_SERVER)
    M.login(EMAIL_ACCOUNT, PASSWORD)
    # rv, data = M.select(EMAIL_FOLDER)
    # if rv == 'OK':
    #     print("Processing mailbox: ", EMAIL_FOLDER)
        # process_mailbox(M)
        # print(M.list())
        # print(directories(M))
    save_grouped_mails(M)
    #     M.close()
    # else:
    #     print("ERROR: Unable to open mailbox ", rv)
    M.logout()

if __name__ == "__main__":
    main()
