Temperature Server
===================

Listens for TCP data in specific format and saves it to mysql

It expects packets in following format

float;int;int;int$

It will recognize the packet regardless of the delivery time of the elements. In other words, they can be delivered in different points of the time, byte by byte. 



