import Riffle

class TourRegClient: Domain, Delegate  {
    
    override func onJoin() {
        /////////////////////////////////////////////////////////////////////////////////////
        // Example Tour Reg/Call Lesson 1 - our first basic example
        backend.call("myFirstFunc", "Hello").then { (s: String) in
            print(s) // Expects a String, like "Hello World"
        }
        // End Example Tour Reg/Call Lesson 1
        
        /////////////////////////////////////////////////////////////////////////////////////
        // Example Tour Reg/Call Lesson 2 Works - type enforcement good
        backend.call("iWantStrings", "Hi").then { (s: String) in
            print(s) // Expects a String, like "Thanks for saying Hi"
        }
        // End Example Tour Reg/Call Lesson 2 Works
        
        // Example Tour Reg/Call Lesson 2 Fails - type enforcement bad
        backend.call("iWantInts", "Hi").then { (s: String) in
            print(s)
        }.error { reason in
            print(reason) // Expects a String, like "Cumin: expecting primitive int, got str"
        }
        // End Example Tour Reg/Call Lesson 2 Fails
        
        // Example Tour Reg/Call Lesson 2 Wait Check - type enforcement on wait
        backend.call("iGiveInts", "Hi").then { (s: String) in
            print(s)
        }.error { reason in
            print(reason) // Expects a String, like "Cumin: expecting primitive str, got int"
        }
        // End Example Tour Reg/Call Lesson 2 Wait Check
        
        /////////////////////////////////////////////////////////////////////////////////////
        // Example Tour Reg/Call Lesson 3 Works - collections of types
        backend.call("iWantManyStrings", ["This", "is", "cool"]).then { (s: String) in
            print(s) // Expects a String, like "Thanks for 3 strings!"
        }
        // End Example Tour Reg/Call Lesson 3 Works
        
        // Example Tour Reg/Call Lesson 3 Fails - collections of types
        let notAllInts: [Any] = [0, 1, "two"]
        backend.call("iWantManyInts", notAllInts).then { (s: String) in
            print(s)
        }.error { reason in
            print(reason) // Expects a String, like "wamp.error.invalid_argument: Cumin: expecting primitive int, got string"
        }
        // End Example Tour Reg/Call Lesson 3 Fails
        
        /////////////////////////////////////////////////////////////////////////////////////
        // xample Tour Reg/Call Lesson 4 Basic Student - intro to classes
        class Student: Model {
            var name: String = "Student Name"
            var age: Int = 20
            var studentID: Int = 0
        }
        let s = Student()
        s.name = "John Smith"
        s.age = 18
        s.studentID = 1234
        backend.call("sendStudent", s)
        // nd Example Tour Reg/Call Lesson 4 Basic Student
        
        print("___SETUPCOMPLETE___")
    }
    
    override func onLeave() {
        print("Sender left")
    }
}
