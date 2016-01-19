//
//  main.swift
//  RiffleTest
//
//  Created by Mickey Barboi on 11/22/15.
//  Copyright © 2015 exis. All rights reserved.
//

/*
TODO:

    Integrate with main swiftRiffle lib for testing
    Make conditional compilers for ios and osx
    Cleanup and integrate new changes with goRiffle
    Implement Domain class in goRiffle
    Implment Call, Unreg, Unsub
*/

import Foundation
import CoreFoundation

#if os(Linux)
    import mantle
    import SwiftGlibc
    import Glibc
#else
    import Darwin.C
#endif

public protocol Delegate {
    func onJoin()
    func onLeave()
}


public class Domain {
    public var mantleDomain: UnsafeMutablePointer<Void>
    public var delegate: Delegate?
    
    public var invocations: [UInt64: [Any] -> ()] = [:]
    public var registrations: [UInt64: [Any] -> Any?] = [:]
    
    var deferreds: [UInt64: Deferred] = [:]
    var handlers: [UInt64: [Any] -> ()] = [:]
    
    
    public init(name: String) {
        mantleDomain = NewDomain(name.cString())
    }
    
    public init(name: String, superdomain: Domain) {
        mantleDomain = Subdomain(superdomain.mantleDomain, name.cString())
    }
    
    public func _subscribe(endpoint: String, _ types: [Any], fn: [Any] -> ()) -> Deferred {
        let hn = CBID()
        handlers[hn] = fn

        let d = Deferred(domain: self)
        Subscribe(self.mantleDomain, endpoint.cString(), d.cb, d.eb, hn, marshall(serializeArguments(types)))
        return d
    }
    
    public func _register(endpoint: String, _ types: [Any], fn: [Any] -> Any) -> Deferred {
        let hn = CBID()
        registrations[hn] = fn

        let d = Deferred(domain: self)
        Register(self.mantleDomain, endpoint.cString(), d.cb, d.eb, hn, marshall(types))
        return d
    }

    public func publish(endpoint: String, _ args: Any...) -> Deferred {
        let d = Deferred(domain: self)
        Publish(self.mantleDomain, endpoint.cString(), d.cb, d.eb, marshall(serializeArguments(args)))
        return d
    }
    
    public func call(endpoint: String, _ args: Any...) -> HandlerDeferred {
        let d = HandlerDeferred(domain: self)
        d.mantleDomain = self.mantleDomain
        Call(self.mantleDomain, endpoint.cString(), d.cb, d.eb, marshall(serializeArguments(args)))
        return d
    }
    
    public func receive() {
        while true {
            var (i, args) = decode(Receive(self.mantleDomain))
            
            if let d = deferreds[i] {
                // remove the deferred (should this ever be optional?)
                deferreds[d.cb] = nil
                deferreds[d.eb] = nil
                
                if d.cb == i {
                    d.callback(args)
                }
                
                if d.eb == i {
                    d.errback(args)
                }
            } else if let fn = handlers[i] {
                fn(args)
            } else if let fn = invocations[i] {
                fn(args)
            } else if let fn = registrations[i] {
                let resultId = args.removeAtIndex(0) as! Double
                
                // Optional serialization has some problems. This unwraps the result to avoid that particular issue
                if let ret = fn(args) {
                    // TODO: handle tuple returns
                    Yield(mantleDomain, UInt64(resultId), marshall([ret]))
                } else {
                    let empty: [Any] = []
                    Yield(mantleDomain, UInt64(resultId), marshall(empty))
                }
            }
        }
    }
    
    public func join() {
        let cb = CBID()
        let eb = CBID()
        
        Join(mantleDomain, cb, eb)
        
        handlers[cb] = { a in
            if let d = self.delegate {
                d.onJoin()
            } else {
                self.onJoin()
            }
        }

        handlers[eb] = { a in
            if let d = self.delegate {
                d.onLeave()
            } else {
                self.onLeave()
            }
        }
        
        handlers[eb] = { (a: Any) in
            print("Unable to join!")
        }
        
        receive()
    }
    
    
    // MARK: Delegate methods
    public func onJoin() { }
    public func onLeave() { }
}