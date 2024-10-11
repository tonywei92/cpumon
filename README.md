# cpumon

![CPU Monitor Screenshot](screenshot.png)

`cpumon` is a utility designed to monitor your CPU usage and provide notifications through beeps and/or by executing a specified command.

For example, it can be useful for sending an email by triggering an API with a `curl` command that corresponds to this action.

## Usage

Example:

```sh
# call `echo` command if CPU is more than 80% usage
cpumon -notify-on-cpu-more-than=80 -no-beep -notify-with="echo CPU usage high > output.log"
```

### Flags

- `-interval`: Interval in seconds (default: 1)
- `-single`: Monitor single CPU (default: true)
- `-notify-on-cpu-more-than`: Notify when CPU usage exceeds this percentage (default: 0)
- `-no-beep`: Disable beep notification (default: false)
- `-notify-with`: Command to run on notification (default: "")

## Contributing

Thank you for considering contributing to `cpumon`! We welcome contributions in the form of bug reports, feature requests, and pull requests. Please see [CONTRIBUTING.md](CONTRIBUTING.md) for guidelines on how to contribute to this project.

## License

This project is licensed under the MIT License. See the [LICENSE](LICENSE) file for details.
