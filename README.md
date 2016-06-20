Golang client for tonomi.com
============================

![Build status](https://travis-ci.org/chemikadze/gonomi.svg?branch=master)

Manifest AST
------------

Provided by gonomi/manifest package. For example, the following manifest:

        application:
            components:
                x:
                    type: test.Component
                    interfaces:
                        myinterface:
                            mypin1: publish-signal(string)
                            mypin2: consume-signal(string)
                            mypin3: send-command()
                            mypin4: receive-command()
                        myrequired:
                            mypin: publish-signal(string)
                    required: [myrequired]

can be encoded like this:

    Application{CompositeComponent{
        Components: map[string]Component{
            "x": LeafComponent{
                Type:          Type{"test.Component"},
                Configuration: Configuration{},
                Interfaces: map[string]LeafInterface{
                    "myinterface": LeafInterface{
                        Pins: map[string]DirectedPinType{
                            "mypin1": {Sends, SignalPin{datatype.String{}}},
                            "mypin2": {Receives, SignalPin{datatype.String{}}},
                            "mypin3": {Sends, CommandPin{datatype.Record{}, datatype.Record{}, datatype.Record{}}},
                            "mypin4": {Receives, CommandPin{}},
                        },
                    },
                    "myrequired": LeafInterface{
                        Pins: map[string]DirectedPinType{
                            "mypin": {Sends, SignalPin{datatype.String{}}},
                        },
                        Required: true,
                    },
                },
            },
        }}}

Parser
------

Using parser is easy:

    import "github.com/chemikadze/gonomi/manifest"

    // ...
    app, err := Parse(":")
    if err == nil {
        t.Error("Error expected")
    }
    // use AST
    fmt.Println(app)

