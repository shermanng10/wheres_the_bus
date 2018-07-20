## What Is This?

I live in New York City and rely heavily on public transportation (since I neither like driving nor have ever really driven seriously in my life).

This is an Alexa skill that wraps around the NYC MTA Bus Time API which gives real time data on bus times and distances for a particular bus stop (all bus stops have a code).

There are phone apps that you can type a bus stop code in and get the same data, but I'd rather have a voice interface so I made this Alexa skill instead.

### TODO's

- Add data persistence (probably DynamoDB) so that users can store their stop preferences.
- Actually deploy the application to Amazon so that other people can use it.
- Think of a better name and invocation to activate the command on the Alexa.

### Installing dependencies

First make sure that you have dep installed.

For Mac:
```bash
$ brew install dep
$ brew upgrade dep
```

Then just install the dependencies
```bash
$ dep ensure
```

### Hosting it yourself

Since I haven't actually deployed this application for public use on Amazon yet, you'd need to compile, then package the binary and zip it.

```bash
$ make build
```

Then create a Lambda function and upload it. Once that's done you create an Alexa skill and hook it up to your Lambda function (you can even specify your own invocations as long as the intent command is wheres_the_bus which has an optional "slot" called stopCode).

You can find the instructions online but I'm probably going to actually deploy it in the upcoming month/week after my vacation.
