// cfx is a set of tools that I use to work with cloudformation, particularly
// in a CI environment.
//
// Included currently are three commands: `show-parameters`, `update-parameters`,
// and `settle`. `show-parameters` and `update-parameters` do exactly what they
// sound like they do. `settle` will follow the events of a stack until the stack
// "settles" with a `_COMPLETE` or `_FAILED` event. In the event that it settles
// in a `FAILED` state, a non-zero exit code will be returned. This makes it
// suitable for use in a CI environment.
//
//   usage: cfx [<flags>] <command> [<args> ...]
//
//   Cloudformation Toolkit
//
//   Flags:
//     --help  Show context-sensitive help (also try --help-long and --help-man).
//
//   Commands:
//     help [<command>...]
//       Show help.
//
//     show-parameters --stack-name=STACK-NAME
//       Show parameters for a stack.
//
//       --stack-name=STACK-NAME  Name of the stack.
//
//     update-parameters --stack-name=STACK-NAME [<flags>]
//       Update parameters for a stack.
//
//       --stack-name=STACK-NAME    Name of the stack.
//       --capabilities=CAPABILITIES ...
//                                  Capabilities required to perform changes.
//       --parameter=PARAMETER ...  Parameters to change.
//
//     settle --stack-name=STACK-NAME [<flags>]
//       Wait for a stack to settle, tailing events.
//
//       --stack-name=STACK-NAME  Name of the stack.
//       --timeout=10m            Maximum time to wait until the stack is considered settled.
//       --poll-interval=5s       Interval at which to poll AWS for events
//
// License
//
// 3-clause BSD. A copy is included with the source.
package main
