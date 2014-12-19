#!/usr/bin/env python
# encoding: utf-8

from gevent import socket, spawn, joinall,sleep
import random
import os


def ss_listen(s):
    
    while True:
        sleep(random.randint(0,5))
        s.send("hi "+os.urandom(15).encode('hex')+'\n')
        s.recv(1024)

jobs = []

print "Connectting...",
for x in xrange(3000):

    ss = socket.socket()
    ss.connect(('localhost', 12345))

    jobs.append(spawn(ss_listen, ss))

print "Done"
joinall(jobs)
