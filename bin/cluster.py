#!/bin/python



import os
import sys
import sets
import re
from collections import defaultdict
from functools import total_ordering


def letters(init):
    first = ord(init)
    num = first
    prefix, iter = "", letters(init)
    while True:
        for x in xrange(0, 26):
            yield prefix + chr(num + x)
        prefix = iter.next()

class Node:
    NAMES = letters('a')

    def __init__(self, path):
        self.short_name = Node.NAMES.next()
        self.path = path
        tokens = self._tokens(path)
        self.tokens = list(self._filter(tokens))
        self.groups = []
        self.cluser = None

    def _tokens(self, path):        
        print "_________"
        print path
        path = re.sub("\.\w\w\w(\sFolder)?$", "", path)
        print path
        path = re.sub("[{}\._\-[\]]+", " ", path)
        print path
        path = re.sub("\s+", " ", path)
        print path
        return path.strip().split(" ")

    def _filter(self, tokens):        
        for token in tokens:
            if re.search("\d", token):
                continue
            if token.lower() == "playboy":
                continue 
            yield token

    def polygrams(self):
        total = len(self.tokens)
        for i in xrange(0, len(self.tokens)):
            for n in xrange(2, 4):
                j = i + n
                if j > len(self.tokens):
                    continue

                bigrams = self.tokens[i:j]
                bigram = ",".join(bigrams)
                yield Link(self, 20 - i, bigram)

@total_ordering
class Link:
    def __init__(self, node, order, bigram):
        self.node = node
        self.order = order
        self.bigram = bigram

    def __lt__(self, o):
        return self.order < o.order
    def __repr__(self):
        return "L<(%s) ->(%d) {%s}>" % (self.node.path, self.order, self.bigram)


@total_ordering
class Group:
    NAMES = letters('A')
    
    def __init__(self, bigram, links):
        self.short_name = Group.NAMES.next()
        self.bigram = bigram
        self.links = links
        self.nodes = set([l.node for l in self.links]) 
    def __lt__(self, o):
        return self.bigram < o.bigram or len(self.links) > len(o.links)
    def __eq__(self, o):
        return self.bigram == self.bigram
    def __hash__(self):
        return hash(self.bigram)
    def __repr__(self):
        return "G{%s(%d)}" % (self.bigram, len(self.links))


class Cluster:
    num = 0
    def __init__(self, group):
        self.name = "C(%d)" % Cluster.num
        Cluster.num += 1
        self.groups = set([ group ])
    def __repr__(self):
        return self.name


def MakeGroups(nodes):
    print 
    print "----NODES----"
    print
    m = defaultdict(set)

    for n in nodes:
        print "------"
        print n.path, " -> ", n.tokens
        for link in n.polygrams():
            print "   ", link
            m[link.bigram].add(link)

    return [ Group(k, v) for k, v in m.items() if len(v) > 1 ]


def BuildClusters(nodes):
    print 
    print "----NODES----"
    print
    clusters = []    
    for n in nodes:
        print "-------"
        print n

        best_cluster = None
        best_count = 0
        for c in clusters:
            print "   ", c
            for cg in c.groups:
                common = cg.nodes.intersection(g.nodes)
                print "  ", common
                if common > best_count:
                    best_count = common
                    best_cluster = cc

        if best_cluster is None:
            best_cluster = Cluster(g)
            clusters.append(best_cluster)

        if g.cluster != best_cluster:
            print "NEW CLUSTER: %s -> %s" % (g, best_cluster)
            g.cluster = best_cluster



nodes = [ Node(f) for f in os.listdir(os.getcwd()) ]
groups = MakeGroups(nodes)


def MakeDot(groups):
    result = ""
    result += "digraph {\n"
    for g in groups:
        result += """  %s [shape=box, label="%s"]; \n""" % (g.short_name, g.bigram)
        for l in g.links:
            result += """  %s; \n""" % (l.node.short_name)

    for g in groups:        
        for l in g.links:
            result += """ %s -> %s; \n""" % (l.node.short_name, g.short_name)

    result += "}\n"
    return result

dot = MakeDot(groups)

with open("graph.dot", 'w') as f:
    f.write(dot)
    

#BuildClusters(groups)


    
        


