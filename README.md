# blinkythingy

The Swiss Army knife of IoT lighting

## example config

Server:

```yaml
fetchers:
    - type: color
      color: red
      number: 1
    - type: color
      color: green
      number: 1
displays:
    - type: http
      addr: :2000
```

Client:

```yaml
fetchers:
    - type: http
      url: http://localhost:2000
displays:
    - type: debug
```

A Jenkins job and Github PRs displayed on a blinky tape, separated by 3 empty spaces:

```yaml
fetchers:
    - type: jenkins
      host: 'example.com'
      username: 'aaaaaaaaaaaaa'
      password: 'xxxxxxxxxx'
      job: 'myjob'
    - type: color
      number: 3
    - type: github
      project: docker/docker
      username: 'aaaaaaaaaa'
      password: 'xxxxxxxxxxxx'
      query: 'label:status/2-code-review'
displays:
    - type: blinky
      path: /dev/ttyACM0
```
