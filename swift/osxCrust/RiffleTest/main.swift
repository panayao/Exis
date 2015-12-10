// Testing

import Foundation

SetLogLevelDebug()
SetLocalFabric()

class TestingDomain: Domain {
    
    override func onJoin() {
        print("Subclass joined!")
        
        register("reg") { (args: Any) -> Any? in
            print("Received call! Args: \(args)")
            return nil
        }
        
        subscribe("sub") { (args: Any) in
            print("Received publish! \(args)")
        }
    }
    
    override func onLeave() {
        print("Subclass left!")
    }
}

TestingDomain(name: "xs.damouse").join()