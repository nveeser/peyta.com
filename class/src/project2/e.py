
f = open("1000.txt")


n = [ x.strip() for x in f.readlines() ]

for i in xrange(100):
	sub = n[i:i+100]
	print ",".join(sub) + ","

