#!/usr/bin/env python
# coding: utf-8

# In[ ]:


lines = open('nodes2.yaml', 'r').readlines()
num1 = lines[11].split(':')
num2 = lines[12].split(':')
num1[1]=int(num1[1])+1
num2[1]=int(num2[1])+1
lines[11]=num1[0]
lines[11] += ": " + str(num1[1])
lines[12]=num2[0]
lines[12] += ": " + str(num2[1])
out = open('nodes2.yaml', 'w')
for i in  lines:
    out.write(i)
    out.write("\n")

out.close()

