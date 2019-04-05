#!/usr/bin/env python
# coding: utf-8

# In[ ]:


lines = open('/users/kevin/nodes.yml', 'r').readlines()
num1 = lines[10].split(':')
num2 = lines[11].split(':')
num1[1]=int(num1[1])+1
num2[1]=int(num2[1])+1
lines[10]=num1[0]
lines[10] += ": " + str(num1[1])
lines[11]=num2[0]
lines[11] += ": " + str(num2[1])
out = open('/users/kevin/nodes.yml', 'w')
for i in  lines:
    out.write(i)
    out.write("\n")

out.close()

