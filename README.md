## What Is This?

I live in New York City and rely heavily on public transportation (since I neither like driving nor have ever really driven seriously in my life).

This is an Alexa skill that wraps around the NYC MTA Bus Time API which gives real time data on bus times and distances for a particular bus stop (all bus stops have a code).

There are phone apps that you can type a bus stop code in and get the same data, but I'd rather have a voice interface so I made this Alexa skill instead.

### TODO's

- Actually deploy the application to Amazon so that other people can use it.

### Installing dependencies for local development

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

