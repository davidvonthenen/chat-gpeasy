# Project: chat-gpeasy

This project provides to main goals:

- hooks the [OpenAI API](https://platform.openai.com/docs/api-reference) by creating a proxy interface which leverages the [Go SDK](github.com/sashabaranov/go-openai)
- an opinionated but easy to use wrapper for the [Go OpenAI SDK](github.com/sashabaranov/go-openai) that makes creating workflows in Go easier

This addresses issues raised with wanting to consume the [OpenAI API](https://platform.openai.com/docs/api-reference) but also preserve some level of content/data ownership with queries posed by clients using the API. You can find more details in this presentation (TODO: link to presentation) and [slides](bit.ly/3HdcEXL).

## Components

### OpenAI Proxy

There is an example proxy service located in the `./cmd/bin` folder which implements the [service hook interface](https://github.com/dvonthenen/chat-gpeasy/blob/proof-of-concept/pkg/proxy/interfaces/interface.go#L10-L40) by dumping the request/response parameters for each [OpenAI API](https://platform.openai.com/docs/api-reference) call to the console.

To run the example proxy:

```sh
$ cd ./cmd/bin
go run cmd.go
```

Then run one of the example clients in the `./examples` folder, such as the [Cumulative](./examples/cumulative) client which by default will connect to the proxy on your localhost.

```sh
$ cd ./examples/cumulative
go run cmd.go
```

### Opinionated OpenAI Clients

For convenience, this repository also implements a number of opinionated clients, which are called `personas`, that wrap the [Go OpenAI SDK](github.com/sashabaranov/go-openai) by are conceptually easier to wrap your brain around.

#### Simple Persona

This provides a very simple Q&A client. 1 question gives you 1 answer. Each Question or Answer does not affect or influence the next.

To create a client:

```go
persona, err := personas.NewSimpleChat()
if err != nil {
    fmt.Printf("personas.NewSimpleChat error: %v\n", err)
}
```

Initialize the client with the model you intent on using:

```go
(*persona).Init(interfaces.SkillTypeGeneric, openai.GPT3Dot5Turbo) // openai.GPT3Dot5Turbo is the default or use "" empty string
```

Ask ChatGPT a question:

```go
prompt := "Hello! How are you doing?"
choices, err := (*persona).Query(ctx, prompt)
if err != nil {
    fmt.Printf("persona.Query error: %v\n", err)
    // exit!
}
fmt.Printf("Me:\n%s\n", prompt)
fmt.Printf("\n\nChatGPT:\n%s\n", choices[0].Message.Content)
```

#### Cumulative Persona

Simple interface for your typical chatbot style (cumulative conversation building chat) client where the context of the conversation (aka the questions and answers) affect or influence the next. Like a real conversation...

This provides a very simple Q&A client.

To create a client:

```go
persona, err := personas.NewCumulativeChat()
if err != nil {
    fmt.Printf("personas.NewCumulativeChat error: %v\n", err)
}
```

Initialize the client with the model you intent on using:

```go
(*persona).Init(interfaces.SkillTypeGeneric, "")
```

Ask it question:

```go
prompt = "Tell me about Long Beach, CA."
choices, err = (*persona).Query(ctx, prompt)
if err != nil {
    fmt.Printf("persona.Query error: %v\n", err)
    os.Exit(1)
}
fmt.Printf("Me:\n%s\n", prompt)
fmt.Printf("\n\nChatGPT:\n%s\n", choices[0].Message.Content)
```

Refine the conversation with more instructions to ChatGPT:

```go
err = (*persona).AddDirective("I want more factual type data")
if err != nil {
    fmt.Printf("persona.AddDirective error: %v\n", err)
    os.Exit(1)
}
```

Refine the initial question to get a different response:

```go
prompt = "Now... tell me about Long Beach, CA."
choices, err = (*persona).Query(ctx, prompt)
if err != nil {
    fmt.Printf("persona.Query error: %v\n", err)
    os.Exit(1)
}
fmt.Printf("Me:\n%s\n", prompt)
fmt.Printf("\n\nChatGPT:\n%s\n", choices[0].Message.Content)
```

#### Advanced Persona

Provides more capabilities/functions on the `Cumulative` client.

To create a client:

```go
persona, err := personas.NewAdvancedChat()
if err != nil {
    fmt.Printf("personas.NewCumulativeChat error: %v\n", err)
}
```

Ask ChatGPT a question like before:

```go
prompt = "Tell me about Long Beach, CA."
choices, err = (*persona).Query(ctx, openai.ChatMessageRoleUser, prompt)
if err != nil {
    fmt.Printf("persona.Query error: %v\n", err)
    os.Exit(1)
}
fmt.Printf("Me:\n%s\n", prompt)
fmt.Printf("\n\nChatGPT:\n%s\n", choices[0].Message.Content)
```

Make a mistake, edit the conversation and regenerate the responses:

```go
fmt.Printf("Oooops... I goofed. I need to edit this...\n\n\n")
conversation, err := (*persona).GetConversation()
if err != nil {
    fmt.Printf("persona.GetConversation error: %v\n", err)
    os.Exit(1)
}

for pos, msg := range conversation {
    if strings.Contains(msg.Content, "Long Beach, CA") {
        prompt = "Tell me about Laguna Beach, CA."
        choices, err := (*persona).EditConversation(pos, prompt)
        if err != nil {
            fmt.Printf("persona.EditConversation error: %v\n", err)
            os.Exit(1)
        }
        fmt.Printf("Me:\n%s\n", prompt)
        fmt.Printf("\n\nChatGPT:\n%s\n", choices[0].Message.Content)
    }
}
```

## Examples

You can find a list of very simple main-style examples to consume this SDK in the [examples folder][examples-folder]. To run these examples, you need to change directory into an example you wish to run and then execute the `go` file in that directory. For example:

```sh
$ cd examples/advanced
go run cmd.go
```

Examples include:

- `Simple/`
  - [Simple](./examples/simple/simple) - A very simple Q then A client
  - [Expert](./examples/simple/expert) - A very simple Q then A client but posing as an expert
  - [DAN](./examples/simple/dan) - Easter egg!
- [Cumulative](./examples/cumulative) - Simple interface for your typical chatbot style (cumulative conversation building chat) client
- [Advanced](./examples/advanced) - Provides more capabilities/functions on the `Cumulative` client
- `Vanilla/`
  - [Chat Completion](./examples/vanilla/chatcompletion) - This is a plain vanilla OpenAI client of a Q and A client (just an example)
  - [Models](./examples/vanilla/models) - This is plain pure vanilla API call to the `GET models` API (just an example)
- [Symbl.ai](./examples/symbl) - a real-time streaming [Symbl.ai](https://symbl.ai/) + [ChatGPT](https://chat.openai.com/) intergration example

## Community

If you have any questions, feel free to contact me on Twitter [@dvonthenen][dvonthenen_twitter] or through our [Community Slack][slack].

This SDK is actively developed, and we love to hear from you! Please feel free to [create an issue][issues] or [open a pull request][pulls] with your questions, comments, suggestions, and feedback. If you liked our integration guide, please star our repo!

This library is released under the [Apache 2.0 License][license]

[examples-folder]: examples/
[issues]: https://github.com/dvonthenen/chat-gpeasy/issues
[pulls]: https://github.com/dvonthenen/chat-gpeasy/pulls
[license]: LICENSE
[slack]: https://join.slack.com/t/symbldotai/shared_invite/zt-4sic2s11-D3x496pll8UHSJ89cm78CA
