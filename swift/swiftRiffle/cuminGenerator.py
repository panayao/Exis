'''
Generator for Cumin functions. 
'''

import os

# Header for the cumin file 
header = '''
// 
// Cumin generic wrapper functions. Each overloaded function corresponds to a handler with a different number of 
// arguments and return types
//
// Generated by cuminGenerator.py
// 

import Foundation

// Converter operator. Attempts to convert the object on the right to the type given on the left
// Just here to make the cumin conversion functions just the smallest bit clearer
infix operator <- {
associativity right
precedence 155
}

func <- <T: Property> (t:T.Type, object: Any) -> T {
    // Deserialize is implemented as part of the Convertible protocol. All properties implement Convertible
    return T.deserialize(object) as! T
}

// Used only in this file to shorten the length of the method signatures
public typealias PR = Property

public extension Domain {

'''

PRODUCTION = 'Riffle/Cumin.swift'
DEV = 'cumin.txt'

generics = ['A', 'B', 'C', 'D', 'E', 'F', 'G', 'H', 'I', 'J']
returns = ['R', 'S', 'T', 'U', 'V', 'X', 'Y', 'Z']

handlerTemplate = '\tpublic func %s<%s>(endpoint: String, _ fn: (%s) -> (%s)) -> Deferred {\n\t\treturn _%s(endpoint) { a in return fn(%s) }\n\t}'
callTemplate = '\tpublic func %s<%s>(endpoint: String, _ fn: (%s) -> (%s)) -> Deferred {\n\t\treturn _%s(endpoint) { a in return fn(%s) }\n\t}'


def renderCaller(template, name, args, ret, renderingArrays):
    cumin = ', '.join(["%s.self <- a[%s]" % (x, i) for i, x in enumerate(args)])
    both = ', '.join([x + ": PR" for x in args] + ret)
    args = ', '.join(args)
    ret = ', '.join(ret)

    return (template % (name, both, args, ret, name, cumin)).replace("<>", "")

def main():
    r, s, n = [], [], []

    for j in range(2):  # The number of return types
        for i in range(0, 7):  # Number of parameters
            if j == 0:
                s.append(renderCaller(handlerTemplate, 'subscribe', generics[:i], returns[:j], False))
                # n.append(renderCaller(callTemplate, 'call', generics[:i], returns[:j], False))

            r.append(renderCaller(handlerTemplate, 'register', generics[:i], returns[:j], False))

    with open(os.path.join(os.getcwd(), PRODUCTION), 'w') as f:
        f.write(header)
        e = r + s + n

        [f.write(x + '\n\n') for x in e]
        f.write("}\n\n")

if __name__ == '__main__':
    main()